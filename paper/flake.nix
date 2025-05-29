{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  };

  outputs =
    { self, nixpkgs }:
    let
      pkgs = import nixpkgs {
        system = "x86_64-linux";
        config.allowUnfree = true;
      };

      pdfProcessScript = pkgs.writeShellScriptBin "pdf-process" ''
        #!/usr/bin/env bash
        COVER="./latex/LaTeX/EPL-master-thesis-covers-template-EN.pdf"
        MAIN="./main.pdf"
        set -e
        (cd ./latex/LaTeX && ${pkgs.texliveFull}/bin/latexmk -pdf EPL-master-thesis-covers-template-EN.tex)
        typst compile ./main.typ
        ${pkgs.pdftk}/bin/pdftk $COVER cat 1 output front.pdf
        ${pkgs.pdftk}/bin/pdftk $COVER cat 2 output back.pdf
        ${pkgs.pdftk}/bin/pdftk front.pdf $MAIN back.pdf cat output Fraiture_16702000_2025.pdf
        rm front.pdf back.pdf main.pdf ./latex/LaTeX/EPL-master-thesis-covers-template-EN.pdf
      '';
    in
    {
      devShells."x86_64-linux".default = pkgs.mkShell {

        buildInputs = with pkgs; [
        ];

        packages = with pkgs; [
          typst
          texliveFull
          pdftk
        ];

        DIRENV = "typst";
      };

      apps."x86_64-linux".pdf-process = {
        type = "app";
        program = "${pdfProcessScript}/bin/pdf-process";
      };
    };
}
