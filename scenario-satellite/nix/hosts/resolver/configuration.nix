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
      openssh.authorizedKeys.keys = [ ];
      password = "env-admin";
    };
  };

  nixpkgs.config.allowUnfree = true;

  environment.systemPackages = with pkgs; [
    helix
    tcpdump
    dig
  ];

  system.stateVersion = "24.11";

  # === SPECIFIC ===

  networking.hostName = "resolver";

  networking.firewall = {
    enable = false;
    allowedTCPPorts = [ 53 ];
    allowedUDPPorts = [ 53 ];
  };

  services.cron = {
    enable = true;
    systemCronJobs = [
      "*/5 * * * *      root    /tmp/log.sh >> /tmp/cron.log"
    ];
  };

  services.unbound = {
    enable = true;
    settings = {
      server = {
        interface = [ "0.0.0.0" ];
        port = 53;
        access-control = [ "0.0.0.0/0 allow" ];

        # Cache settings
        cache-max-ttl = 2; # 24 hours
        cache-min-ttl = 1; # 1 hour
      };

      forward-zone = [
        {
          name = ".";
          forward-addr = [ config.autoritary ];
        }
      ];
    };
  };

  users.users = {
    admin = {
      isNormalUser = true;
      description = "resolver admin";
      extraGroups = [
        "networkmanager"
        "wheel"
      ];
      password = "admin";
      openssh.authorizedKeys.keys = [ config.web_public ];
    };
  };
}
