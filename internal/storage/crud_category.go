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

type CategoryStore struct {
	db     *pgxpool.Pool
	logger logger.Logger
}

func NewCategoryStore(db *pgxpool.Pool, logger logger.Logger) *CategoryStore {
	return &CategoryStore{db: db, logger: logger}
}

func (s *CategoryStore) CreateCat(ctx context.Context, p *entities.Category) error {
	query := `
		INSERT INTO categories (name, description)
		VALUES ($1, $2)
		RETURNING id`

	err := s.db.QueryRow(ctx, query, p.Name, p.Description).Scan(&p.ID)
	if err != nil {
		s.logger.Error(ctx, "failed to create category",
			"error", err,
			"category_ID", p.ID,
			"category_name", p.Name,
		)
		return errors.ErrDatabase.WithError(err)
	}
	s.logger.Debug(ctx, "category created",
		"category_id", p.ID,
	)
	return nil
}

func (s *CategoryStore) ReadCat(ctx context.Context, id int) (*entities.Category, error) {
	p := &entities.Category{}

	err := s.db.QueryRow(ctx,
		"SELECT ID, name, description, created_at FROM categories WHERE id = $1", id).
		Scan(&p.ID, &p.Name, &p.Description, &p.CreatedAt)

	if err != nil {
		if stderrors.Is(err, pgx.ErrNoRows) {
			s.logger.Error(ctx, "category not found",
				"category_id", id,
			)
			return nil, errors.ErrNotFound.WithError(err)
		}
	}
	s.logger.Debug(ctx, "category retrieved",
		"category_id", id,
	)
	return p, nil
}

func (s *CategoryStore) UpdateCat(ctx context.Context, category *entities.Category) error {
	query := `
		UPDATE categories
		SET name = $1,
			description = $2,
			updated_at = NOW()
		WHERE id = $3
		RETURNING id, name, description, created_at, updated_at`

	err := s.db.QueryRow(
		ctx,
		query,
		category.Name,
		category.Description,
		category.ID,
	).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			s.logger.Error(ctx, "category not found",
				"category_id", category.ID,
			)
			return errors.ErrNotFound
		}
		s.logger.Error(ctx, "failed to update category",
			"category_id", category.ID,
		)
		return errors.ErrBadRequest
	}
	s.logger.Info(ctx, "category updated successfully",
		"category_id", category.ID,
	)
	return nil
}

func (s *CategoryStore) DeleteCat(ctx context.Context, id int) (*entities.Category, error) {
	p, err := s.ReadCat(ctx, id)
	if err != nil {
		if stderrors.Is(err, errors.ErrNotFound) {
			s.logger.Error(ctx, "category not found",
				"category_id", id,
			)
			return nil, errors.ErrNotFound
		}
		s.logger.Error(ctx, "failed to check category",
			"category_id", id,
			"error", err,
		)
		return nil, errors.ErrDatabase
	}

	_, err = s.db.Exec(ctx,
		"DELETE FROM categories WHERE id = $1", id)

	if err != nil {
		s.logger.Error(ctx, "failed to delete category",
			"category_id", id,
			"error", err,
		)
		return nil, errors.ErrDatabase
	}

	s.logger.Info(ctx, "category deleted successfully",
		"category_id", id,
		"category_name", p.Name,
	)
	return p, nil
}
