# Daedalus

A code-driven agent orchestration engine.

## Prerequisites

- **[Go 1.22](https://go.dev/dl/)** or later
- **[Node.js 20](https://nodejs.org/)** (current LTS) or later

## Quick Start

### Install dependencies

```bash
make install
```

Downloads Go modules and installs frontend npm packages (using `npm ci` with the lockfile when available).

### Run locally

```bash
make dev
```

This starts both services concurrently:

| Service | URL | Port |
|---------|-----|------|
| Go backend | http://localhost:8080 | 8080 |
| Vite frontend | http://localhost:5173 | 5173 |

The Vite dev server proxies all `/api` requests to the Go backend, so the frontend and backend communicate seamlessly during development.

Ctrl-C stops both processes.

### Run a single service

```bash
make backend    # Go backend only (:8080)
make frontend   # Vite dev server only (:5173)
```

## Production / Docker

### Build

Compile the Go binary and bundle the frontend assets:

```bash
make build
```

### Docker

Build and run with Docker Compose:

```bash
make docker-build   # Build the image
make docker-up      # Start the container
```

Production runs a single Go process that serves both the API and the static frontend on port `8080`.

## Directory Structure

```
├── backend/            # Go server
│   ├── cmd/            # Entry points
│   └── internal/       # Application logic
├── frontend/           # React + TypeScript + Vite
│   └── src/            # Frontend source
├── Dockerfile          # Multi-stage Docker build
├── docker-compose.yml  # Docker Compose config
├── Makefile            # Build and dev targets
└── README.md           # You are here
```

## Go Module Path

The Go module is declared as `github.com/cjlapao/daedalus/backend`. If you fork this project, update the module path in `backend/go.mod` to match your own GitHub username.
