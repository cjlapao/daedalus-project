# ==============================================================================
# Daedalus — Development Makefile
# ==============================================================================
#
# Quick start:
#   make          # Show this help
#   make install  # Install dependencies
#   make dev      # Run backend + frontend concurrently (Ctrl-C to stop both)
#
# Targets:
#   help           Show this help message (default target)
#   install        Install backend Go and frontend npm dependencies
#   dev            Run backend and frontend concurrently in dev mode
#   backend        Run only the Go backend (localhost:8080)
#   frontend       Run only the Vite frontend dev server (localhost:5173)
#   build          Compile Go binary and build frontend assets
#   docker-build   Build Docker image(s) via docker compose
#   docker-up      Start services via docker compose up
#   clean          Remove build artifacts (dist, node_modules, compiled binary)
#
# Notes:
#   - Ctrl-C in `make dev` terminates both the backend and frontend processes.
#     A trap on INT/TERM ensures child processes receive the signal.
#   - The backend listens on port 8080 (configurable via PORT env var).
#   - The frontend dev server listens on port 5173 (Vite default).
#   - `make install` runs Go module download first, then npm ci (lockfile).
#     Fall back to `npm install` if no lockfile exists.
#
# ==============================================================================

.SILENT:

# ── Variables ──────────────────────────────────────────────────────────────────

BACKEND_DIR   := backend
FRONTEND_DIR  := frontend
BINARY_NAME   := daedalus
NODE_MODULES  := $(FRONTEND_DIR)/node_modules

# ── Default target ────────────────────────────────────────────────────────────

.DEFAULT_GOAL := help

# ── Help ──────────────────────────────────────────────────────────────────────

.PHONY: help
help:
	@echo "Daedalus — Development Makefile"
	@echo "==============================="
	@echo ""
	@echo "  make install     Install backend Go and frontend npm dependencies"
	@echo "  make dev         Run backend + frontend concurrently (Ctrl-C to stop both)"
	@echo "  make backend     Run only the Go backend (localhost:8080)"
	@echo "  make frontend    Run only the Vite dev server (localhost:5173)"
	@echo "  make build       Compile Go binary and build frontend assets"
	@echo "  make docker-build Build Docker image(s) via docker compose"
	@echo "  make docker-up   Start services via docker compose up"
	@echo "  make clean       Remove build artifacts (dist, node_modules, binary)"
	@echo ""

# ── Install ───────────────────────────────────────────────────────────────────

.PHONY: install
install:
	@echo "==> Installing backend Go dependencies..."
	@cd $(BACKEND_DIR) && go mod download
	@echo "==> Installing frontend npm dependencies..."
	@cd $(FRONTEND_DIR) && if [ -f package-lock.json ]; then npm ci; else npm install; fi
	@echo "Done."

# ── Backend ───────────────────────────────────────────────────────────────────

.PHONY: backend
backend:
	@echo "==> Starting Go backend on :8080..."
	@cd $(BACKEND_DIR) && go run ./cmd/$(BINARY_NAME)/

# ── Frontend ──────────────────────────────────────────────────────────────────

.PHONY: frontend
frontend:
	@echo "==> Starting Vite dev server on :5173..."
	@cd $(FRONTEND_DIR) && npm run dev

# ── Dev (concurrent backend + frontend) ──────────────────────────────────────

.PHONY: dev
dev:
	@echo "==> Starting Daedalus dev environment..."
	@echo "    Backend  → http://localhost:8080"
	@echo "    Frontend → http://localhost:5173"
	@echo "    Ctrl-C will stop both processes."
	@echo ""
	@# Run backend and frontend concurrently.
	@# Trap INT/TERM so both children are killed on Ctrl-C.
	@trap 'kill $${BG_PID} $${FG_PID} 2>/dev/null; wait $${BG_PID} $${FG_PID} 2>/dev/null' INT TERM; \
		cd $(BACKEND_DIR) && go run ./cmd/$(BINARY_NAME)/ & \
		BG_PID=$$!; \
		cd $(FRONTEND_DIR) && npm run dev & \
		FG_PID=$$!; \
		wait $${BG_PID} $${FG_PID} 2>/dev/null

# ── Build ─────────────────────────────────────────────────────────────────────

.PHONY: build
build:
	@echo "==> Building frontend..."
	@cd $(FRONTEND_DIR) && ./node_modules/.bin/tsc && ./node_modules/.bin/vite build
	@echo "==> Compiling Go binary..."
	@cd $(BACKEND_DIR) && go build -o ../$(BINARY_NAME) ./cmd/$(BINARY_NAME)/
	@echo "Build complete."

# ── Docker ────────────────────────────────────────────────────────────────────

.PHONY: docker-build docker-up
docker-build:
	@echo "==> Building Docker image(s)..."
	@docker compose build

docker-up:
	@echo "==> Starting services via docker compose..."
	@docker compose up

# ── Clean ─────────────────────────────────────────────────────────────────────

.PHONY: clean
clean:
	@echo "==> Removing build artifacts..."
	@rm -rf $(FRONTEND_DIR)/dist $(BINARY_NAME) $(NODE_MODULES)
	@echo "Clean complete."
