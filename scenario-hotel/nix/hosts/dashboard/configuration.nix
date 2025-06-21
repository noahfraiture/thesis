{
  config,
  pkgs,
  ...
}:
let
  pkgs-dashboard =
    import
      (builtins.fetchTarball {
        url = "https://github.com/NixOS/nixpkgs/archive/0c159930e7534aa803d5cf03b27d5c86ad6050b7.tar.gz";
        sha256 = "0bqyjsgyf8cqq1w4ags8smmd1da40lrh57q8f3d8hmyp35aph7p9";
      })
      {
        config = config.nixpkgs.config;
        system = "x86_64-linux";
      };
in
{
  imports = [
    ../hardware-configuration.nix
  ];

  nix.settings = {
    experimental-features = "nix-command flakes";
    auto-optimise-store = true; # optimise at every build
  };

  boot.loader.systemd-boot.enable = true;
  boot.loader.efi.canTouchEfiVariables = true;

  networking.networkmanager.enable = true;
  time.timeZone = "Europe/Brussels";
  i18n.defaultLocale = "en_US.UTF-8";
  services.xserver.xkb = {
    layout = "us";
    variant = "";
  };

  environment.variables = {
    EDITOR = "hx";
  };

  services.openssh = {
    enable = true;
    ports = [ 22 ];
    settings = {
      PasswordAuthentication = true;
      AllowUsers = [
        "env-admin"
        "operator"
      ];
    };
  };

  users.users = {
    env-admin = {
      isNormalUser = true;
      description = "Admin of the CTF";
      extraGroups = [
        "networkmanager"
        "wheel"
      ];
      password = "env-admin";
    };
  };

  nixpkgs.config.allowUnfree = true;

  environment.systemPackages = [
    pkgs.helix
    pkgs.node-red
    pkgs.netcat-gnu
    pkgs-dashboard.mosquitto
    pkgs.python311
  ];

  system.stateVersion = "24.11";

  # === SPECIFIC ===

  networking.hostName = "dashboard";

  users.users.operator = {
    isNormalUser = true;
    description = "Operator of dashboards";
    extraGroups = [
    ];
    password = "password123";
  };

  networking = {
    firewall = {
      enable = true;
      allowedTCPPorts = [
        1880
        1883
      ];
      allowedUDPPorts = [ 1881 ];
    };
  };

  systemd.services.mosquittoStart = {
    enable = true;
    wantedBy = [ "multi-user.target" ];
    description = "Web server";
    path = [ "/run/current-system/sw" ];
    script = "mosquitto -v";
    serviceConfig = {
      User = "env-admin";
    };
  };

  services.node-red = {
    enable = true;
    openFirewall = true;
  };

  systemd.services.turnOff = {
    enable = true;
    wantedBy = [ "multi-user.target" ];
    description = "Web server";
    path = [ "/run/current-system/sw" ];
    script = ''
      while true; do
          mosquitto_pub -h localhost -p 1883 -t light/control -m "off"
          sleep 2
      done
    '';
    serviceConfig = {
      User = "env-admin";
    };
  };

  home-manager = {
    users."operator" = {
      programs.home-manager.enable = true;
      programs.bash.enable = true;
      home = {
        username = "operator";
        homeDirectory = "/home/operator";
        stateVersion = "24.11";
      };
    };
  };
  environment.variables.BEARER = "0c497d24da14cde5f5b947f7920d4df20189d9b1ad3302b9b31afd1ec1918f1b";

  system.activationScripts.history = pkgs.lib.stringAfter [ "users" ] ''
    echo 'curl http://dashboard/next -H "Authorization: Bearer $BEARER"' > /home/operator/.bash_history
    rm -r /home/operator/robot
    cp /home/env-admin/nixos/hosts/dashboard/robot /home/operator/robot
    chmod +x /home/operator/robot
  '';

}
