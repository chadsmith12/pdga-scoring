package database

import (
	"context"

	"github.com/chadsmith12/pdga-scoring/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context) (*pgxpool.Pool, error) {
    dbConfig, err := config.LoadDatabase()
    if err != nil {
        return nil, err
    }
    pgConfig, err := pgxpool.ParseConfig(dbConfig.String())
    if err != nil {
        return nil, err
    }
    conn, err := pgxpool.NewWithConfig(ctx, pgConfig)

    return conn, err
}
