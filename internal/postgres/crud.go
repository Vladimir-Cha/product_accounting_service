package postgres

import (
	"context"
	"database/sql"

	"github.com/Vladimir-Cha/product_accounting_service/internal/entities"
	"github.com/labstack/gommon/log"
)

type ProductStore struct {
	db *sql.DB
}

func NewProductStore(db *sql.DB) *ProductStore {
	return &ProductStore{db: db}
}

func (s *ProductStore) Create(ctx context.Context, p *entities.Product) error {
	query := `
INSERT INTO products (name, price)
VALUES ($1, $2)
RETURNING id, created_at`

	return s.db.QueryRowContext(ctx, query, p.Name, p.Price).Scan(&p.ID, &p.CreatedAt)
}

func (s *ProductStore) Read(ctx context.Context, id int) (*entities.Product, error) {
	p := &entities.Product{}

	err := s.db.QueryRowContext(ctx,
		"SELECT ID, name, price, created_at FROM products WHERE id = $1", id).
		Scan(&p.ID, &p.Name, &p.Price, &p.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return p, err
}

func (s *ProductStore) Update(ctx context.Context, product *entities.Product) error {
	query := `
UPDATE products
SET name = $1,
price = $2,
updated_at = $3
RETURNING id, name, price, created_at, updated_at`

	err := s.db.QueryRowContext(
		ctx,
		query,
		product.Name,
		product.Price,
		product.ID,
	).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Errorf("product with ID %d not found", product.ID)
			return nil
		}
		log.Errorf("failed to update: %v", err)
		return nil
	}
	log.Printf("product id: %d updated sucsessfully", product.ID)
	return nil
}

func (s *ProductStore) Delete(ctx context.Context, id int) (*entities.Product, error) {
	p := &entities.Product{}

	err := s.db.QueryRowContext(ctx,
		"DELETE FROM products WHERE id = $1", id).
		Scan(&p.ID, &p.Name, &p.Price, &p.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return p, err
}
