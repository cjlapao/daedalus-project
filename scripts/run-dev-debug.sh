#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

echo "=== Daedalus Dev Loop (Debug Mode) ==="
echo ""

# --- Backend: Delve headless ---
echo "[1/2] Starting Go backend (Delve headless on :40000)..."
cd "$ROOT_DIR/backend"

if ! command -v dlv &>/dev/null; then
    echo "ERROR: Delve (dlv) is not installed."
    echo "Install it with: go install github.com/go-delve/delve/cmd/dlv@latest"
    exit 1
fi

dlv debug --listen=:40000 --headless --api-version=2 --accept-multiclient \
    --continue --check-out-go-version=false ./cmd/daedalus/ &
BACKEND_PID=$!

# Wait for Delve to accept connections
for i in $(seq 1 20); do
    if curl -s http://localhost:40000/version >/dev/null 2>&1; then
        echo "  Delve ready on :40000"
        break
    fi
    if ! kill -0 $BACKEND_PID 2>/dev/null; then
        echo "ERROR: Delve process exited unexpectedly."
        exit 1
    fi
    sleep 0.5
done

# --- Frontend: Vite dev server ---
echo "[2/2] Starting frontend dev server (Vite on :5173)..."
cd "$ROOT_DIR/frontend"

if ! command -v npm &>/dev/null; then
    echo "ERROR: npm is not installed."
    exit 1
fi

npm run dev &
FRONTEND_PID=$!

echo ""
echo "========================================"
echo "  Backend debug : http://localhost:40000"
echo "  Frontend dev  : http://localhost:5173"
echo "  Press Ctrl+C to stop everything"
echo "========================================"
echo ""

cleanup() {
    echo ""
    echo "Shutting down..."
    kill $BACKEND_PID 2>/dev/null || true
    kill $FRONTEND_PID 2>/dev/null || true
    wait $BACKEND_PID 2>/dev/null || true
    wait $FRONTEND_PID 2>/dev/null || true
    exit 0
}
trap cleanup INT TERM

wait
