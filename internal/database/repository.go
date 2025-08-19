package database

import (
	"context"

	"github.com/Vladimir-Cha/product_accounting_service/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/gommon/log"
)

// New создает и возвращает *pgxpool.Pool напрямую
func New(ctx context.Context, cfg config.DBConfig) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DBUrl)
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = cfg.DBMaxConns
	poolConfig.MinConns = cfg.DBMinConns
	poolConfig.MaxConnLifetime = cfg.DBConnLife

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	// Проверка подключения
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	log.Printf("PostgreSQL connected (Max: %d, Min: %d)", cfg.DBMaxConns, cfg.DBMinConns)
	return pool, nil
}
