{
  pkgs,
  config,
  ...
}:
let
  systemdPkgs =
    import
      (builtins.fetchTarball {
        url = "https://github.com/NixOS/nixpkgs/archive/3b05df1d13c1b315cecc610a2f3180f6669442f0.tar.gz";
        sha256 = "1dr7kfdl4wvxhml4hd9k77xszl55vbjbb6ssirs2qv53mgw8c24w";
      })
      {
        config = config.nixpkgs.config;
        system = "x86_64-linux";
      };
  serverPkgs = pkgs.buildGoPackage {
    pname = "server";
    version = "0.0.1";
    src = ./server;
    vendorHash = null;
    subPackages = [ "." ];
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
    pkgs.go
    pkgs.gcc
    systemdPkgs.systemd
  ];

  system.stateVersion = "24.11";

  # === SPECIFIC ===

  networking.hostName = "server";

  environment.variables = {
    CGO_ENABLED = "1";
  };

  networking = {
    firewall = {
      enable = true;
      allowedTCPPorts = [
        22
        8081
      ];
    };
  };

  users.users = {
    web = {
      isNormalUser = true;
      description = "Web server";
      password = "web";
    };
    admin = {
      isNormalUser = true;
      description = "Admin";
      password = "admin";
    };
  };

  # Start server as web user
  system.activationScripts.move_server = pkgs.lib.stringAfter [ "users" ] ''
    rm -fr /home/web/server /home/web/go
    cp /home/env-admin/nixos/hosts/server/server/db.sqlite /home/web/server/
    cp "${serverPkgs}/bin/main" /home/web/server
    chown -R web:users /home/web/server
  '';

  systemd.services = {
    server = {
      enable = true;
      wantedBy = [ "multi-user.target" ];
      description = "Web server";
      path = [ "/run/current-system/sw" ];
      script = ''
        cd /home/web/server &&
        ./main
      '';
      serviceConfig = {
        User = "web";
        Restart = "always";
      };
    };

    elann = {
      enable = true;
      wantedBy = [ "multi-user.target" ];
      description = "Web server";
      path = [ "/run/current-system/sw" ];
      script = ''
        while true
        do
          curl localhost:8081/jump/?name=elann
          sleep 1
        done
      '';
      serviceConfig = {
        User = "web";
        Restart = "always";
      };
    };

    jimmy = {
      enable = true;
      wantedBy = [ "multi-user.target" ];
      description = "Web server";
      path = [ "/run/current-system/sw" ];
      script = ''
        while true
        do
          curl localhost:8081/to-pocket/?name=jimmy &
          curl localhost:8081/to-pocket/?name=jimmy &
          curl localhost:8081/to-bank/?name=jimmy
        done
      '';
      serviceConfig = {
        User = "web";
        Restart = "always";
      };
    };
  };

  # ssh into web from path traversal to find public key
  # https://github.com/edwinlinson/ecom/blob/main/key.md.pub
  home-manager = {
    users."web" = {
      programs.home-manager.enable = true;
      home = {
        username = "web";
        homeDirectory = "/home/web";
        stateVersion = "24.11";
        file = {
          "public_key" = {
            text = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAILROvEP/drEaoyD6o2pXLTwUqq6+HXi//iKYl7hrXXxY web";
            target = ".ssh/authorized_keys";
          };
        };
      };
    };
  };

  # Elevation to admin
  services.cron = {
    enable = true;
    systemCronJobs = [ "*/5 * * * * web rsync -a /tmp/server/ /home/admin" ];
  };

  # Elevation to root to restart server CVE-2023-26604
  security.sudo.extraRules = [
    {
      users = [ "admin" ];
      commands = [
        {
          command = "/run/current-system/sw/bin/systemctl status";
          options = [
            "SETENV"
            "NOPASSWD"
          ];
        }
      ];
    }
  ];

}
