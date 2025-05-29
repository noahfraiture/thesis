{
  description = "A NixOS Docker image with custom configuration";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.11";
  };

  outputs =
    { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };

      # Helper function to build Docker image
      buildDockerImage =
        {
          name,
          config,
          packages,
        }:
        let
          nixos = import "${nixpkgs}/nixos" {
            configuration = config;
            system = system;
          };
          systemClosure = nixos.config.system.build.toplevel;
        in
        pkgs.dockerTools.buildLayeredImage {
          inherit name;
          tag = "latest";
          contents = [
            systemClosure
          ] ++ packages;
          config = {
            Cmd = [ "${systemClosure}/init" ];
          };
          extraCommands = ''
            ${pkgs.lib.concatStringsSep "\n" (
              map (pkg: ''
                cp -r ${pkg}/* tmp/
                chmod -R +w tmp
              '') packages
            )}
          '';
        };
    in
    {
      packages.${system} = {
        zheng-victim = buildDockerImage {
          name = "zheng-victim";
          config = import ./victim.nix { inherit pkgs; };
          packages = [
            (pkgs.buildGoModule {
              pname = "victim";
              version = "0.0.1";
              src = ../zheng/victim;
              vendorHash = null;
              subPackages = [ "." ];
            })
          ];
        };
        zheng-c2 = buildDockerImage {
          name = "zheng-c2";
          config = import ./c2.nix { inherit pkgs; };
          packages = [
            (pkgs.buildGoModule {
              pname = "c2";
              version = "0.0.1";
              src = ../zheng/c2;
              vendorHash = null;
              subPackages = [ "." ];
            })
          ];
        };
        zheng-host = buildDockerImage {
          name = "zheng-host";
          config = import ./host.nix { inherit pkgs; };
          packages = [ ];
        };
      };
    };
}
