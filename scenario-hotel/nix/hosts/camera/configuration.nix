{
  pkgs,
  ...
}:
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
    pkgs.python311
    pkgs.python311Packages.requests
    pkgs.python311Packages.flask
    pkgs.python311Packages.pycryptodome
  ];

  system.stateVersion = "24.11";

  # === SPECIFIC ===

  networking.hostName = "camera";

  networking = {
    firewall = {
      enable = true;
      allowedTCPPorts = [ 8888 ];
    };
  };

  system.activationScripts.move_web = pkgs.lib.stringAfter [ "users" ] ''
    rm -f /home/env-admin/hikvision
    cp /home/env-admin/nixos/hosts/camera/hikvision /home/env-admin/hikvision
    chown -R env-admin:users /home/env-admin/hikvision
    chmod +x /home/env-admin/hikvision
  '';

  systemd.services.hikivision = {
    enable = true;
    wantedBy = [ "multi-user.target" ];
    description = "Camera server";
    path = [
      "/run/current-system/sw"
    ];
    script = ''/home/env-admin/hikvision'';
    serviceConfig = {
      User = "env-admin";
    };
  };

}
