package database

import (
	"context"

	"github.com/Vladimir-Cha/product_accounting_service/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/gommon/log"
)

// Repository определяет контракт для работы с БД
type Repository interface {
	Pool() *pgxpool.Pool            // пул соединений
	Ping(ctx context.Context) error // проверка соединения
	Close()                         // закрытие соединения
}

type repo struct {
	pool *pgxpool.Pool
}

// подключение к БД
func New(ctx context.Context, cfg config.DBConfig) (Repository, error) {
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

	// проверяем подключение
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	log.Printf("PostgreSQL connected (Max: %d, Min: %d)", cfg.DBMaxConns, cfg.DBMinConns)

	return &repo{pool: pool}, nil
}

func (r *repo) Pool() *pgxpool.Pool {
	return r.pool
}

func (r *repo) Ping(ctx context.Context) error {
	return r.pool.Ping(ctx)
}

func (r *repo) Close() {
	r.pool.Close()
}
