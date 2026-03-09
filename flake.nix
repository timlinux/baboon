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

        # Hugo documentation site
        docs = pkgs.stdenv.mkDerivation {
          pname = "baboon-docs";
          version = "1.0.0";
          src = ./hugo;

          nativeBuildInputs = [ pkgs.hugo ];

          buildPhase = ''
            hugo --minify
          '';

          installPhase = ''
            cp -r public $out
          '';
        };

        # Script to run Hugo dev server
        docs-serve = pkgs.writeShellScriptBin "baboon-docs-serve" ''
          cd ${toString ./hugo}
          ${pkgs.hugo}/bin/hugo server -D --bind 0.0.0.0 --port 1313
        '';

        # Script to build docs
        docs-build = pkgs.writeShellScriptBin "baboon-docs-build" ''
          cd ${toString ./hugo}
          ${pkgs.hugo}/bin/hugo --minify
          echo "Documentation built in hugo/public/"
        '';

        # Script to open docs in browser
        docs-open = pkgs.writeShellScriptBin "baboon-docs-open" ''
          ${pkgs.xdg-utils}/bin/xdg-open http://localhost:1313 2>/dev/null || \
          open http://localhost:1313 2>/dev/null || \
          echo "Open http://localhost:1313 in your browser"
        '';

      in
      {
        packages = {
          default = pkgs.buildGoModule {
            pname = "baboon";
            version = "1.4.0";
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

          # Documentation packages
          docs = docs;
          docs-serve = docs-serve;
          docs-build = docs-build;
          docs-open = docs-open;
        };

        # Apps for `nix run`
        apps = {
          default = {
            type = "app";
            program = "${self.packages.${system}.default}/bin/baboon";
          };

          # nix run .#docs-serve - Start Hugo dev server
          docs-serve = {
            type = "app";
            program = "${docs-serve}/bin/baboon-docs-serve";
          };

          # nix run .#docs-build - Build documentation
          docs-build = {
            type = "app";
            program = "${docs-build}/bin/baboon-docs-build";
          };

          # nix run .#docs-open - Open docs in browser
          docs-open = {
            type = "app";
            program = "${docs-open}/bin/baboon-docs-open";
          };
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            gotools
            go-tools
            hugo
            nodejs_20
          ];

          shellHook = ''
            echo "🐒 Baboon development environment"
            echo ""
            echo "Commands:"
            echo "  make build      - Build Go binary"
            echo "  make run        - Run the app"
            echo "  make docs-dev   - Start Hugo dev server"
            echo "  make docs-build - Build documentation"
            echo ""
            echo "Nix run commands:"
            echo "  nix run .#docs-serve - Start Hugo dev server"
            echo "  nix run .#docs-build - Build documentation"
            echo ""
          '';
        };
      }
    );
}
