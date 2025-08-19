package storage

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Vladimir-Cha/product_accounting_service/internal/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// интерфейс ProductStore для тестирования
type MockProductStore struct {
	mock.Mock
}

func (m *MockProductStore) Create(ctx context.Context, product *entities.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductStore) Read(ctx context.Context, id int) (*entities.Product, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.Product), args.Error(1)
}

func (m *MockProductStore) Update(ctx context.Context, product *entities.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductStore) Delete(ctx context.Context, id int) (*entities.Product, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.Product), args.Error(1)
}

func TestProductStore_Create(t *testing.T) {
	store := new(MockProductStore)
	ctx := context.Background()

	// успешное создание продукта
	product := &entities.Product{
		Name:        "Test Product",
		Description: "Test Description",
	}
	store.On("Create", ctx, product).Return(nil).Run(func(args mock.Arguments) {
		// имитируем заполнение полей
		prod := args.Get(1).(*entities.Product)
		prod.ID = 1
		prod.CreatedAt = time.Now()
	})

	err := store.Create(ctx, product)
	assert.NoError(t, err)
	assert.Equal(t, 1, product.ID)
	assert.NotZero(t, product.CreatedAt)
	store.AssertCalled(t, "Create", ctx, product)

	// ошибка базы данных
	store.On("Create", ctx, mock.AnythingOfType("*entities.Product")).Return(errors.New("database error"))
	err = store.Create(ctx, &entities.Product{Name: "Error Product"})
	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
}

func TestProductStore_Read(t *testing.T) {
	store := new(MockProductStore)
	ctx := context.Background()

	// успешное чтение продукта
	product := &entities.Product{
		ID:          1,
		Name:        "Test Product",
		Description: "Test Description",
		CreatedAt:   time.Now(),
	}
	store.On("Read", ctx, 1).Return(product, nil)

	result, err := store.Read(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, product, result)
	store.AssertCalled(t, "Read", ctx, 1)

	// продукт не найден
	store.On("Read", ctx, 999).Return((*entities.Product)(nil), nil)
	result, err = store.Read(ctx, 999)
	assert.NoError(t, err)
	assert.Nil(t, result)
	store.AssertCalled(t, "Read", ctx, 999)
}

func TestProductStore_Update(t *testing.T) {
	store := new(MockProductStore)

	// успешное обновление
	product := &entities.Product{
		ID:          1,
		Name:        "Updated Product",
		Description: "Updated Description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	store.On("Update", mock.Anything, product).Return(nil)

	err := store.Update(context.Background(), product)
	assert.NoError(t, err)
	store.AssertCalled(t, "Update", mock.Anything, product)

	// продукт не найден
	store.On("Update", mock.Anything, mock.AnythingOfType("*entities.Product")).
		Return(errors.New("product not found"))

	err = store.Update(context.Background(), &entities.Product{ID: 999})
	assert.Error(t, err)
	assert.Equal(t, "product not found", err.Error())
}

func TestProductStore_Delete(t *testing.T) {
	store := new(MockProductStore)
	ctx := context.Background()

	// успешное удаление продукта
	product := &entities.Product{
		ID:          1,
		Name:        "Test Product",
		Description: "Test Description",
		CreatedAt:   time.Now(),
	}
	store.On("Read", ctx, 1).Return(product, nil)
	store.On("Delete", ctx, 1).Return(product, nil)

	result, err := store.Delete(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, product, result)
	store.AssertCalled(t, "Delete", ctx, 1)

	// продукт не найден
	store.On("ReadCat", ctx, 999).Return((*entities.Product)(nil), nil)
	store.On("Delete", ctx, 999).Return((*entities.Product)(nil), nil)
	result, err = store.Delete(ctx, 999)
	assert.NoError(t, err)
	assert.Nil(t, result)
	store.AssertCalled(t, "Delete", ctx, 999)
}
