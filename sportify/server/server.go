package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/TheVovchenskiy/sportify-backend/api"
	"github.com/TheVovchenskiy/sportify-backend/app"
	"github.com/TheVovchenskiy/sportify-backend/db"
	"github.com/TheVovchenskiy/sportify-backend/pkg/mylogger"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const basicTimeout = 10 * time.Second

func (s *Server) runTgHandler(ctx context.Context, cfg *Config, handler api.Handler, logger *mylogger.MyLogger) error {
	r := chi.NewRouter()
	r.Route(cfg.APIPrefix, func(r chi.Router) {
		r.Use(middleware.Recoverer)
		r.Use(middleware.Logger)

		r.Post("/message", handler.TryCreateEvent)
	})

	s.serverTg = http.Server{ //nolint:exhaustruct
		Addr:                         cfg.PortTg,
		Handler:                      r,
		DisableGeneralOptionsHandler: false,
		ReadTimeout:                  basicTimeout,
		WriteTimeout:                 basicTimeout,
	}

	logger.WithCtx(ctx).Infof("listen bot input %s\n", cfg.PortTg)
	if err := s.serverTg.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

type Server struct {
	serverPublic http.Server
	serverTg     http.Server
	configFile   string
	logger       *mylogger.MyLogger
}

//nolint:funlen
func (s *Server) Run(ctx context.Context, configFile string) error {
	s.configFile = configFile

	cfg, err := NewConfig(configFile)
	if err != nil {
		return err
	}

	logger, err := mylogger.New(cfg.LoggerOutput, cfg.LoggerErrOutput, cfg.ProductionMode)
	if err != nil {
		return err
	}

	defer logger.Sync()
	s.logger = logger

	postgresStorage, err := db.NewPostgresStorage(ctx, cfg.URLDatabase)
	if err != nil {
		return err
	}

	handler := api.NewHandler(app.NewApp(postgresStorage), logger)

	r := chi.NewRouter()
	r.Route(cfg.APIPrefix, func(r chi.Router) {
		r.Use(middleware.Recoverer)
		r.Use(middleware.Logger)
		r.Get("/events", handler.GetEvents)
		r.Get("/event/{id}", handler.GetEvent)
		r.Put("/event/sub/{id}", handler.SubscribeEvent)
		r.Post("/event", handler.CreateEventSite)

		r.Get("/img/*", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/img/" {
				http.Error(w, "404 page not found", http.StatusNotFound)
				return
			}

			fs := http.StripPrefix(cfg.APIPrefix+"img/", http.FileServer(http.Dir(cfg.PathPhotos)))

			fs.ServeHTTP(w, r)
		})
	})

	go func() {
		if err := s.runTgHandler(ctx, cfg, handler, logger); err != nil {
			logger.WithCtx(ctx).Error(err)
		}
	}()

	s.serverPublic = http.Server{ //nolint:exhaustruct
		Addr:                         cfg.PortPublic,
		Handler:                      r,
		DisableGeneralOptionsHandler: false,
		ReadTimeout:                  basicTimeout,
		WriteTimeout:                 basicTimeout,
	}

	logger.WithCtx(ctx).Infof("listen %s\n", cfg.PortPublic)
	if err := s.serverPublic.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.WithCtx(ctx).Infof("shotdown server")

	ctxWithTimeout, cancel := context.WithTimeout(ctx, basicTimeout)
	defer cancel()

	err := s.serverPublic.Shutdown(ctxWithTimeout)
	if err != nil {
		return fmt.Errorf("to shutdown server public: %w", err)
	}

	ctxWithTimeoutTg, cancelTg := context.WithTimeout(ctx, basicTimeout)
	defer cancelTg()

	err = s.serverTg.Shutdown(ctxWithTimeoutTg)
	if err != nil {
		return fmt.Errorf("to shutdown server tg: %w", err)
	}

	return nil
}
