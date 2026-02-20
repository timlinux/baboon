{
  description = "Baboon - A terminal typing practice app";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages = {
          default = pkgs.buildGoModule {
            pname = "baboon";
            version = "1.1.0";
            src = ./.;
            vendorHash = null;

            meta = with pkgs.lib; {
              description = "A terminal typing practice app with ASCII art";
              homepage = "https://github.com/timlinux/baboon";
              license = licenses.mit;
              maintainers = [ ];
              mainProgram = "baboon";
            };
          };
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            gotools
            go-tools
          ];
        };
      }
    );
}
