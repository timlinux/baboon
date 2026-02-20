#!/usr/bin/env bash
#
# start-backend.sh - Start the Baboon backend server
#
# Usage:
#   ./scripts/start-backend.sh           # Start on default port 8787
#   ./scripts/start-backend.sh -p        # Start with punctuation mode
#   ./scripts/start-backend.sh -port 9000  # Start on custom port
#
# The server runs in the background and writes its PID to:
#   $XDG_RUNTIME_DIR/baboon.pid (or /tmp/baboon.pid)
#
# Use stop-backend.sh to stop the server.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

# Determine PID file location
PID_DIR="${XDG_RUNTIME_DIR:-/tmp}"
PID_FILE="$PID_DIR/baboon.pid"
LOG_FILE="$PID_DIR/baboon.log"

# Check if already running
if [[ -f "$PID_FILE" ]]; then
    PID=$(cat "$PID_FILE")
    if kill -0 "$PID" 2>/dev/null; then
        echo "Baboon backend is already running (PID: $PID)"
        echo "Use ./scripts/stop-backend.sh to stop it first"
        exit 1
    else
        # Stale PID file
        rm -f "$PID_FILE"
    fi
fi

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

echo "Starting Baboon backend server..."
echo "Binary: $BABOON_BIN"
echo "Log file: $LOG_FILE"

# Start server in background with nohup
nohup "$BABOON_BIN" -server "$@" > "$LOG_FILE" 2>&1 &

# Wait a moment for the server to start and write PID file
sleep 0.5

if [[ -f "$PID_FILE" ]]; then
    PID=$(cat "$PID_FILE")
    echo "Backend started successfully (PID: $PID)"
    echo ""
    echo "To check status: ./scripts/status-backend.sh"
    echo "To stop:         ./scripts/stop-backend.sh"
    echo "To view logs:    tail -f $LOG_FILE"
else
    echo "Error: Backend failed to start"
    echo "Check logs: cat $LOG_FILE"
    exit 1
fi
