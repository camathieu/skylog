# Build stage: frontend
FROM node:22-alpine AS frontend
WORKDIR /app/webapp
COPY webapp/package*.json ./
RUN npm ci
COPY webapp/ ./
RUN npm run build

# Build stage: Go binary (frontend dist embedded via go:embed)
FROM golang:1.24-alpine AS backend
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
COPY --from=frontend /app/webapp/dist ./webapp/dist
# Pure Go SQLite — no CGO needed
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o skylog ./cmd/skylog

# Runtime: minimal Alpine image
FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=backend /app/skylog /usr/local/bin/
EXPOSE 8080
VOLUME ["/data"]
ENTRYPOINT ["skylog"]
CMD ["--data-dir", "/data"]
