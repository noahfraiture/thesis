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

  users.users = {
    env-admin = {
      isNormalUser = true;
      password = "env-admin";
      extraGroups = [ "wheel" ];
    };
    admin = {
      isNormalUser = true;
      password = "admin";
      extraGroups = [ "wheel" ];
    };
  };

  environment.systemPackages = with pkgs; [ bash ];

  networking.firewall.enable = false;
}
