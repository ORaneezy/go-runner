package database

import (
	"context"
	"fmt"
	"time"

	"github.com/ORaneezy/go-runner/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgres(config *config.Config) (*pgxpool.Pool, error) {
	ctx := context.Background()
	pool, err := pgxpool.New(
		ctx,
		fmt.Sprintf(
			"postgres://%s:%s@%s/%s?sslmode=disable",
			config.DatabaseUser, config.DatabasePassword, config.DatabaseHost, config.DatabaseName,
		),
	)
	if err != nil {
		return nil, err
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = pool.Ping(timeoutCtx)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
