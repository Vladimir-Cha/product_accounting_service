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

// интерфейс CategoryStore для тестирования
type MockCategoryStore struct {
	mock.Mock
}

func (m *MockCategoryStore) CreateCat(ctx context.Context, category *entities.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryStore) ReadCat(ctx context.Context, id int) (*entities.Category, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.Category), args.Error(1)
}

func (m *MockCategoryStore) UpdateCat(ctx context.Context, category *entities.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryStore) DeleteCat(ctx context.Context, id int) (*entities.Category, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.Category), args.Error(1)
}

func TestCategoryStore_CreateCat(t *testing.T) {
	store := new(MockCategoryStore)
	ctx := context.Background()

	// успешное создание категории
	category := &entities.Category{
		Name:        "Test Category",
		Description: "Test Description",
	}
	store.On("CreateCat", ctx, category).Return(nil).Run(func(args mock.Arguments) {
		// имитируем заполнение полей
		cat := args.Get(1).(*entities.Category)
		cat.ID = 1
		cat.CreatedAt = time.Now()
	})

	err := store.CreateCat(ctx, category)
	assert.NoError(t, err)
	assert.Equal(t, 1, category.ID)
	assert.NotZero(t, category.CreatedAt)
	store.AssertCalled(t, "CreateCat", ctx, category)

	// ошибка базы данных
	store.On("CreateCat", ctx, mock.AnythingOfType("*entities.Category")).Return(errors.New("database error"))
	err = store.CreateCat(ctx, &entities.Category{Name: "Error Category"})
	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
}

func TestCategoryStore_ReadCat(t *testing.T) {
	store := new(MockCategoryStore)
	ctx := context.Background()

	// успешное чтение категории
	category := &entities.Category{
		ID:          1,
		Name:        "Test Category",
		Description: "Test Description",
		CreatedAt:   time.Now(),
	}
	store.On("ReadCat", ctx, 1).Return(category, nil)

	result, err := store.ReadCat(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, category, result)
	store.AssertCalled(t, "ReadCat", ctx, 1)

	// категория не найдена
	store.On("ReadCat", ctx, 999).Return((*entities.Category)(nil), nil)
	result, err = store.ReadCat(ctx, 999)
	assert.NoError(t, err)
	assert.Nil(t, result)
	store.AssertCalled(t, "ReadCat", ctx, 999)
}

func TestCategoryStore_UpdateCat(t *testing.T) {
	store := new(MockCategoryStore)

	// успешное обновление
	category := &entities.Category{
		ID:          1,
		Name:        "Updated Category",
		Description: "Updated Description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	store.On("UpdateCat", mock.Anything, category).Return(nil)

	err := store.UpdateCat(context.Background(), category)
	assert.NoError(t, err)
	store.AssertCalled(t, "UpdateCat", mock.Anything, category)

	// категория не найдена
	store.On("UpdateCat", mock.Anything, mock.AnythingOfType("*entities.Category")).
		Return(errors.New("category not found"))

	err = store.UpdateCat(context.Background(), &entities.Category{ID: 999})
	assert.Error(t, err)
	assert.Equal(t, "category not found", err.Error())
}

func TestCategoryStore_DeleteCat(t *testing.T) {
	store := new(MockCategoryStore)
	ctx := context.Background()

	// успешное удаление категории
	category := &entities.Category{
		ID:          1,
		Name:        "Test Category",
		Description: "Test Description",
		CreatedAt:   time.Now(),
	}
	store.On("ReadCat", ctx, 1).Return(category, nil)
	store.On("DeleteCat", ctx, 1).Return(category, nil)

	result, err := store.DeleteCat(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, category, result)
	store.AssertCalled(t, "DeleteCat", ctx, 1)

	// категория не найдена
	store.On("ReadCat", ctx, 999).Return((*entities.Category)(nil), nil)
	store.On("DeleteCat", ctx, 999).Return((*entities.Category)(nil), nil)
	result, err = store.DeleteCat(ctx, 999)
	assert.NoError(t, err)
	assert.Nil(t, result)
	store.AssertCalled(t, "DeleteCat", ctx, 999)
}
