# SkyLog — ARCHITECTURE.md

## Overview

SkyLog is a self-hosted personal skydiving logbook with a Go backend and a Vue 3 frontend embedded into the binary at compile time via `//go:embed`.

## Project Layout

```
skylog/
├── assets.go              # package assets — //go:embed all:webapp/dist
├── cmd/skylog/main.go     # Entrypoint: flags, DB init, start server
├── server/
│   ├── server.go          # Chi router, CORS, API routes, SPA fallback
│   ├── handlers/jump.go   # Jump CRUD HTTP handlers
│   └── models/
│       ├── database.go    # InitDB(): opens SQLite, sets WAL pragmas, AutoMigrate
│       └── jump.go        # Jump GORM model
├── webapp/                # Vue 3 + Vite + Tailwind CSS frontend
│   ├── src/
│   │   ├── api/jumps.js   # fetch() API client
│   │   ├── components/    # AppLayout (sidebar + bottom nav)
│   │   └── views/         # JumpList, JumpForm
│   └── dist/              # Built output — embedded into Go binary
├── Dockerfile             # Multi-stage: Node → Go → Alpine
└── Makefile
```

## Key Design: go:embed

```
webapp/dist/ (build output)
      ↓
assets.go → //go:embed all:webapp/dist → assets.FS (embed.FS)
      ↓
server/server.go → fs.Sub(assets.FS, "webapp/dist") → http.FileServer
```

The root-level `assets.go` (package `assets`) must live alongside `webapp/` so `//go:embed` can reference `webapp/dist` without `..` traversal (which Go forbids).

## HTTP Routing

| Method | Path           | Handler              |
|--------|----------------|----------------------|
| GET    | /api/jumps     | ListJumps            |
| POST   | /api/jumps     | CreateJump           |
| GET    | /api/jumps/{id}| GetJump              |
| PUT    | /api/jumps/{id}| UpdateJump           |
| DELETE | /api/jumps/{id}| DeleteJump           |
| GET    | /*             | Embedded SPA (Vue)   |

Non-`/api` routes serve `index.html` as a fallback for Vue Router history mode.

## Database

- **Driver**: `github.com/glebarez/sqlite` (pure Go wrapper over modernc/sqlite — no CGO)
- **ORM**: GORM v2 with AutoMigrate
- **Pragmas**: `journal_mode=WAL`, `foreign_keys=ON`
- **Location**: `{data-dir}/skylog.db`

## Frontend

- **Vue Router**: history mode, SPA
- **Tailwind CSS v4** via `@tailwindcss/vite` plugin
- **Dev proxy**: Vite proxies `/api` → `http://localhost:8080`
- **Color palette**: CSS variables — deep sky blues + sunset orange

## Data Volume

| Path         | Contents                |
|--------------|-------------------------|
| `/data`      | Docker volume           |
| `/data/skylog.db` | SQLite database    |

## Phase Roadmap

| Phase | Status |
|-------|--------|
| PoC   | ✅ Jump CRUD, responsive UI, go:embed, Docker |
| MVP   | ⬜ Search/filter, tags, stats, documents |
| GA    | ⬜ Gemini OCR import, DZ map, polish |
| V2    | ⬜ PWA, auth, cloud sync |
