package postgres

import (
	"context"
	"fmt"

	"github.com/Vladimir-Cha/product_accounting_service/internal/entities"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/gommon/log"
)

type CategoryStore struct {
	db *pgxpool.Pool
}

func NewCategoryStore(db *pgxpool.Pool) *CategoryStore {
	return &CategoryStore{db: db}
}

func (s *CategoryStore) CreateCat(ctx context.Context, p *entities.Category) error {
	query := `
		INSERT INTO categories (name, description)
		VALUES ($1, $2)
		RETURNING id, description, created_at`

	return s.db.QueryRow(ctx, query, p.Name, p.Description).Scan(&p.ID, &p.Description, &p.CreatedAt)
}

func (s *CategoryStore) ReadCat(ctx context.Context, id int) (*entities.Category, error) {
	p := &entities.Category{}

	err := s.db.QueryRow(ctx,
		"SELECT ID, name, description, created_at FROM categories WHERE id = $1", id).
		Scan(&p.ID, &p.Name, &p.Description, &p.CreatedAt)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return p, err
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
			log.Errorf("category with ID %d not found", category.ID)
			return fmt.Errorf("category not found")
		}
		log.Errorf("failed to update: %v", err)
		return fmt.Errorf("failed to update category")
	}
	log.Printf("category id: %d updated sucsessfully", category.ID)
	return nil
}

func (s *CategoryStore) DeleteCat(ctx context.Context, id int) (*entities.Category, error) {
	p, err := s.ReadCat(ctx, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, nil // категория не найдена
	}

	_, err = s.db.Exec(ctx,
		"DELETE FROM categories WHERE id = $1", id)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Errorf("database error: %v", err)
		return nil, err
	}

	log.Printf("category deleted: ID=%d, Name=%q", p.ID, p.Name)

	return p, nil
}
