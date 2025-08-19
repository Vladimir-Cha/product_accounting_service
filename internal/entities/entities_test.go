package entities

import (
	"testing"

	"github.com/Vladimir-Cha/product_accounting_service/internal/validator"
	"github.com/stretchr/testify/assert"
)

func TestProductValidation(t *testing.T) {
	p := Product{Name: "A", Price: -1} // Невалидные данные
	err := validator.New().Validate(p)
	assert.Error(t, err)
}

func TestCategoryValidation(t *testing.T) {
	p := Category{Name: "A", CategoryID: -1} // Невалидные данные
	err := validator.New().Validate(p)
	assert.Error(t, err)
}
