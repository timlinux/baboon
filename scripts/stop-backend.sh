#!/usr/bin/env bash
#
# stop-backend.sh - Stop the Baboon backend server
#
# Usage:
#   ./scripts/stop-backend.sh        # Graceful shutdown (SIGTERM)
#   ./scripts/stop-backend.sh -f     # Force kill (SIGKILL)

set -euo pipefail

FORCE=false
if [[ "${1:-}" == "-f" ]] || [[ "${1:-}" == "--force" ]]; then
    FORCE=true
fi

# Determine PID file location
PID_DIR="${XDG_RUNTIME_DIR:-/tmp}"
PID_FILE="$PID_DIR/baboon.pid"

if [[ ! -f "$PID_FILE" ]]; then
    echo "Baboon backend is not running (no PID file found)"
    exit 0
fi

PID=$(cat "$PID_FILE")

if ! kill -0 "$PID" 2>/dev/null; then
    echo "Baboon backend is not running (stale PID file)"
    rm -f "$PID_FILE"
    exit 0
fi

if $FORCE; then
    echo "Force killing Baboon backend (PID: $PID)..."
    kill -9 "$PID" 2>/dev/null || true
else
    echo "Stopping Baboon backend (PID: $PID)..."
    kill "$PID" 2>/dev/null || true
fi

# Wait for process to terminate
for i in {1..10}; do
    if ! kill -0 "$PID" 2>/dev/null; then
        echo "Backend stopped successfully"
        rm -f "$PID_FILE"
        exit 0
    fi
    sleep 0.5
done

# If still running after 5 seconds
if kill -0 "$PID" 2>/dev/null; then
    echo "Backend did not stop gracefully, force killing..."
    kill -9 "$PID" 2>/dev/null || true
    rm -f "$PID_FILE"
fi

echo "Backend stopped"
