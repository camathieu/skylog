.PHONY: dev dev-backend dev-frontend build test docker-build docker-run clean

# ─── Development ────────────────────────────────────────────────────────────────
# Run Go backend and Vite dev server in parallel.
# Vite proxies /api to the Go backend at :8080.
dev:
	$(MAKE) -j2 dev-backend dev-frontend

dev-backend:
	mkdir -p data
	go run ./cmd/skylog --data-dir ./data

dev-frontend:
	cd webapp && npm run dev

# ─── Production build ────────────────────────────────────────────────────────────
# 1. Build the Vue app (output → webapp/dist)
# 2. Compile Go binary (webapp/dist is embedded via //go:embed)
build:
	cd webapp && npm ci && npm run build
	CGO_ENABLED=0 go build -ldflags="-s -w" -o skylog ./cmd/skylog

# ─── Testing ────────────────────────────────────────────────────────────────────
test:
	go test ./...

# ─── Docker ─────────────────────────────────────────────────────────────────────
docker-build:
	docker build -t skylog .

docker-run:
	docker run -d -p 8080:8080 -v skylog-data:/data --name skylog skylog

docker-stop:
	docker rm -f skylog

docker-logs:
	docker logs -f skylog

# ─── Cleanup ────────────────────────────────────────────────────────────────────
clean:
	rm -f skylog
	rm -rf webapp/dist
