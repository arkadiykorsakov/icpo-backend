package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	authHandler "icpo-backend/internal/api/handlers/auth"
	cityHandler "icpo-backend/internal/api/handlers/city"
	countryHandler "icpo-backend/internal/api/handlers/country"
	crimecategoryHandler "icpo-backend/internal/api/handlers/crime-category"
	eyeColorHandler "icpo-backend/internal/api/handlers/eye-color"
	hairColorHandler "icpo-backend/internal/api/handlers/hair-color"
	healthHandler "icpo-backend/internal/api/handlers/health"
	languageHandler "icpo-backend/internal/api/handlers/language"
	profileHandler "icpo-backend/internal/api/handlers/profile"
	swaggerHandler "icpo-backend/internal/api/handlers/swagger"
	"icpo-backend/internal/api/middleware"
	"icpo-backend/internal/api/routes"
	"icpo-backend/internal/config"
	authDomain "icpo-backend/internal/domain/auth"
	cityDomain "icpo-backend/internal/domain/city"
	countryDomain "icpo-backend/internal/domain/country"
	crimeCategoryDomain "icpo-backend/internal/domain/crime-category"
	eyeColorDomain "icpo-backend/internal/domain/eye-color"
	hairColorDomain "icpo-backend/internal/domain/hair-color"
	languageDomain "icpo-backend/internal/domain/language"
	profileDomain "icpo-backend/internal/domain/profile"
	tokenDomain "icpo-backend/internal/domain/token"
	userDomain "icpo-backend/internal/domain/user"
	"icpo-backend/internal/infrastructure/database"
	"icpo-backend/internal/infrastructure/database/migrator"
	"icpo-backend/pkg/cookies"
)

func main() {
	if err := run(); err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := migrator.Up(ctx, cfg.DatabaseURL); err != nil {
		return err
	}

	pool, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	userRepo := userDomain.NewRepository(pool)
	tokenRepo := tokenDomain.NewRepository(pool)
	countryRepo := countryDomain.NewRepository(pool)
	languageRepo := languageDomain.NewRepository(pool)
	crimeCategoryRepo := crimeCategoryDomain.NewRepository(pool)
	hairColorRepo := hairColorDomain.NewRepository(pool)
	eyeColorRepo := eyeColorDomain.NewRepository(pool)
	cityRepo := cityDomain.NewRepository(pool)

	tokenSvc := tokenDomain.NewService(tokenRepo, tokenDomain.Config{
		AccessTTL:  cfg.AccessTTL,
		RefreshTTL: cfg.RefreshTTL,
		JWTSecret:  cfg.JWTSecret,
	})
	authSvc := authDomain.NewService(
		userRepo,
		tokenSvc,
	)
	profileSvc := profileDomain.NewService(
		userRepo,
		tokenRepo,
	)
	countrySvc := countryDomain.NewService(countryRepo)
	languageSvc := languageDomain.NewService(languageRepo)
	crimeCategorySvc := crimeCategoryDomain.NewService(crimeCategoryRepo)
	hairColorSvc := hairColorDomain.NewService(hairColorRepo)
	eyeColorSvc := eyeColorDomain.NewService(eyeColorRepo)
	citySvc := cityDomain.NewService(cityRepo, countryRepo)

	cookiesManager := cookies.NewManager(cookies.Options{
		Path:     "/api/v1/auth",
		Domain:   cfg.CookieDomain,
		Secure:   cfg.CookieSecure,
		HttpOnly: true,
		SameSite: cfg.CookieSameSite,
	})

	healthH := healthHandler.NewHandler()
	swaggerH := swaggerHandler.NewHandler()
	authH := authHandler.NewHandler(authSvc, cookiesManager)
	profileH := profileHandler.NewHandler(profileSvc)
	countryH := countryHandler.NewHandler(countrySvc)
	languageH := languageHandler.NewHandler(languageSvc)
	crimeCategoryH := crimecategoryHandler.NewHandler(crimeCategorySvc)
	hairColorH := hairColorHandler.NewHandler(hairColorSvc)
	eyeColorH := eyeColorHandler.NewHandler(eyeColorSvc)
	cityH := cityHandler.NewHandler(citySvc)

	router := routes.NewRouter(routes.Deps{
		Health:        healthH,
		Swagger:       swaggerH,
		Auth:          authH,
		Profile:       profileH,
		Country:       countryH,
		Language:      languageH,
		CrimeCategory: crimeCategoryH,
		HairColor:     hairColorH,
		EyeColor:      eyeColorH,
		City:          cityH,
		Authenticate:  middleware.Authenticate(tokenSvc),
		CORS:          middleware.CORS(cfg.CORSAllowedOrigins),
		TrustedOrigin: middleware.RequireTrustedOrigin(cfg.CORSAllowedOrigins),
	})

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)

	go func() {
		slog.Info("server started", "addr", cfg.HTTPAddr)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	case <-ctx.Done():
		slog.Info("shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	}

	return nil
}
