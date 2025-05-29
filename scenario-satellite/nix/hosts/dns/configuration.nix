{ config, pkgs, ... }:
{
  imports = [
    ../../config.nix
    ../hardware-configuration.nix
  ];

  nix.settings = {
    experimental-features = "nix-command flakes";
    auto-optimise-store = true;
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
    tcpdump
  ];

  system.stateVersion = "24.11";

  # === Specific configuration ===

  networking.hostName = "dns";

  networking.firewall = {
    enable = false;
    allowedTCPPorts = [ 53 ];
    allowedUDPPorts = [ 53 ];
  };

  # https://calomel.org/unbound_dns.html
  # https://search.nixos.org/options?channel=unstable&show=services.unbound.settings&from=0&size=50&sort=alpha_asc&type=packages&query=services.unbound
  # https://wiki.nixos.org/wiki/Unbound
  # https://github.com/NLnetLabs/unbound/blob/master/doc/example.conf.in
  services.unbound = {
    enable = true;
    settings = {
      server = {
        interface = [ "0.0.0.0" ];
        port = 53;
        access-control = [ "0.0.0.0/0 allow" ];

        local-zone = [
          ''"satellite.app" static''
        ];
        local-data = [
          "\"satellite.app A ${config.satellite}\""
        ];
      };

      forward-zone = [
        {
          name = ".";
          forward-addr = [
            "1.0.0.1@53#one.one.one.one"
            "1.1.1.1@53#one.one.one.one"
            "8.8.4.4@53#dns.google"
            "8.8.8.8@53#dns.google"
            "9.9.9.9@53#dns.quad9.net"
            "149.112.112.112@53#dns.quad9.net"
          ];
        }
      ];
    };
  };
}
