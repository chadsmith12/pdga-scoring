package database

import (
	"context"
	"fmt"

	"github.com/chadsmith12/pdga-scoring/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var numberConnections = 0

func Connect(ctx context.Context) (*pgxpool.Pool, error) {
    dbConfig, err := config.LoadDatabase()
    if err != nil {
        return nil, err
    }
    pgConfig, err := pgxpool.ParseConfig(dbConfig.String())
    if err != nil {
        return nil, err
    }
    pgConfig.BeforeClose = func(c *pgx.Conn) {
	numberConnections--
	fmt.Printf("CONNECTION HAS BEEN CLOSED: %d\n", numberConnections)
    }
    pgConfig.BeforeConnect = func(ctx context.Context, cc *pgx.ConnConfig) error {
	numberConnections++
	fmt.Printf("CONNECTION IS HAPPENING: %d\n", numberConnections)

	return nil
    }
    conn, err := pgxpool.NewWithConfig(ctx, pgConfig)
    stat := conn.Stat()
    fmt.Printf("STAT MAX CONN: %d\n", stat.MaxConns())
    return conn, err
}
