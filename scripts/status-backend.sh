#!/usr/bin/env bash
#
# status-backend.sh - Check the status of the Baboon backend server
#
# Usage:
#   ./scripts/status-backend.sh              # Check default port 8787
#   ./scripts/status-backend.sh -port 9000   # Check custom port
#
# Exit codes:
#   0 - Backend is running and healthy
#   1 - Backend is not running

set -euo pipefail

PORT=8787

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -port)
            PORT="$2"
            shift 2
            ;;
        *)
            shift
            ;;
    esac
done

# Determine PID file location
PID_DIR="${XDG_RUNTIME_DIR:-/tmp}"
PID_FILE="$PID_DIR/baboon.pid"
LOG_FILE="$PID_DIR/baboon.log"

echo "Baboon Backend Status"
echo "====================="
echo ""

# Check PID file
if [[ -f "$PID_FILE" ]]; then
    PID=$(cat "$PID_FILE")
    if kill -0 "$PID" 2>/dev/null; then
        echo "Process:  Running (PID: $PID)"
    else
        echo "Process:  Not running (stale PID file)"
        rm -f "$PID_FILE"
    fi
else
    echo "Process:  Not running (no PID file)"
fi

# Check health endpoint
echo -n "Health:   "
HEALTH_URL="http://127.0.0.1:$PORT/api/health"
if HEALTH=$(curl -s --connect-timeout 2 "$HEALTH_URL" 2>/dev/null); then
    STATUS=$(echo "$HEALTH" | grep -o '"status":"[^"]*"' | cut -d'"' -f4)
    SESSIONS=$(echo "$HEALTH" | grep -o '"active_sessions":[0-9]*' | cut -d':' -f2)
    echo "Healthy"
    echo "Sessions: $SESSIONS active"
    echo ""
    echo "API URL:  http://127.0.0.1:$PORT"
else
    echo "Not responding"
    echo ""
    if [[ -f "$LOG_FILE" ]]; then
        echo "Recent logs:"
        echo "------------"
        tail -5 "$LOG_FILE" 2>/dev/null || echo "(no logs)"
    fi
    exit 1
fi

# Show session details if any
if [[ "$SESSIONS" -gt 0 ]]; then
    echo ""
    echo "Active Sessions:"
    echo "----------------"
    SESSIONS_URL="http://127.0.0.1:$PORT/api/sessions"
    curl -s "$SESSIONS_URL" 2>/dev/null | \
        grep -o '"id":"[^"]*"' | \
        cut -d'"' -f4 | \
        while read -r id; do
            echo "  - $id"
        done
fi

echo ""
exit 0
