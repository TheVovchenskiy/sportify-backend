package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/TheVovchenskiy/sportify-backend/api"
	"github.com/TheVovchenskiy/sportify-backend/app"
	"github.com/TheVovchenskiy/sportify-backend/app/botapi"
	"github.com/TheVovchenskiy/sportify-backend/app/config"
	"github.com/TheVovchenskiy/sportify-backend/app/telegramapi"
	"github.com/TheVovchenskiy/sportify-backend/db"
	sportifymiddleware "github.com/TheVovchenskiy/sportify-backend/pkg/middleware"
	"github.com/TheVovchenskiy/sportify-backend/pkg/mylogger"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-pkgz/auth"
	"github.com/go-pkgz/auth/avatar"
	authmiddleware "github.com/go-pkgz/auth/middleware"
	"github.com/go-pkgz/auth/provider"
	"github.com/go-pkgz/auth/token"
)

const basicTimeout = 10 * time.Second

func (s *Server) runTgHandler(ctx context.Context, cfg *config.Config, handler api.Handler, logger *mylogger.MyLogger) error {
	r := chi.NewRouter()
	r.Route(cfg.App.APIPrefix, func(r chi.Router) {
		r.Use(middleware.Recoverer)
		r.Use(middleware.Logger)

		r.Post("/message", handler.TryCreateEvent)
		r.Put("/events/{id}/subscribers", handler.SubscribeEvent)
		r.Get("/events/{event_id}/subscribers", handler.UserIsSubscribed)
		r.Post("/users", handler.LoginUserFromTg)
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

	logger.Debugf("Config: %v", cfg)

	postgresStorage, _, err := db.NewPostgresStorage(ctx, cfg.Postgres.URL)
	if err != nil {
		return err
	}

	fsStorage, err := db.NewFileSystemStorage(cfg.App.PathPhotos)
	if err != nil {
		return fmt.Errorf("to new fs storage: %w", err)
	}

	mapTokenStorage := db.NewMapTokenStorage(time.Minute * 30)

	botAPI, err := botapi.NewBotAPI(cfg.BotAPI.BaseURL, cfg.BotAPI.Port)
	if err != nil {
		return fmt.Errorf("to new bot api: %w", err)
	}

	url := cfg.App.Domain + cfg.App.Port
	appSportify := app.NewApp(
		cfg.App.URLPrefixFile, fsStorage, postgresStorage, postgresStorage, mapTokenStorage, logger, botAPI,
	)

	tgAPI := telegramapi.NewTelegramAPIDummy()
	handler := api.NewHandler(appSportify, logger, cfg.App.FolderID, cfg.App.IAMToken, cfg.App.Domain, cfg.App.Port, cfg.App.APIPrefix, tgAPI)
	checkCredFunc := handler.NewCredCheckFunc(ctx)
	authMiddleware, authHandler, tokenServiceProvider := s.prepareAuthProviders(
		ctx,
		cfg.App.AuthSecret, url, cfg.Bot.Token,
		checkCredFunc, http.DefaultClient, mapTokenStorage, &handler, &handler,
		logger, tgAPI,
	)

	handler.AddTokenServiceProvider(tokenServiceProvider)

	r := chi.NewRouter()
	r.Route(cfg.App.APIPrefix, func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Use(middleware.RequestID)
		r.Use(sportifymiddleware.Config)
		r.Use(sportifymiddleware.ConvertErrUnknownToOurType)
		r.Get("/healthcheck", handler.Healthcheck)
		r.Get("/events", handler.FindEvents)
		r.Get("/event/{id}", handler.GetEvent)
		r.Get("/profiles/{id}", handler.GetProfile)
		r.With(authMiddleware.Auth).Put("/event/{id}", handler.EditEventSite)
		r.With(authMiddleware.Auth).Delete("/event/{id}", handler.DeleteEvent)
		r.With(authMiddleware.Auth).Put("/event/sub/{id}", handler.SubscribeEvent)
		r.With(authMiddleware.Auth).Post("/event", handler.CreateEventSite)
		r.With(authMiddleware.Auth).Get("/users/{id}/events", handler.GetUsersEvents)
		r.With(authMiddleware.Auth).Get("/users/{id}/sub_active/events", handler.GetUsersSubActiveEvents)
		r.With(authMiddleware.Auth).Get("/users/{id}/sub_archive/events", handler.GetUsersSubArchiveEvents)
		r.With(authMiddleware.Auth).Post("/upload", handler.UploadFile)
		r.With(authMiddleware.Auth).Put("/profiles/{user_id}", handler.UpdateProfile)

		r.Mount("/auth",
			sportifymiddleware.ConvertLoginResponseToCheck(
				&handler,
				sportifymiddleware.PostOnlyRestriction("/logout", authHandler)))
		r.Post("/auth/register", handler.Register)
		r.With(authMiddleware.Auth).Get("/auth/check", handler.Check)

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

func (s *Server) prepareAuthProviders(
	_ context.Context,
	authSecret, url, tgToken string,
	credCheckerFunc provider.CredCheckerFunc,
	httpClient *http.Client,
	storageToken StorageToken,
	checkHandler sportifymiddleware.CheckHandler,
	claimsUpdater token.ClaimsUpdater,
	logger *mylogger.MyLogger,
	tgAPI provider.TelegramAPI,
) (authmiddleware.Authenticator, http.Handler, *token.Service) {
	options := auth.Opts{
		ClaimsUpd: claimsUpdater,
		SecretReader: token.SecretFunc(func(id string) (string, error) {
			// TODO: refresh every day
			return authSecret, nil
		}),
		SecureCookies:  true,
		TokenDuration:  time.Minute * 5,
		CookieDuration: time.Hour * 24,
		AvatarStore:    &avatar.NoOp{},
		DisableXSRF:    true,
		SameSiteCookie: http.SameSiteLaxMode,
		Issuer:         "move-life.ru",
		URL:            url,
		Logger:         logger,
	}

	service := auth.NewService(options)

	service.AddDirectProvider("my", credCheckerFunc)

	// tgAPI := NewTelegramAPIDummy()

	tgHandler := provider.TelegramHandler{
		ProviderName: "telegram",
		ErrorMsg:     "❌ Произошла ошибка входа, попробуйте еще раз.",
		SuccessMsg:   "✅ Вы успешно вошли, вернитесь на сайт",
		Telegram:     tgAPI,
		L:            logger,
		TokenService: service.TokenService(),
		AvatarSaver:  nil,
	}

	go func() {
		err := tgHandler.Run(context.Background())
		if err != nil {
			logger.Errorf("FAIL: to start tg auth provider: %+v", err)
		}
	}()

	// add wrapper for rewrite response on Login
	wrapTgHandler := NewWrapperTgHandler(storageToken, checkHandler, &tgHandler)

	service.AddCustomHandler(wrapTgHandler)

	middlewares := service.Middleware()
	authRoutes, _ := service.Handlers()

	return middlewares, authRoutes, service.TokenService()
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
