package server

import (
	"context"
	"net/http"

	"github.com/TheVovchenskiy/sportify-backend/api"
	"github.com/TheVovchenskiy/sportify-backend/app"
	"github.com/TheVovchenskiy/sportify-backend/db"
	"github.com/TheVovchenskiy/sportify-backend/pkg/mylogger"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func runTgHandler(ctx context.Context, cfg *Config, handler api.Handler, logger *mylogger.MyLogger) error {
	r := chi.NewRouter()
	r.Route(cfg.APIPrefix, func(r chi.Router) {
		r.Use(middleware.Recoverer)
		r.Use(middleware.Logger)

		r.Post("/message", handler.TryCreateEvent)
	})

	logger.WithCtx(ctx).Infof("listen bot input %s\n", cfg.PortTg)

	return http.ListenAndServe(cfg.PortTg, r)
}

type Server struct {
	httpServer http.Server
}

func (s *Server) Run(ctx context.Context) error {
	cfg, err := NewConfig()
	if err != nil {
		return err
	}

	logger, err := mylogger.New(cfg.LoggerOutput, cfg.LoggerErrOutput, cfg.ProductionMode)
	if err != nil {
		return err
	}

	simpleEventStorage, err := db.NewSimpleEventStorage()
	if err != nil {
		return err
	}

	handler := api.NewHandler(app.NewApp(simpleEventStorage), logger)

	r := chi.NewRouter()
	r.Route(cfg.APIPrefix, func(r chi.Router) {
		r.Use(middleware.Recoverer)
		r.Use(middleware.Logger)
		r.Get("/events", handler.GetEvents)
		r.Get("/event/{id}", handler.GetEvent)
		r.Put("/event/sub/{id}", handler.SubscribeEvent)

		r.Get("/img/*", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/img/" {
				http.Error(w, "404 page not found", http.StatusNotFound)
				return
			}

			fs := http.StripPrefix("/img/", http.FileServer(http.Dir(cfg.PathPhotos)))

			fs.ServeHTTP(w, r)
		})
	})

	go func() {
		if err := runTgHandler(ctx, cfg, handler, logger); err != nil {
			logger.WithCtx(ctx).Error(err)
		}
	}()

	logger.WithCtx(ctx).Infof("listen %s\n", cfg.PortPublic)

	if err := http.ListenAndServe(cfg.PortPublic, r); err != nil {
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx) //nolint:wrapcheck
}

func (s *Server) ReRun(ctx context.Context) error {
	err := s.Shutdown(ctx)
	if err != nil {
		return err
	}

	return s.Run(ctx)
}
