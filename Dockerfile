# syntax=docker/dockerfile:1.7

# ---------------------------------------------------------------------------
# Stage 1: Frontend build
# ---------------------------------------------------------------------------
FROM node:20-alpine@sha256:fb4cd12c85ee03686f6af5362a0b0d56d50c58a04632e6c0fb8363f609372293 AS frontend-build

WORKDIR /app/frontend

# Cache node_modules across builds
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci

# Copy source and build
COPY frontend/ ./
RUN npm run build

# ---------------------------------------------------------------------------
# Stage 2: Go binary build
# ---------------------------------------------------------------------------
FROM golang:1.22-alpine@sha256:1699c10032ca2582ec89a24a1312d986a3f094aed3d5c1147b19880afe40e052 AS backend-build

WORKDIR /app/backend

# Cache Go module downloads
COPY backend/go.mod go.sum ./
RUN go mod download

# Copy source and build a static binary
COPY backend/ ./
RUN go build -trimpath -ldflags="-s -w" -o /usr/local/bin/daedalus ./cmd/daedalus/

# ---------------------------------------------------------------------------
# Stage 3: Runtime — minimal image
# ---------------------------------------------------------------------------
FROM alpine:3.20@sha256:d9e853e87e55526f6b2917df91a2115c36dd7c696a35be12163d44e6e2a4b6bc AS runtime

# Install CA certificates so Go HTTP client can verify TLS connections
RUN apk add --no-cache ca-certificates

WORKDIR /

# Copy the Go binary from the build stage
COPY --from=backend-build /usr/local/bin/daedalus /usr/local/bin/daedalus

# Copy the built frontend static files
# The Go server resolves "frontend/dist" relative to WORKDIR:
#   ../../frontend/dist  from WORKDIR=/  →  /frontend/dist
COPY --from=frontend-build /app/frontend/dist /frontend/dist

EXPOSE 8080

ENV PORT=8080

USER nobody:nogroup

STOPSIGNAL SIGTERM

ENTRYPOINT ["/usr/local/bin/daedalus"]
