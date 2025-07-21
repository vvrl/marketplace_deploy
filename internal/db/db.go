package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"marketplace/internal/config"
	"marketplace/internal/logger"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

func NewDatabase(cfg *config.Config) (*sql.DB, error) {

	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Dbname,
	)

	dbConfig, err := pgx.ParseConfig(url)
	if err != nil {
		logger.Logger.Fatalf("failed to parse config for pgx: %v", err)
	}

	var lastErr error

	for i := 0; i < cfg.Database.MaxAttempts; i++ {

		db := stdlib.OpenDB(*dbConfig)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			lastErr = err
			logger.Logger.Warnf("database connection attempt %d/%d failed: %v", i+1, cfg.Database.MaxAttempts, err)
			time.Sleep(2 * time.Second)
			continue
		}
		return db, nil
	}

	return nil, fmt.Errorf("unable to connect to database: %w", lastErr)

}
