#!/usr/bin/env bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Get project root (where flake.nix lives)
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
FLAKE_FILE="$PROJECT_DIR/flake.nix"

# Verify we're in the right directory
if [ ! -f "$FLAKE_FILE" ]; then
    echo -e "${RED}❌ Cannot find flake.nix at $FLAKE_FILE${NC}"
    echo "Please run this script from the baboon project."
    exit 1
fi

cd "$PROJECT_DIR"

# Extract current version from flake.nix (the one after pname = "baboon")
CURRENT_VERSION=$(grep -A1 'pname = "baboon";' "$FLAKE_FILE" | grep 'version' | sed 's/.*version = "\([^"]*\)".*/\1/')

# Get latest git tag
LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "none")

# Get git status
UNCOMMITTED=$(git status --porcelain 2>/dev/null | grep -v "^??" | wc -l | tr -d ' ')

echo ""
echo -e "${BLUE}╔════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║${NC}           ${GREEN}🐒 Baboon Release Manager${NC}                     ${BLUE}║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${YELLOW}Current State:${NC}"
echo -e "  📦 Version in flake.nix: ${GREEN}$CURRENT_VERSION${NC}"
echo -e "  🏷️  Latest git tag:       ${GREEN}$LATEST_TAG${NC}"
echo -e "  📝 Uncommitted changes:   ${UNCOMMITTED}"
echo ""

# Check for uncommitted changes
if [ "$UNCOMMITTED" -gt 0 ]; then
    echo -e "${YELLOW}⚠️  You have uncommitted changes:${NC}"
    git status --short | grep -v "^??"
    echo ""
    read -p "Continue anyway? (y/N): " CONTINUE
    if [[ ! "$CONTINUE" =~ ^[Yy]$ ]]; then
        echo -e "${RED}Aborted.${NC}"
        exit 1
    fi
fi

# Suggest next version
IFS='.' read -ra VERSION_PARTS <<< "$CURRENT_VERSION"
MAJOR=${VERSION_PARTS[0]}
MINOR=${VERSION_PARTS[1]}
PATCH=${VERSION_PARTS[2]}

SUGGESTED_PATCH="$MAJOR.$MINOR.$((PATCH + 1))"
SUGGESTED_MINOR="$MAJOR.$((MINOR + 1)).0"
SUGGESTED_MAJOR="$((MAJOR + 1)).0.0"

echo -e "${YELLOW}Version Options:${NC}"
echo -e "  [1] Patch (bugfix):    ${GREEN}$SUGGESTED_PATCH${NC}"
echo -e "  [2] Minor (feature):   ${GREEN}$SUGGESTED_MINOR${NC}"
echo -e "  [3] Major (breaking):  ${GREEN}$SUGGESTED_MAJOR${NC}"
echo -e "  [4] Custom version"
echo ""
read -p "Select option (1-4) or enter version directly: " VERSION_CHOICE

case "$VERSION_CHOICE" in
    1)
        NEW_VERSION="$SUGGESTED_PATCH"
        ;;
    2)
        NEW_VERSION="$SUGGESTED_MINOR"
        ;;
    3)
        NEW_VERSION="$SUGGESTED_MAJOR"
        ;;
    4)
        read -p "Enter custom version (e.g., 2.0.0): " NEW_VERSION
        ;;
    *)
        # Check if they entered a version directly
        if [[ "$VERSION_CHOICE" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
            NEW_VERSION="$VERSION_CHOICE"
        else
            echo -e "${RED}❌ Invalid option: $VERSION_CHOICE${NC}"
            exit 1
        fi
        ;;
esac

# Validate version format
if [[ ! "$NEW_VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo -e "${RED}❌ Invalid version format: $NEW_VERSION${NC}"
    echo "Version must be in format: X.Y.Z (e.g., 1.4.0)"
    exit 1
fi

# Check if tag already exists
if git rev-parse "v$NEW_VERSION" >/dev/null 2>&1; then
    echo -e "${RED}❌ Tag v$NEW_VERSION already exists!${NC}"
    exit 1
fi

echo ""
echo -e "${YELLOW}Summary:${NC}"
echo -e "  Old version: ${RED}$CURRENT_VERSION${NC}"
echo -e "  New version: ${GREEN}$NEW_VERSION${NC}"
echo -e "  Git tag:     ${GREEN}v$NEW_VERSION${NC}"
echo ""
read -p "Proceed with release? (y/N): " CONFIRM
if [[ ! "$CONFIRM" =~ ^[Yy]$ ]]; then
    echo -e "${RED}Aborted.${NC}"
    exit 1
fi

# Update flake.nix (only the baboon package version, not docs version)
echo ""
echo -e "${BLUE}📝 Updating flake.nix...${NC}"
# Use awk to only replace the version after pname = "baboon"
awk -v old="$CURRENT_VERSION" -v new="$NEW_VERSION" '
    /pname = "baboon";/ { found=1 }
    found && /version = "/ { sub(/version = "[^"]*";/, "version = \"" new "\";"); found=0 }
    { print }
' "$FLAKE_FILE" > "$FLAKE_FILE.tmp" && mv "$FLAKE_FILE.tmp" "$FLAKE_FILE"
echo -e "${GREEN}✓ Updated version in flake.nix${NC}"

# Commit the change
echo -e "${BLUE}📦 Creating commit...${NC}"
git add "$FLAKE_FILE"
git commit -m "Bump version to $NEW_VERSION"
echo -e "${GREEN}✓ Committed version bump${NC}"

# Create tag
echo -e "${BLUE}🏷️  Creating tag v$NEW_VERSION...${NC}"
git tag "v$NEW_VERSION"
echo -e "${GREEN}✓ Created tag v$NEW_VERSION${NC}"

echo ""
echo -e "${GREEN}╔════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║${NC}              ${GREEN}✅ Release Prepared!${NC}                       ${GREEN}║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${YELLOW}To publish the release, run:${NC}"
echo ""
echo -e "  ${BLUE}git push origin main${NC}"
echo -e "  ${BLUE}git push origin v$NEW_VERSION${NC}"
echo ""
echo -e "Or push both at once:"
echo -e "  ${BLUE}git push origin main --tags${NC}"
echo ""
echo -e "${YELLOW}This will trigger GitHub Actions to:${NC}"
echo "  • Build binaries for Linux, macOS, Windows"
echo "  • Create DEB and RPM packages"
echo "  • Create Flatpak bundle"
echo "  • Publish GitHub release with all artifacts"
echo ""

read -p "Push now? (y/N): " PUSH_NOW
if [[ "$PUSH_NOW" =~ ^[Yy]$ ]]; then
    echo ""
    echo -e "${BLUE}🚀 Pushing to remote...${NC}"
    git push origin main
    git push origin "v$NEW_VERSION"
    echo ""
    echo -e "${GREEN}✅ Release v$NEW_VERSION published!${NC}"
    echo -e "Watch the release at: ${BLUE}https://github.com/timlinux/baboon/actions${NC}"
else
    echo ""
    echo -e "${YELLOW}Release prepared but not pushed.${NC}"
    echo "Run the push commands above when ready."
fi
