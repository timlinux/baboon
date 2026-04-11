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

        # Script to record demo with asciinema
        demo-record = pkgs.writeShellScriptBin "baboon-demo-record" ''
          #!/usr/bin/env bash
          set -e

          # Use current working directory (should be project root)
          PROJECT_DIR="$(pwd)"
          DEMO_DIR="$PROJECT_DIR/demo"
          CAST_FILE="$DEMO_DIR/baboon-demo.cast"
          GIF_FILE="$DEMO_DIR/baboon-demo.gif"

          # Verify we're in the right directory
          if [ ! -f "$PROJECT_DIR/flake.nix" ]; then
            echo "❌ Please run this command from the baboon project root directory."
            exit 1
          fi

          mkdir -p "$DEMO_DIR"

          echo "🎬 Recording Baboon demo..."
          echo ""
          echo "Tips for a good demo:"
          echo "  - Start baboon and begin typing"
          echo "  - Show the real-time colour feedback"
          echo "  - Complete a few words to show the flow"
          echo "  - Optionally complete a round to show statistics"
          echo "  - Keep it under 30 seconds"
          echo ""
          echo "Press Enter to start recording, type 'exit' when done..."
          read -r

          ${pkgs.asciinema}/bin/asciinema rec --overwrite "$CAST_FILE"

          echo ""
          echo "✅ Recording saved to $CAST_FILE"
          echo ""
          echo "Converting to GIF for README..."

          ${pkgs.asciinema-agg}/bin/agg --theme monokai "$CAST_FILE" "$GIF_FILE"
          echo "✅ GIF saved to $GIF_FILE"

          # Copy GIF to hugo static folder
          mkdir -p "$PROJECT_DIR/hugo/static/images"
          cp "$GIF_FILE" "$PROJECT_DIR/hugo/static/images/baboon-demo.gif"
          echo "✅ GIF copied to hugo/static/images/"

          echo ""
          echo "🎉 Demo recording complete!"
          echo ""
          echo "The demo is now available at:"
          echo "  - demo/baboon-demo.cast (asciinema format)"
          echo "  - demo/baboon-demo.gif (animated GIF)"
          echo "  - hugo/static/images/baboon-demo.gif (for docs)"
          echo ""
          echo "README.md and docs will automatically use the new demo."
        '';

        # Script to play demo locally
        demo-play = pkgs.writeShellScriptBin "baboon-demo-play" ''
          #!/usr/bin/env bash
          PROJECT_DIR="$(pwd)"
          CAST_FILE="$PROJECT_DIR/demo/baboon-demo.cast"

          if [ ! -f "$CAST_FILE" ]; then
            echo "❌ No demo recording found at $CAST_FILE"
            echo "Run 'nix run .#demo-record' to create one."
            exit 1
          fi

          echo "🎬 Playing Baboon demo..."
          ${pkgs.asciinema}/bin/asciinema play "$CAST_FILE"
        '';

        # Script to manage releases
        release = pkgs.writeShellScriptBin "baboon-release" ''
          #!/usr/bin/env bash
          PROJECT_DIR="$(pwd)"
          if [ ! -f "$PROJECT_DIR/scripts/release.sh" ]; then
            echo "❌ Please run this command from the baboon project root directory."
            exit 1
          fi
          exec "$PROJECT_DIR/scripts/release.sh"
        '';

        # Script to install web dependencies
        web-install = pkgs.writeShellScriptBin "baboon-web-install" ''
          #!/usr/bin/env bash
          PROJECT_DIR="$(pwd)"
          if [ ! -f "$PROJECT_DIR/web/package.json" ]; then
            echo "❌ Please run this command from the baboon project root directory."
            exit 1
          fi
          cd "$PROJECT_DIR/web"
          ${pkgs.nodejs_22}/bin/npm install
        '';

        # Script to build web frontend
        web-build = pkgs.writeShellScriptBin "baboon-web-build" ''
          #!/usr/bin/env bash
          PROJECT_DIR="$(pwd)"
          if [ ! -f "$PROJECT_DIR/web/package.json" ]; then
            echo "❌ Please run this command from the baboon project root directory."
            exit 1
          fi
          cd "$PROJECT_DIR/web"
          ${pkgs.nodejs_22}/bin/npm install
          ${pkgs.nodejs_22}/bin/npm run build
          # Copy build to dist for server
          mkdir -p dist
          cp -r build/* dist/
          echo "Web frontend built in web/dist/"
        '';

        # Script to run web dev server
        web-dev = pkgs.writeShellScriptBin "baboon-web-dev" ''
          #!/usr/bin/env bash
          PROJECT_DIR="$(pwd)"
          if [ ! -f "$PROJECT_DIR/web/package.json" ]; then
            echo "❌ Please run this command from the baboon project root directory."
            exit 1
          fi
          cd "$PROJECT_DIR/web"
          ${pkgs.nodejs_22}/bin/npm install
          ${pkgs.nodejs_22}/bin/npm run dev
        '';

      in
      {
        packages = {
          default = pkgs.buildGoModule {
            pname = "baboon";
            version = "1.7.0";
            src = ./.;
            vendorHash = null; # Will be updated after first build

            # Required for go-sqlite3 CGO driver
            nativeBuildInputs = [ pkgs.pkg-config ];
            buildInputs = [ pkgs.sqlite ];
            CGO_ENABLED = 1;

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
          demo-record = demo-record;
          demo-play = demo-play;
          release = release;
          # Web frontend packages
          web-install = web-install;
          web-build = web-build;
          web-dev = web-dev;
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

          # nix run .#demo-record - Record a demo with asciinema
          demo-record = {
            type = "app";
            program = "${demo-record}/bin/baboon-demo-record";
          };

          # nix run .#demo-play - Play the demo locally
          demo-play = {
            type = "app";
            program = "${demo-play}/bin/baboon-demo-play";
          };

          # nix run .#release - Version bump and release
          release = {
            type = "app";
            program = "${release}/bin/baboon-release";
          };

          # nix run .#web-install - Install web dependencies
          web-install = {
            type = "app";
            program = "${web-install}/bin/baboon-web-install";
          };

          # nix run .#web-build - Build web frontend
          web-build = {
            type = "app";
            program = "${web-build}/bin/baboon-web-build";
          };

          # nix run .#web-dev - Start web dev server
          web-dev = {
            type = "app";
            program = "${web-dev}/bin/baboon-web-dev";
          };
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            gotools
            go-tools
            hugo
            nodejs_22
            asciinema
            asciinema-agg
            # Required for CGo (SQLite driver)
            gcc
            pkg-config
            sqlite
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
            echo "  nix run .#web-build   - Build web frontend"
            echo "  nix run .#web-dev     - Start web dev server"
            echo "  nix run .#docs-serve  - Start Hugo dev server"
            echo "  nix run .#docs-build  - Build documentation"
            echo "  nix run .#demo-record - Record demo with asciinema"
            echo "  nix run .#demo-play   - Play recorded demo"
            echo "  nix run .#release     - Version bump and release"
            echo ""
          '';
        };
      }
    );
}
