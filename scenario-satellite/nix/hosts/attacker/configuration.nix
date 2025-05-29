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
    iptables
  ];

  system.stateVersion = "24.11";

  # === SPECIFIC ===

  services = {
    xserver = {
      enable = true;
      desktopManager = {
        xterm.enable = true;
        xfce.enable = true;
      };
    };
    displayManager.defaultSession = "xfce";
  };

  networking.hostName = "nixos";

  networking = {

    firewall = {
      enable = true;
      allowedTCPPorts = [
        22
        80
        443
        6666
      ];
      extraCommands = ''
        iptables -F
        iptables -t nat -A OUTPUT -p tcp --dport 80 -j DNAT --to-destination ${config.web}
      '';
    };
  };

  virtualisation.docker = {
    enable = true;
  };

  users.users.attacker = {
    isNormalUser = true;
    description = "Attacker";
    createHome = true;
    password = "attacker";
    extraGroups = [
      "wheel"
      "docker"
    ];
    packages = with pkgs; [
      firefox-unwrapped
    ];
  };
}
