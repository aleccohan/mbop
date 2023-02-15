package main

import (
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/redhatinsights/mbop/internal/handlers"
	l "github.com/redhatinsights/mbop/internal/logger"
	"github.com/redhatinsights/mbop/internal/middleware"
)

func main() {
	if err := l.Init(); err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	// Emulating the log message at the beginning of mainHandler()
	r.Use(middleware.Logging)

	// TODO: move these to actual handler functions as we figure out which paths
	// are get vs post
	r.Get("/", handlers.Status)
	r.Get("/v*", handlers.CatchAll)
	r.Post("/v*", handlers.CatchAll)
	r.Get("/api/entitlements*", handlers.CatchAll)
	r.Get("/v1/jwt", handlers.JWTV1Handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}

	srv := http.Server{
		Addr:              ":" + port,
		ReadHeaderTimeout: 2 * time.Second,
		Handler:           r,
	}

	l.Log.Info("Starting MBOP Server on", "port", port)
	if err := srv.ListenAndServe(); err != nil {
		l.Log.Error(err, "reason", "server couldn't start")
	}
}
