package server

import (
	"fmt"
	"io/fs"
	"net/http"

	assets "github.com/camathieu/skylog"
	"github.com/camathieu/skylog/server/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"
)

// New creates and configures the HTTP server.
func New(db *gorm.DB, listen string) (*http.Server, error) {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// CORS for Vite dev server hot-reload
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	// --- API routes ---
	jumpHandler := handlers.NewJumpHandler(db)
	r.Route("/api", func(r chi.Router) {
		r.Route("/jumps", func(r chi.Router) {
			r.Get("/", jumpHandler.ListJumps)
			r.Post("/", jumpHandler.CreateJump)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", jumpHandler.GetJump)
				r.Put("/", jumpHandler.UpdateJump)
				r.Delete("/", jumpHandler.DeleteJump)
			})
		})
	})

	// --- Embedded frontend (SPA fallback) ---
	// assets.FS is the go:embed FS rooted at the repo root, so subtract "webapp/dist"
	distFS, err := fs.Sub(assets.FS, "webapp/dist")
	if err != nil {
		return nil, fmt.Errorf("sub embed fs: %w", err)
	}
	fileServer := http.FileServer(http.FS(distFS))

	// For any non-/api path: serve static files or fall back to index.html (Vue Router).
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		// Check if the file exists in the embedded FS
		path := r.URL.Path
		if path == "/" || path == "" {
			path = "index.html"
		} else {
			path = path[1:] // strip leading /
		}

		if _, err := distFS.Open(path); err != nil {
			// Not found → serve index.html for client-side routing
			r.URL.Path = "/"
		}
		fileServer.ServeHTTP(w, r)
	})

	return &http.Server{
		Addr:    listen,
		Handler: r,
	}, nil
}
