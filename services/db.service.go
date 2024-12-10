package services

import (
	"rudder/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDbServices(c config.AppConfig, p *pgxpool.Pool) *DbServices {
	return &DbServices{
		Config: c,
		Pool:   p,
	}
}

type DbServices struct {
	Config config.AppConfig
	Pool   *pgxpool.Pool
}
