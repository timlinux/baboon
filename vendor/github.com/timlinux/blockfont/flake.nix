{
  description = "blockfont - Unicode block letter rendering for terminal applications";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        version = "0.1.0";

        # Script to record demo with asciinema
        demo-record = pkgs.writeShellScriptBin "blockfont-demo-record" ''
          #!/usr/bin/env bash
          set -e

          # Use current working directory (should be project root)
          PROJECT_DIR="$(pwd)"
          DEMO_DIR="$PROJECT_DIR/demo"
          CAST_FILE="$DEMO_DIR/blockfont-demo.cast"
          GIF_FILE="$DEMO_DIR/blockfont-demo.gif"

          # Verify we're in the right directory
          if [ ! -f "$PROJECT_DIR/flake.nix" ]; then
            echo "Please run this command from the blockfont project root directory."
            exit 1
          fi

          mkdir -p "$DEMO_DIR"

          echo "Recording blockfont demo..."
          echo ""
          echo "Tips for a good demo:"
          echo "  - Run 'make demo' to start the demo application"
          echo "  - Navigate through different screens (1-8)"
          echo "  - Show vim editing on the editing screen"
          echo "  - Show animations on the animation screen"
          echo "  - Keep it under 60 seconds"
          echo ""
          echo "Press Enter to start recording, type 'exit' when done..."
          read -r

          ${pkgs.asciinema}/bin/asciinema rec --overwrite "$CAST_FILE"

          echo ""
          echo "Recording saved to $CAST_FILE"
          echo ""
          echo "Converting to GIF for README..."

          ${pkgs.asciinema-agg}/bin/agg --theme monokai "$CAST_FILE" "$GIF_FILE"
          echo "GIF saved to $GIF_FILE"

          # Copy GIF to docs static folder
          mkdir -p "$PROJECT_DIR/docs/docs/images"
          cp "$GIF_FILE" "$PROJECT_DIR/docs/docs/images/blockfont-demo.gif"
          echo "GIF copied to docs/docs/images/"

          echo ""
          echo "Demo recording complete!"
          echo ""
          echo "The demo is now available at:"
          echo "  - demo/blockfont-demo.cast (asciinema format)"
          echo "  - demo/blockfont-demo.gif (animated GIF)"
          echo "  - docs/docs/images/blockfont-demo.gif (for docs)"
          echo ""
          echo "README.md and docs will automatically use the new demo."
        '';

        # Script to play demo locally
        demo-play = pkgs.writeShellScriptBin "blockfont-demo-play" ''
          #!/usr/bin/env bash
          PROJECT_DIR="$(pwd)"
          CAST_FILE="$PROJECT_DIR/demo/blockfont-demo.cast"

          if [ ! -f "$CAST_FILE" ]; then
            echo "No demo recording found at $CAST_FILE"
            echo "Run 'nix run .#demo-record' to create one."
            exit 1
          fi

          echo "Playing blockfont demo..."
          ${pkgs.asciinema}/bin/asciinema play "$CAST_FILE"
        '';

        # Script to manage releases
        release = pkgs.writeShellScriptBin "blockfont-release" ''
          #!/usr/bin/env bash
          set -e

          PROJECT_DIR="$(pwd)"

          # Verify we're in the right directory
          if [ ! -f "$PROJECT_DIR/flake.nix" ]; then
            echo "Please run this command from the blockfont project root directory."
            exit 1
          fi

          # Get current version from flake.nix
          CURRENT_VERSION=$(grep 'version = "' "$PROJECT_DIR/flake.nix" | head -1 | sed 's/.*version = "\([^"]*\)".*/\1/')
          echo "Current version: $CURRENT_VERSION"
          echo ""

          # Parse version components
          MAJOR=$(echo "$CURRENT_VERSION" | cut -d. -f1)
          MINOR=$(echo "$CURRENT_VERSION" | cut -d. -f2)
          PATCH=$(echo "$CURRENT_VERSION" | cut -d. -f3)

          echo "Select version bump type:"
          echo "  1) Patch ($MAJOR.$MINOR.$((PATCH + 1))) - Bug fixes"
          echo "  2) Minor ($MAJOR.$((MINOR + 1)).0) - New features"
          echo "  3) Major ($((MAJOR + 1)).0.0) - Breaking changes"
          echo "  4) Custom version"
          echo ""
          read -p "Choice [1-4]: " choice

          case $choice in
            1) NEW_VERSION="$MAJOR.$MINOR.$((PATCH + 1))" ;;
            2) NEW_VERSION="$MAJOR.$((MINOR + 1)).0" ;;
            3) NEW_VERSION="$((MAJOR + 1)).0.0" ;;
            4) read -p "Enter version: " NEW_VERSION ;;
            *) echo "Invalid choice"; exit 1 ;;
          esac

          echo ""
          echo "New version will be: $NEW_VERSION"
          read -p "Continue? [y/N] " confirm
          if [ "$confirm" != "y" ] && [ "$confirm" != "Y" ]; then
            echo "Aborted."
            exit 0
          fi

          # Update version in flake.nix
          sed -i "s/version = \"$CURRENT_VERSION\"/version = \"$NEW_VERSION\"/" "$PROJECT_DIR/flake.nix"

          # Commit and tag
          ${pkgs.git}/bin/git add -A
          ${pkgs.git}/bin/git commit -m "Release v$NEW_VERSION"
          ${pkgs.git}/bin/git tag -a "v$NEW_VERSION" -m "Release v$NEW_VERSION"

          echo ""
          echo "Release v$NEW_VERSION prepared!"
          echo ""
          echo "To publish:"
          echo "  git push origin main --tags"
          echo ""
          echo "Then create a GitHub release at:"
          echo "  https://github.com/timlinux/blockfont/releases/new?tag=v$NEW_VERSION"
        '';

        # Script to serve docs
        docs-serve = pkgs.writeShellScriptBin "blockfont-docs-serve" ''
          cd ${toString ./.}/docs
          ${pkgs.python312Packages.mkdocs}/bin/mkdocs serve
        '';

        # Script to build docs
        docs-build = pkgs.writeShellScriptBin "blockfont-docs-build" ''
          cd ${toString ./.}/docs
          ${pkgs.python312Packages.mkdocs}/bin/mkdocs build
          echo "Documentation built in docs/site/"
        '';

        # Script to open docs in browser
        docs-open = pkgs.writeShellScriptBin "blockfont-docs-open" ''
          ${pkgs.xdg-utils}/bin/xdg-open http://localhost:8000 2>/dev/null || \
          open http://localhost:8000 2>/dev/null || \
          echo "Open http://localhost:8000 in your browser"
        '';

      in
      {
        packages = {
          default = pkgs.buildGoModule {
            pname = "blockfont";
            inherit version;
            src = ./.;
            vendorHash = null;

            meta = with pkgs.lib; {
              description = "Unicode block letter rendering for terminal applications";
              homepage = "https://github.com/timlinux/blockfont";
              license = licenses.mit;
              maintainers = [ ];
            };
          };

          # Demo binary
          demo = pkgs.buildGoModule {
            pname = "blockfont-demo";
            inherit version;
            src = ./.;
            vendorHash = null;
            subPackages = [ "examples/demo" ];
          };

          # Documentation and scripts
          docs-serve = docs-serve;
          docs-build = docs-build;
          docs-open = docs-open;
          demo-record = demo-record;
          demo-play = demo-play;
          release = release;
        };

        # Apps for `nix run`
        apps = {
          default = {
            type = "app";
            program = "${self.packages.${system}.demo}/bin/demo";
          };

          # nix run .#demo - Run the demo application
          demo = {
            type = "app";
            program = "${self.packages.${system}.demo}/bin/demo";
          };

          # nix run .#docs-serve - Start MkDocs dev server
          docs-serve = {
            type = "app";
            program = "${docs-serve}/bin/blockfont-docs-serve";
          };

          # nix run .#docs-build - Build documentation
          docs-build = {
            type = "app";
            program = "${docs-build}/bin/blockfont-docs-build";
          };

          # nix run .#docs-open - Open docs in browser
          docs-open = {
            type = "app";
            program = "${docs-open}/bin/blockfont-docs-open";
          };

          # nix run .#demo-record - Record a demo with asciinema
          demo-record = {
            type = "app";
            program = "${demo-record}/bin/blockfont-demo-record";
          };

          # nix run .#demo-play - Play the demo locally
          demo-play = {
            type = "app";
            program = "${demo-play}/bin/blockfont-demo-play";
          };

          # nix run .#release - Version bump and release
          release = {
            type = "app";
            program = "${release}/bin/blockfont-release";
          };

          # nix run .#build - Build the library
          build = {
            type = "app";
            program = toString (pkgs.writeShellScript "build" ''
              cd ${toString ./.}
              ${pkgs.go}/bin/go build ./...
            '');
          };

          # nix run .#test - Run tests
          test = {
            type = "app";
            program = toString (pkgs.writeShellScript "test" ''
              cd ${toString ./.}
              ${pkgs.go}/bin/go test -v ./...
            '');
          };

          # nix run .#lint - Run linters
          lint = {
            type = "app";
            program = toString (pkgs.writeShellScript "lint" ''
              cd ${toString ./.}
              ${pkgs.golangci-lint}/bin/golangci-lint run ./...
            '');
          };
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            # Go development
            go
            gopls
            golangci-lint
            gotools
            go-tools
            delve

            # Pre-commit and linting
            pre-commit

            # Documentation
            python312
            python312Packages.mkdocs
            python312Packages.mkdocs-material
            python312Packages.mkdocs-minify-plugin

            # Demo recording
            asciinema
            asciinema-agg

            # Git and utilities
            git
            gh
            xdg-utils
          ];

          shellHook = ''
            echo "blockfont development environment"
            echo "Go version: $(go version)"
            echo ""
            echo "Make commands:"
            echo "  make build         - Build/verify the library"
            echo "  make test          - Run tests"
            echo "  make demo          - Run the interactive demo"
            echo "  make docs-dev      - Start MkDocs dev server"
            echo "  make lint          - Lint code"
            echo ""
            echo "Nix run commands:"
            echo "  nix run .#demo        - Run demo application"
            echo "  nix run .#docs-serve  - Start MkDocs dev server"
            echo "  nix run .#demo-record - Record demo with asciinema"
            echo "  nix run .#demo-play   - Play recorded demo"
            echo "  nix run .#release     - Version bump and release"
            echo ""
          '';
        };
      }
    );
}
