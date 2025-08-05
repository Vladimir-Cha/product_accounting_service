package entities

import "time"

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name" validate:"required,min=2"`
	Price       float64   `json:"price" validate:"required,gt=0"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}
