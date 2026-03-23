# SkyLog — AGENTS.md

AI coding assistant context for the SkyLog project.

## Quick Start

```bash
# Install frontend dependencies
cd webapp && npm install && cd ..

# Development (Go backend on :8080 + Vite dev server on :5173)
make dev

# Production build (embeds frontend into binary)
make build

# Run locally
./skylog --data-dir ./data

# Docker
make docker-build
make docker-run   # runs on http://localhost:8080
```

## Project Structure

See [ARCHITECTURE.md](ARCHITECTURE.md) for the full layout.

## Key Conventions

### Backend (Go)
- Module: `github.com/camathieu/skylog`
- Router: Chi (`github.com/go-chi/chi/v5`)
- ORM: GORM v2 + `github.com/glebarez/sqlite` (pure Go, **no CGO**)
- Always `CGO_ENABLED=0` for builds
- JSON responses: `{"error": "message"}` on failure
- HTTP status codes: 200 OK, 201 Created, 204 No Content, 400 Bad Request, 404 Not Found, 500 Internal Server Error
- Pagination params: `?page=1&per_page=20&order=desc`

### Frontend (Vue 3)
- Vite + Vue 3 (Composition API / `<script setup>`)
- Tailwind CSS v4 via `@tailwindcss/vite`
- API client: `webapp/src/api/jumps.js` — thin `fetch()` wrappers
- Color palette: CSS variables in `src/style.css` (sky blues + sunset orange)
- Responsive: desktop sidebar (≥768px) + mobile bottom nav (<768px)

### go:embed
- The `assets.go` file at the **repo root** (package `assets`) holds the `//go:embed all:webapp/dist` directive
- This file must stay at the root alongside `webapp/` — never move it
- `server/server.go` imports `github.com/camathieu/skylog` (the root `assets` package)

## Adding a New Feature

1. **New model**: add struct to `server/models/`, register in `InitDB()` AutoMigrate list
2. **New API route**: add handler in `server/handlers/`, register route in `server/server.go`
3. **New view**: add `.vue` file to `webapp/src/views/`, add route in `webapp/src/main.js`
4. **New API call**: add method to `webapp/src/api/jumps.js` (or create a new api file)

## Testing

```bash
# Go tests
make test

# Manual API test
curl http://localhost:8080/api/jumps
curl -X POST http://localhost:8080/api/jumps \
  -H 'Content-Type: application/json' \
  -d '{"jumpNumber":1,"date":"2026-03-23T10:00:00Z","dropzone":"Empuriabrava","jumpType":"Fun"}'
```
