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
    dig
    helix
    iptables
    busybox
  ];

  system.stateVersion = "24.11";

  # === SPECIFIC ===

  networking.hostName = "c2";

  networking = {
    firewall = {
      enable = true;
      allowedTCPPorts = [
        22
        23
        80
        443
        53
        8053
      ];
      allowedUDPPorts = [
        53
        8053
      ];

      # IPTables to redirect DNS request to port 8053 of the resolver
      extraCommands = ''
        iptables -F
        iptables -t nat -A OUTPUT -p tcp --dport domain -j DNAT --to-destination ${config.resolver}
        iptables -t nat -A OUTPUT -p udp --dport domain -j DNAT --to-destination ${config.resolver}
      '';
    };

    # not ideal, seems to only generate first time
    resolvconf.enable = false;
    # networkmanager.insertNameservers = [ "127.0.0.1" ];
  };

  users.users.admin = {
    isNormalUser = true;
    description = "admin of the server sending message to satellite";
    extraGroups = [ "wheel" ];
    createHome = true;
    password = "admin";
  };

  systemd.services.messages = {
    enable = true;
    wantedBy = [ "multi-user.target" ];
    description = "Web server";
    path = [ "/run/current-system/sw" ];
    script = config.sh_messages;
    serviceConfig = {
      User = "admin";
    };
  };

  systemd.services.telnet = {
    enable = true;
    wantedBy = [ "multi-user.target" ];
    description = "Telnet server";
    path = [ "/run/current-system/sw" ];
    script = ''telnetd -l /run/current-system/sw/bin/login -F'';
    serviceConfig = {
      User = "root";
    };
  };

  home-manager = {
    users."admin" = {
      programs.home-manager.enable = true;
      home = {
        username = "admin";
        homeDirectory = "/home/admin";
        stateVersion = "24.11";
        file."test.sh" = {
          text = config.sh_messages;
        };
      };
    };
  };
}
