package storage

import (
	"context"

	stderrors "errors"

	"github.com/Vladimir-Cha/product_accounting_service/internal/entities"
	"github.com/Vladimir-Cha/product_accounting_service/internal/errors"
	"github.com/Vladimir-Cha/product_accounting_service/internal/logger"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductStore struct {
	db     *pgxpool.Pool
	logger logger.Logger
}

func NewProductStore(db *pgxpool.Pool, logger logger.Logger) *ProductStore {
	return &ProductStore{db: db, logger: logger}
}

func (s *ProductStore) Create(ctx context.Context, p *entities.Product) error {
	query := `
		INSERT INTO products (name, price, description)
		VALUES ($1, $2, $3)
		RETURNING id`

	err := s.db.QueryRow(ctx, query, p.Name, p.Price, p.Description).Scan(&p.ID)
	if err != nil {
		s.logger.Error(ctx, "failed to create product",
			"error", err,
			"product_ID", p.ID,
			"product_name", p.Name,
		)
		return errors.ErrDatabase.WithError(err)
	}
	s.logger.Debug(ctx, "product created",
		"product_id", p.ID,
	)
	return nil
}

func (s *ProductStore) Read(ctx context.Context, id int) (*entities.Product, error) {
	p := &entities.Product{}

	err := s.db.QueryRow(ctx,
		"SELECT ID, name, price, description, created_at FROM products WHERE id = $1", id).
		Scan(&p.ID, &p.Name, &p.Price, &p.Description, &p.CreatedAt)

	if err != nil {
		if stderrors.Is(err, pgx.ErrNoRows) {
			s.logger.Error(ctx, "product not found",
				"product_id", id,
			)
			s.logger.Error(ctx, "database error",
				"operation", "read product",
				"product_id", id,
				"error", err,
			)
			return nil, errors.ErrNotFound.WithError(err)
		}
	}
	s.logger.Debug(ctx, "product retrieved",
		"product_id", id,
	)
	return p, nil
}

func (s *ProductStore) Update(ctx context.Context, product *entities.Product) error {
	query := `
		UPDATE products
		SET name = $1,
			price = $2,
			description = $3,
			updated_at = NOW()
		WHERE id = $4
		RETURNING id, name, price, description, created_at, updated_at`

	err := s.db.QueryRow(
		ctx,
		query,
		product.Name,
		product.Price,
		product.Description,
		product.ID,
	).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.Description,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			s.logger.Error(ctx, "product not found",
				"product_id", product.ID,
			)
			return errors.ErrNotFound
		}
		s.logger.Error(ctx, "failed to update product",
			"product_id", product.ID,
		)
		return errors.ErrBadRequest
	}
	s.logger.Info(ctx, "product updated successfully",
		"product_id", product.ID,
	)
	return nil
}

func (s *ProductStore) Delete(ctx context.Context, id int) (*entities.Product, error) {
	p, err := s.Read(ctx, id)
	if err != nil {
		if stderrors.Is(err, errors.ErrNotFound) {
			s.logger.Error(ctx, "product not found",
				"product_id", id,
			)
			return nil, errors.ErrNotFound
		}
		s.logger.Error(ctx, "failed to check product",
			"product_id", id,
			"error", err,
		)
		return nil, errors.ErrDatabase
	}

	_, err = s.db.Exec(ctx, "DELETE FROM products WHERE id = $1", id)

	if err != nil {
		s.logger.Error(ctx, "failed to delete product",
			"product_id", id,
			"error", err,
		)
		return nil, errors.ErrDatabase
	}

	s.logger.Info(ctx, "product deleted successfully",
		"product_id", id,
		"product_name", p.Name,
	)
	return p, nil
}
