package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"goshrest/internal/config"
	"goshrest/internal/gate/userprofile"
	"goshrest/internal/gate/usertoken"
	"goshrest/internal/story/signin"
	"goshrest/web"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	dbConf := config.NewDbConfig(appConf.DbConfig)
	db, err := sql.Open(dbConf.Driver, dbConf.DSN())
	if err != nil {
		slog.Error("Error connecting DB", slog.Any("error", err))
		panic(err)
	}

	defer db.Close()

	err = migrateSQL(dbConf, db)
	if err != nil {
		slog.Error("SQL migrations failed", slog.Any("error", err))
		panic(err)
	}

	logger := initLogger()

	dbx := sqlx.NewDb(db, dbConf.Driver)
	txManager := manager.Must(trmsqlx.NewDefaultFactory(dbx))
	txGetter := func(ctx context.Context) trmsqlx.Tr {
		return trmsqlx.DefaultCtxGetter.DefaultTrOrDB(ctx, dbx)
	}

	userProfileGate := userprofile.NewGate(txGetter)
	userTokenGate := usertoken.NewGate(txGetter)
	signInStory := signin.NewStory(userProfileGate, userTokenGate, txManager)

	webRoutes := web.NewWebRoutes(signInStory, config.NewGAuthConfig(appConf.GoogleAuthConfig))

	r := chi.NewRouter()
	r.Use(httplog.RequestLogger(logger))

	routes := webRoutes.Routes()
	for _, route := range routes {
		r.Method(route.Method, route.Pattern, route.Handler)
	}

	slog.Debug("Starting an app", slog.Int64("port", appConf.Port))
	http.ListenAndServe(fmt.Sprintf("localhost:%d", appConf.Port), r)
}

func migrateSQL(dbConf config.DbConfig, db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return errors.Wrap(err, "Error initting SQL migrator")
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations/sql", dbConf.Catalog, driver)
	if err != nil {
		return errors.Wrap(err, "Error locating SQL migrations")
	}

	err = m.Up()
	if err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return errors.Wrap(err, "Error executing SQL migrations")
		}
	}

	return nil
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
