package postgres

import (
	"context"
	"fmt"

	"github.com/Vladimir-Cha/product_accounting_service/internal/entities"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/gommon/log"
)

type ProductStore struct {
	db *pgxpool.Pool
}

func NewProductStore(db *pgxpool.Pool) *ProductStore {
	return &ProductStore{db: db}
}

func (s *ProductStore) Create(ctx context.Context, p *entities.Product) error {
	query := `
		INSERT INTO products (name, price, description)
		VALUES ($1, $2, $3)
		RETURNING id, name, price, description, created_at`

	return s.db.QueryRow(ctx, query, p.Name, p.Price, p.Description).Scan(&p.ID, &p.Name, &p.Price, &p.Description, &p.CreatedAt)
}

func (s *ProductStore) Read(ctx context.Context, id int) (*entities.Product, error) {
	p := &entities.Product{}

	err := s.db.QueryRow(ctx,
		"SELECT ID, name, price, description, created_at FROM products WHERE id = $1", id).
		Scan(&p.ID, &p.Name, &p.Price, &p.Description, &p.CreatedAt)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return p, err
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
			log.Errorf("product with ID %d not found", product.ID)
			return fmt.Errorf("product not found")
		}
		log.Errorf("failed to update: %v", err)
		return fmt.Errorf("failed to update product")
	}
	log.Printf("product id: %d updated sucsessfully", product.ID)
	return nil
}

func (s *ProductStore) Delete(ctx context.Context, id int) (*entities.Product, error) {
	p, err := s.Read(ctx, id)
	//p := &entities.Product{}
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, nil // Продукт не найден
	}

	_, err = s.db.Exec(ctx,
		"DELETE FROM products WHERE id = $1", id)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Errorf("database error: %v", err)
		return nil, err
	}

	return p, nil
}
