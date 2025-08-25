package storage

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Vladimir-Cha/product_accounting_service/internal/entities"
	"github.com/Vladimir-Cha/product_accounting_service/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCategoryStore_CreateCat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockCategoryStore(ctrl)
	ctx := context.Background()

	t.Run("successful creation", func(t *testing.T) {
		category := &entities.Category{
			Name:        "Test Category",
			Description: "Test Description",
		}

		// Настраиваем ожидание
		mockStore.EXPECT().
			CreateCat(ctx, category).
			DoAndReturn(func(ctx context.Context, cat *entities.Category) error {
				// Имитируем заполнение полей БД
				cat.ID = 1
				cat.CreatedAt = time.Now().Truncate(time.Microsecond)
				return nil
			})

		err := mockStore.CreateCat(ctx, category)
		assert.NoError(t, err)
		assert.Equal(t, 1, category.ID)
		assert.NotZero(t, category.CreatedAt)
	})

	t.Run("database error", func(t *testing.T) {
		errorCategory := &entities.Category{
			Name: "Error Category",
		}

		mockStore.EXPECT().
			CreateCat(ctx, errorCategory).
			Return(errors.New("database error"))

		err := mockStore.CreateCat(ctx, errorCategory)
		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
	})
}

func TestCategoryStore_ReadCat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockCategoryStore(ctrl)
	ctx := context.Background()

	t.Run("successful read", func(t *testing.T) {
		expectedCategory := &entities.Category{
			ID:          1,
			Name:        "Test Category",
			Description: "Test Description",
			CreatedAt:   time.Now().Truncate(time.Microsecond),
		}

		mockStore.EXPECT().
			ReadCat(ctx, 1).
			Return(expectedCategory, nil)

		result, err := mockStore.ReadCat(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, expectedCategory, result)
	})

	t.Run("category not found", func(t *testing.T) {
		mockStore.EXPECT().
			ReadCat(ctx, 999).
			Return(nil, nil)

		result, err := mockStore.ReadCat(ctx, 999)
		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("database error on read", func(t *testing.T) {
		mockStore.EXPECT().
			ReadCat(ctx, 500).
			Return(nil, errors.New("database connection error"))

		result, err := mockStore.ReadCat(ctx, 500)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "database connection error", err.Error())
	})
}

func TestCategoryStore_UpdateCat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockCategoryStore(ctrl)
	ctx := context.Background()

	t.Run("successful update", func(t *testing.T) {
		category := &entities.Category{
			ID:          1,
			Name:        "Updated Category",
			Description: "Updated Description",
			CreatedAt:   time.Now().Add(-24 * time.Hour).Truncate(time.Microsecond),
			UpdatedAt:   time.Now().Truncate(time.Microsecond),
		}

		mockStore.EXPECT().
			UpdateCat(ctx, category).
			Return(nil)

		err := mockStore.UpdateCat(ctx, category)
		assert.NoError(t, err)
	})

	t.Run("category not found", func(t *testing.T) {
		nonExistentCategory := &entities.Category{
			ID:   999,
			Name: "Non-existent Category",
		}

		mockStore.EXPECT().
			UpdateCat(ctx, nonExistentCategory).
			Return(errors.New("category not found"))

		err := mockStore.UpdateCat(ctx, nonExistentCategory)
		assert.Error(t, err)
		assert.Equal(t, "category not found", err.Error())
	})

	t.Run("database error on update", func(t *testing.T) {
		category := &entities.Category{
			ID:   1,
			Name: "Test Category",
		}

		mockStore.EXPECT().
			UpdateCat(ctx, category).
			Return(errors.New("update failed"))

		err := mockStore.UpdateCat(ctx, category)
		assert.Error(t, err)
		assert.Equal(t, "update failed", err.Error())
	})
}

func TestCategoryStore_DeleteCat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockCategoryStore(ctrl)
	ctx := context.Background()

	t.Run("successful delete", func(t *testing.T) {
		expectedCategory := &entities.Category{
			ID:          1,
			Name:        "Test Category",
			Description: "Test Description",
			CreatedAt:   time.Now().Truncate(time.Microsecond),
		}

		mockStore.EXPECT().
			DeleteCat(ctx, 1).
			Return(expectedCategory, nil)

		result, err := mockStore.DeleteCat(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, expectedCategory, result)
	})

	t.Run("category not found", func(t *testing.T) {
		mockStore.EXPECT().
			DeleteCat(ctx, 999).
			Return(nil, nil)

		result, err := mockStore.DeleteCat(ctx, 999)
		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("database error on delete", func(t *testing.T) {
		mockStore.EXPECT().
			DeleteCat(ctx, 500).
			Return(nil, errors.New("delete failed"))

		result, err := mockStore.DeleteCat(ctx, 500)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "delete failed", err.Error())
	})
}
