{ pkgs }:
{
  system.stateVersion = "24.11";
  boot.isContainer = true;
  services.openssh = {
    enable = true;
    settings = {
      PasswordAuthentication = true;
    };
  };

  users.users.env-admin = {
    isNormalUser = true;
    password = "env-admin";
    extraGroups = [ "wheel" ];
  };

  environment.systemPackages = with pkgs; [ bash ];

  networking.firewall.enable = false;

  systemd.services.c2 = {
    enable = true;
    wantedBy = [ "multi-user.target" ];
    path = [ "/run/current-system/sw" ];
    script = "/bin/main";
    serviceConfig = {
      User = "env-admin";
    };
  };
}
