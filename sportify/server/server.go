package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/TheVovchenskiy/sportify-backend/api"
	"github.com/TheVovchenskiy/sportify-backend/app"
	"github.com/TheVovchenskiy/sportify-backend/app/config"
	"github.com/TheVovchenskiy/sportify-backend/db"
	sportifymiddleware "github.com/TheVovchenskiy/sportify-backend/pkg/middleware"
	"github.com/TheVovchenskiy/sportify-backend/pkg/mylogger"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const basicTimeout = 10 * time.Second

func (s *Server) runTgHandler(ctx context.Context, cfg *config.Config, handler api.Handler, logger *mylogger.MyLogger) error {
	r := chi.NewRouter()
	r.Route(cfg.App.APIPrefix, func(r chi.Router) {
		r.Use(middleware.Recoverer)
		r.Use(middleware.Logger)

		r.Post("/message", handler.TryCreateEvent)
	})

	s.serverTg = http.Server{ //nolint:exhaustruct
		Addr:                         cfg.Bot.Port,
		Handler:                      r,
		DisableGeneralOptionsHandler: false,
		ReadTimeout:                  basicTimeout,
		WriteTimeout:                 basicTimeout,
	}

	logger.Infof("listen bot input %s\n", cfg.Bot.Port)
	if err := s.serverTg.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

type Server struct {
	serverPublic http.Server
	serverTg     http.Server
	configPaths  []string
	logger       *mylogger.MyLogger
}

//nolint:funlen
func (s *Server) Run(ctx context.Context, configFile []string) error {
	s.configPaths = configFile

	err := config.InitConfig(configFile)
	if err != nil {
		return err
	}

	cfg := config.GetGlobalConfig()

	logger, err := mylogger.New(
		cfg.Logger.LoggerOutput,
		cfg.Logger.LoggerErrOutput,
		cfg.Logger.ProductionMode,
	)
	if err != nil {
		return err
	}
	defer logger.Sync()
	s.logger = logger

	config.WatchRemoteConfig(logger)

	postgresStorage, err := db.NewPostgresStorage(ctx, cfg.Postgres.URL)
	if err != nil {
		return err
	}

	fsStorage, err := db.NewFileSystemStorage(cfg.PathPhotos)
	if err != nil {
		return fmt.Errorf("to new fs storage: %w", err)
	}

	handler := api.NewHandler(app.NewApp(cfg.URLPrefixFile, fsStorage, postgresStorage), logger, cfg.FolderID, cfg.IAMToken)

	r := chi.NewRouter()
	r.Route(cfg.App.APIPrefix, func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Use(middleware.RequestID)
		r.Use(sportifymiddleware.Config)
		r.Get("/events", handler.GetEvents)
		r.Get("/event/{id}", handler.GetEvent)
		r.Put("/event/{id}", handler.EditEventSite)
		r.Delete("/event/{id}", handler.DeleteEvent)
		r.Put("/event/sub/{id}", handler.SubscribeEvent)
		r.Post("/event", handler.CreateEventSite)
		r.Post("/upload", handler.UploadFile)

		r.Get("/img/*", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/img/" {
				http.Error(w, "404 page not found", http.StatusNotFound)
				return
			}

			fs := http.StripPrefix(cfg.App.APIPrefix+"img/", http.FileServer(http.Dir(cfg.App.PathPhotos)))

			fs.ServeHTTP(w, r)
		})
	})

	go func() {
		if err := s.runTgHandler(ctx, cfg, handler, logger); err != nil {
			logger.Error(err)
		}
	}()

	s.serverPublic = http.Server{ //nolint:exhaustruct
		Addr:                         cfg.App.Port,
		Handler:                      r,
		DisableGeneralOptionsHandler: false,
		ReadTimeout:                  basicTimeout,
		WriteTimeout:                 basicTimeout,
	}

	logger.Infof("listen %s\n", cfg.App.Port)
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
