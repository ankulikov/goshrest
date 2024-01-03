package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"goshrest/internal/config"
	"goshrest/web"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file", slog.Any("error", err))
		panic(err)
	}

	// parse configs
	appConf, err := config.NewAppConfigFromEnv(context.Background())
	if err != nil {
		slog.Error("Error parsing env config", slog.Any("error", err))
		panic(err)
	}

	fmt.Println(appConf)

	logger := initLogger()

	webRoutes := web.NewWebRoutes(config.NewGAuthConfig(appConf.GoogleAuthConfig))

	r := chi.NewRouter()
	r.Use(httplog.RequestLogger(logger))

	routes := webRoutes.Routes()
	for _, route := range routes {
		r.Method(route.Method, route.Pattern, route.Handler)
	}

	slog.Debug("Starting an app", slog.Int64("port", appConf.Port))
	http.ListenAndServe(fmt.Sprintf("localhost:%d", appConf.Port), r)
}

func initLogger() *httplog.Logger {
	logger := httplog.NewLogger("httplog-example", httplog.Options{
		// Logger
		JSON: true,
		// TimeFieldFormat: time.RFC850,
		// SourceFieldName: "source",
		LogLevel:         slog.LevelDebug,
		Concise:          true,
		RequestHeaders:   true,
		MessageFieldName: "message",

		Tags: map[string]string{
			"version": "v1.0-81aa4244d9fc8076a",
			"env":     "dev",
		},
		QuietDownRoutes: []string{
			"/",
			"/ping",
		},
		QuietDownPeriod: 10 * time.Second,
	})

	return logger
}
