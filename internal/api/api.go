package api

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/go-pricing-engine/internal/api/routes"
	"github.com/go-pricing-engine/internal/config"
	"github.com/go-pricing-engine/internal/logger"
)

type API struct {
	Router *chi.Mux
	Config *config.Config
	Logger *zap.SugaredLogger
}

func NewServer(logger *zap.SugaredLogger, cfg *config.Config) *API {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))
	routes.SetupRoutes(r)

	return &API{
		Router: r,
		Config: cfg,
		Logger: logger,
	}
}

func (api *API) Run() error {
	srv := http.Server{
		Addr:         api.Config.Port,
		Handler:      api.Router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sign := <-quit

		logger.Logger.Infow("Caught signal", "signal", sign.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		logger.Logger.Infow("Completing background tasks", "addr", srv.Addr)

		shutdownError <- nil
	}()

	logger.Logger.Infow("Starting server", "addr", srv.Addr, "env", api.Config.Env)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	logger.Logger.Infow("Stopped server", "addr", srv.Addr)

	return nil
}
