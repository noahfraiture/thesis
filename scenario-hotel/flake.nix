{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-24.11";
  };

  outputs =
    { self, nixpkgs }:
    let
      pkgs = import nixpkgs {
        system = "x86_64-linux";
        config.allowUnfree = true;
      };
    in
    {
      devShells."x86_64-linux".default = pkgs.mkShell {

        buildInputs = with pkgs; [
        ];

        packages =
          with pkgs;
          [
            gcc
            nmap
            python312
          ]
          ++ (with pkgs.python312Packages; [
            requests
            pycryptodome
          ]);

        DIRENV = "env";
      };
    };
}
