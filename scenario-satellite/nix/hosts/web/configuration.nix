{ config, pkgs, ... }:
{
  imports = [
    ../../config.nix
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

  environment.systemPackages = with pkgs; [
    helix
    dig
  ];

  system.stateVersion = "24.11";

  # === SPECIFIC ===

  networking.hostName = "web";

  # TODO : enable and change rules in other vms
  networking = {
    firewall = {
      enable = true;
      allowedTCPPorts = [ 80 ];

      # FIXME
      # Use ping sweep ? Or other discovery method
      extraCommands = ''
        iptables -F
        iptables -t nat -A OUTPUT -p tcp --dport domain -j DNAT --to-destination ${config.resolver}
        iptables -t nat -A OUTPUT -p udp --dport domain -j DNAT --to-destination ${config.resolver}
      '';
    };

    resolvconf.enable = false;
  };

  users.users = {
    web = {
      isNormalUser = true;
      description = "web user";
      createHome = true;
      packages = with pkgs; [
        php
        vim
        python3
      ];
      password = "web"; # TODO : replace with hashed password
    };

    admin = {
      isNormalUser = true;
      description = "admin of the web server";
      extraGroups = [
        "wheel"
      ];
      createHome = true;
      password = "admin";
    };
  };

  # TODO : might change to use home.file to be cleaner
  system.activationScripts.move_web = pkgs.lib.stringAfter [ "users" ] ''
    rm -rf /home/web/visitor /home/web/admin
    mv /home/env-admin/nixos/hosts/web/visitor /home/web/
    mv /home/env-admin/nixos/hosts/web/admin /home/web/
    chown -R web:users /home/web/visitor
    chown -R web:users /home/web/admin
  '';

  systemd.services.web = {
    enable = true;
    wantedBy = [ "multi-user.target" ];
    path = with pkgs; [
      php
    ];
    description = "Web server";
    script = ''
      (cd /home/web/visitor && php -S 0.0.0.0:8000) &
      (cd /home/web/admin && php -S 0.0.0.0:8001)
    '';
    serviceConfig = {
      User = "web";
    };
  };

  services.traefik = {
    enable = true;

    staticConfigOptions = {
      entryPoints = {
        web = {
          address = ":80";
          asDefault = true;
        };
      };
    };

    dynamicConfigOptions = {
      # Visitor page
      http.routers.visitor = {
        rule = "Host(`stellar.com`)";
        service = "visitor";
      };
      http.services.visitor.loadBalancer.servers = [ { url = "http://localhost:8000"; } ];

      http.routers.admin = {
        rule = "Host(`admin.com`)";
        service = "admin";
      };
      http.services.admin.loadBalancer.servers = [ { url = "http://localhost:8001"; } ];
    };
  };

  networking.extraHosts = ''
    127.0.0.1 stellar.com
    127.0.0.1 admin.com
  '';

  security.sudo.extraRules = [
    {
      users = [ "web" ];
      commands = [
        {
          command = "/etc/profiles/per-user/web/bin/vim";
          options = [
            "SETENV"
            "NOPASSWD"
          ];
        }
      ];
    }
  ];

  home-manager = {
    users."web" = {
      programs.home-manager.enable = true;
      home = {
        username = "web";
        homeDirectory = "/home/web";
        stateVersion = "24.11";
      };
    };

    users."admin" = {
      programs.home-manager.enable = true;
      home = {
        username = "admin";
        homeDirectory = "/home/admin";
        stateVersion = "24.11";
      };

      programs.ssh = {
        enable = true;
        addKeysToAgent = "yes";
      };

      home.file = {
        "private_key" = {
          text = config.web_private;
          target = ".ssh/resolver";
        };
        "public_key" = {
          text = config.web_public;
          target = ".ssh/resolver.pub";
        };
        "note.txt" = {
          text = "Don't forget to update the resolver at ${config.resolver}, use the user 'admin'";
        };
      };
    };
  };
}
