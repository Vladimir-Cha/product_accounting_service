package entities

import "time"

type Category struct {
	ID          int       `json:"id"`
	Name        string    `json:"name" validate:"required,min=2"`
	Description string    `json:"description"`
	CategoryID  int       `json:"category_id"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}
