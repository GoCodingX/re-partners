package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"

	"github.com/GoCodingX/repartners/internal/config"
	"github.com/GoCodingX/repartners/internal/handlers"
	"github.com/GoCodingX/repartners/internal/repository/pg"
	"github.com/GoCodingX/repartners/pkg/gen/openapi"
	"github.com/GoCodingX/repartners/pkg/logger"
	"github.com/GoCodingX/repartners/pkg/migrate"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func main() {
	// setup logger
	logger.Init()

	// load env vars from .env files
	err := godotenv.Load()
	if err != nil {
		var pathErr *fs.PathError
		if !errors.As(err, &pathErr) {
			logger.Fatal("failed to read env file", err)

			return
		}
	}

	// load config
	cfg, err := env.ParseAs[config.Config]()
	if err != nil {
		logger.Fatal("failed to load config", err)

		return
	}

	// initialize db repo
	repo, err := initializeRepository(&cfg)
	if err != nil {
		logger.Fatal("failed to initialize repository", err)

		return
	}

	swagger, err := openapi.GetSwagger()
	if err != nil {
		logger.Fatal("error loading swagger spec", err)

		return
	}

	// initialize service
	service := handlers.NewPacksService(&handlers.NewPacksServiceParams{
		Repo: repo,
	})

	// initialize router
	srv, err := handlers.NewRouter(service, swagger, &cfg)
	if err != nil {
		logger.Fatal("failed to initialize router", err)

		return
	}

	// start the http server
	logger.Info("starting server", slog.String("port", cfg.Port))

	if err := srv.Start(fmt.Sprintf(":%s", cfg.Port)); err != nil {
		logger.Fatal("failed to start http server", err)
	}
}

func initializeRepository(cfg *config.Config) (*pg.Repository, error) {
	dbConn := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.DatabaseUrl)))

	if err := dbConn.Ping(); err != nil {
		return nil, fmt.Errorf("failed connect to db: %w", err)
	}

	dbClient := bun.NewDB(dbConn, pgdialect.New())

	if err := migrate.Up(cfg.MigrationsDir, cfg.DatabaseUrl); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	repo := pg.NewRepository(dbClient)

	return repo, nil
}
