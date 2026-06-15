# Daedalus — Project Scaffolding Task

You are scaffolding the initial skeleton for **Project Daedalus**, a code-driven agent
orchestration engine. This task creates the empty, runnable monorepo structure only — **no
business logic**. The goal is a skeleton that builds, runs, and serves a health check end
to end, so feature work can begin on a known-good foundation.

Do not implement orchestration, agents, planner, or any domain logic. Build the shell and
prove it runs.

---

## 1. What you are building

A **monorepo** containing:

- A **Go backend** (the orchestrator API — for now just a health endpoint).
- A **React + TypeScript frontend** (the dashboard — for now just a page that calls the
  backend health endpoint and shows the result).
- A **Makefile** that runs both in debug/development mode.
- A **Docker** setup so the whole thing runs in a container.

In **development**, the two run as separate processes: the Vite dev server serves the UI
with hot reload and **proxies** API calls to the Go backend. In **production**, the Go
binary serves the pre-built React static assets as a single process. Scaffold both, but
the day-to-day developer loop is the dev/proxy mode.

---

## 2. Required directory layout

Create exactly this structure at the repository root:

```
daedalus/
├── Makefile
├── README.md
├── .gitignore
├── .dockerignore
├── docker-compose.yml
├── Dockerfile
├── backend/
│   ├── go.mod
│   ├── go.sum
│   ├── cmd/
│   │   └── daedalus/
│   │       └── main.go          # entrypoint: starts HTTP server
│   └── internal/
│       ├── server/
│       │   └── server.go        # router, server setup
│       └── handlers/
│           └── health.go        # GET /api/health → {"status":"ok"}
└── frontend/
    ├── package.json
    ├── tsconfig.json
    ├── vite.config.ts           # includes dev proxy for /api → backend
    ├── index.html
    └── src/
        ├── main.tsx
        ├── App.tsx              # calls /api/health, renders the result
        └── vite-env.d.ts
```

Do not add extra packages, state libraries, UI kits, or tooling beyond what is listed
here. Keep the skeleton minimal.

---

## 3. Backend specifics (Go)

- **Module path:** `github.com/<OWNER>/daedalus/backend` — leave `<OWNER>` as a clearly
  marked placeholder if the owner is unknown, and note it in the README.
- **Go version:** use a current stable Go (1.22+). Pin it in `go.mod`.
- **HTTP:** standard library `net/http` with `http.ServeMux` is sufficient. Do **not**
  pull in a web framework.
- **Port:** backend listens on `:8080`. Make it overridable via a `PORT` env var,
  defaulting to `8080`.
- **Endpoint:** `GET /api/health` returns HTTP 200 with JSON body
  `{"status":"ok","service":"daedalus"}` and `Content-Type: application/json`.
- **Production static serving (scaffold, not active in dev):** include a code path where,
  if a built frontend directory (`frontend/dist`) is present, the server serves it for
  non-`/api` routes. It is fine for this to be a no-op when the directory does not exist.
- **Logging:** log on startup which port it bound to. Keep it to the standard library.
- `main.go` wires config → server → listen. `server.go` builds the mux and routes.
  `health.go` holds the handler. Keep these thin.

---

## 4. Frontend specifics (React + TypeScript)

- **Tooling:** Vite, React, TypeScript. Use the standard Vite React-TS setup.
- **Dev server port:** `5173` (Vite default).
- **Dev proxy (critical):** in `vite.config.ts`, proxy `/api` to `http://localhost:8080`
  so that in development the UI calls `/api/health` and Vite forwards it to the Go backend.
  This is the "app calls the internal backend" wiring.
- **App behavior:** on load, `App.tsx` fetches `/api/health`, then renders either the
  returned status or an error state. Plain, unstyled or minimally styled — this only needs
  to prove the round trip works. No router, no component library.
- **Build output:** `npm run build` emits to `frontend/dist` (Vite default), which is what
  the Go server serves in production mode.

---

## 5. Makefile (debug/dev is the priority)

Provide at least these targets. Keep them readable.

- `make help` — list targets (make this the default goal).
- `make install` — install backend Go deps and frontend npm deps.
- `make dev` — run backend and frontend **together** in debug mode (backend on `:8080`,
  Vite on `:5173` with the proxy). Running both concurrently from one target is the key
  developer-experience requirement; a simple background-process approach or a documented
  two-terminal fallback is acceptable, but prefer single-command concurrent run.
- `make backend` — run only the Go backend in dev.
- `make frontend` — run only the Vite dev server.
- `make build` — build the frontend (`frontend/dist`) and compile the Go binary.
- `make docker-build` — build the Docker image.
- `make docker-up` — run the container via docker-compose.
- `make clean` — remove build artifacts (`frontend/dist`, compiled binary, etc.).

If `make dev` runs two processes, ensure Ctrl-C cleans both up, and document the behavior
in the README.

---

## 6. Docker

- **Dockerfile:** multi-stage.
  1. Stage one builds the frontend (`node` image → `npm ci` → `npm run build` → `dist`).
  2. Stage two builds the Go binary (`golang` image), copying in the built `dist`.
  3. Final minimal runtime stage (small base) holds the Go binary plus the static `dist`,
     runs the single Go process which serves both the API and the static UI on `:8080`.
- **docker-compose.yml:** one service running the image, mapping `8080:8080`, with `PORT`
  configurable. This compose path represents the **production-style single-process** run
  (Go serves the built UI), not the hot-reload dev loop.
- Keep the final image lean. Do not ship the Node toolchain or Go toolchain in the runtime
  stage.

---

## 7. Repo hygiene

- `.gitignore` must cover: Go build artifacts and the compiled binary, `node_modules`,
  `frontend/dist`, common editor/OS files, and any local env files.
- `.dockerignore` must exclude `node_modules`, `frontend/dist`, `.git`, and build
  artifacts so the build context stays small.
- `README.md` must document: what Daedalus is (one or two lines pointing to the roadmap),
  prerequisites (Go version, Node version), and the exact commands to get running —
  `make install` then `make dev` — plus the production/docker path. Note the `<OWNER>`
  placeholder if the module path was left generic.

---

## 8. Definition of done (verify before declaring complete)

Do not report success until **all** of these hold:

1. `make install` completes with no errors on a clean checkout.
2. `make dev` brings up both processes; the Go backend logs that it bound `:8080` and Vite
   serves on `:5173`.
3. Opening the Vite dev URL loads the page, and the page successfully fetches
   `/api/health` **through the proxy** and displays `status: ok`. (The browser calls Vite;
   Vite forwards to Go; Go responds. This proves the full dev round trip.)
4. `curl http://localhost:8080/api/health` returns
   `{"status":"ok","service":"daedalus"}` with a 200 and JSON content type.
5. `make build` produces `frontend/dist` and a compiled Go binary with no errors.
6. `make docker-build` builds successfully, and `make docker-up` serves the app on
   `http://localhost:8080` — the page loads **and** `/api/health` responds from the single
   container process (production-style, Go serving static + API).
7. The directory layout matches section 2.

Report what you created, the commands you ran to verify each item above, and their output.
If any step cannot pass, stop and report the blocker rather than working around it.

---

## 9. Out of scope (do not build)

- Any orchestration, planner, agent, executor, event store, or database code.
- Authentication, websockets/SSE, the live dashboard, swim lanes.
- Styling systems, component libraries, state management.
- CI configuration.

Those come in later phases. This task is the runnable skeleton and nothing more.
