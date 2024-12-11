package database

import (
	"context"
	"fmt"
	"rudder/internal/config"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBConnection struct {
	Pool *pgxpool.Pool
}

func NewDBConnection(appConfig *config.AppConfig) (*DBConnection, error) {
	poolConfig, err := pgxpool.ParseConfig(appConfig.DatabaseUrl)
	if err != nil {
		return nil, fmt.Errorf("unable to parse db config: %v", err)
	}
	poolConfig.AfterConnect = func(ctx context.Context, c *pgx.Conn) error {
		pgxdecimal.Register(c.TypeMap())
		return nil
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	return &DBConnection{Pool: pool}, nil
}

func (dc *DBConnection) Close() {
	if dc.Pool != nil {
		dc.Pool.Close()
	}
}
