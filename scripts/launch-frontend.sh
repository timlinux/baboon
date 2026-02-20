#!/usr/bin/env bash
#
# launch-frontend.sh - Launch a Baboon frontend client
#
# Usage:
#   ./scripts/launch-frontend.sh           # Connect to default port 8787
#   ./scripts/launch-frontend.sh -p        # With punctuation mode
#   ./scripts/launch-frontend.sh -port 9000  # Connect to custom port
#
# This connects to an existing backend server. Start the backend first with:
#   ./scripts/start-backend.sh
#
# Multiple frontends can connect to the same backend simultaneously.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

PORT=8787
EXTRA_ARGS=()

# Parse arguments to extract port
while [[ $# -gt 0 ]]; do
    case $1 in
        -port)
            PORT="$2"
            EXTRA_ARGS+=("$1" "$2")
            shift 2
            ;;
        *)
            EXTRA_ARGS+=("$1")
            shift
            ;;
    esac
done

# Find baboon binary
BABOON_BIN=""
if [[ -x "$PROJECT_DIR/baboon" ]]; then
    BABOON_BIN="$PROJECT_DIR/baboon"
elif [[ -x "$PROJECT_DIR/result/bin/baboon" ]]; then
    BABOON_BIN="$PROJECT_DIR/result/bin/baboon"
elif command -v baboon &>/dev/null; then
    BABOON_BIN="$(command -v baboon)"
else
    echo "Error: baboon binary not found"
    echo "Build it with: nix build (or go build)"
    exit 1
fi

# Check if backend is running
HEALTH_URL="http://127.0.0.1:$PORT/api/health"
if ! curl -s --connect-timeout 2 "$HEALTH_URL" &>/dev/null; then
    echo "Error: Backend is not running on port $PORT"
    echo ""
    echo "Start the backend first with:"
    echo "  ./scripts/start-backend.sh"
    echo ""
    echo "Or run in combined mode:"
    echo "  baboon"
    exit 1
fi

# Launch frontend
exec "$BABOON_BIN" -client "${EXTRA_ARGS[@]}"
