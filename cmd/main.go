package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/go-pricing-engine/internal/config"
	"github.com/go-pricing-engine/internal/logger"
)

func main() {
	logger.InitLogger()
	defer logger.CloseLogger()

	// Load .env file
	if err := godotenv.Load(); err != nil {
		logger.Logger.Fatal("Error loading .env file", zap.Error(err))
	}

	// Load config
	cfg := config.NewConfig()
	if err := cfg.ParseFlags(); err != nil {
		logger.Logger.Fatal("Failed to parse command-line flags", zap.Error(err))
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	http.ListenAndServe(":3000", r)
}
