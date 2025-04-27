package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"good_shortener/internal/config"
	"good_shortener/internal/http-server/handlers/auth/login"
	"good_shortener/internal/http-server/handlers/redirect"
	"good_shortener/internal/http-server/handlers/url/getUrls"
	"good_shortener/internal/http-server/handlers/url/save"
	"good_shortener/internal/lib/logger/sl"
	"good_shortener/internal/middlewares/jwtMiddleware"
	"good_shortener/internal/storage/postgres"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// TODO: init config: cleanenv
	cfg := config.MustLoad()
	fmt.Println(cfg)

	// TODO: init logger: slog
	log := setupLogger(cfg.Env)

	log.Info("starting app...", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// TODO: init storge: postgres
	storage, err := postgres.New(cfg.Name, cfg.Password, cfg.Port)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	e := echo.New()
	myJwtMiddleware := jwtMiddleware.New()

	e.Use(middleware.RequestID())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		AllowMethods: []string{
			http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete,
		},
		AllowOrigins: []string{"http://localhost:3000"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))
	// e.Use(middleware.Logger()) // TODO: create custom logger middleware, because that off main logger
	e.Use(middleware.Recover())

	authGroup := e.Group("/url")
	authGroup.Use(myJwtMiddleware)

	// e.POST("auth/registration", registration.New(log, storage))
	e.POST("auth/login", login.New(log, storage))
	authGroup.POST("/save", save.New(log, storage, storage))
	authGroup.GET("/users_all", getUrls.New(log, storage, storage))

	e.GET("/:alias", redirect.New(log, storage))
	// TODO: add delete handler

	s := http.Server{
		Addr:         cfg.Address,
		Handler:      e,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Error("failed to start server", sl.Err(err))
		os.Exit(1)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
