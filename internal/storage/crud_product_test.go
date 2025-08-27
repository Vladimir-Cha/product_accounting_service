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
	"github.com/stretchr/testify/require"
)

func TestProductStore_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockProductStore(ctrl)
	ctx := context.Background()

	t.Run("successful creation", func(t *testing.T) {
		product := &entities.Product{
			Name:        "Test Product",
			Price:       99.99,
			Description: "Test Description",
		}

		// Настраиваем ожидание
		mockStore.EXPECT().
			Create(ctx, product).
			DoAndReturn(func(ctx context.Context, cat *entities.Product) error {
				// Имитируем заполнение полей БД
				cat.ID = 1
				cat.CreatedAt = time.Now().Truncate(time.Microsecond)
				return nil
			})

		err := mockStore.Create(ctx, product)
		assert.NoError(t, err)
		assert.Equal(t, 1, product.ID)
		assert.NotZero(t, product.CreatedAt)
	})

	t.Run("database error", func(t *testing.T) {
		errorProduct := &entities.Product{
			Name: "Error Product",
		}

		mockStore.EXPECT().
			Create(ctx, errorProduct).
			Return(errors.New("database error"))

		err := mockStore.Create(ctx, errorProduct)
		assert.Error(t, err)
		require.Equal(t, "database error", err.Error())
	})
}

func TestProductStore_Read(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockProductStore(ctrl)
	ctx := context.Background()

	t.Run("successful read", func(t *testing.T) {
		expectedProduct := &entities.Product{
			ID:          1,
			Name:        "Test Product",
			Price:       99.99,
			Description: "Test Description",
			CreatedAt:   time.Now().Truncate(time.Microsecond),
		}

		mockStore.EXPECT().
			Read(ctx, 1).
			Return(expectedProduct, nil)

		result, err := mockStore.Read(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, expectedProduct, result)
	})

	t.Run("product not found", func(t *testing.T) {
		mockStore.EXPECT().
			Read(ctx, 999).
			Return(nil, nil)

		result, err := mockStore.Read(ctx, 999)
		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("database error on read", func(t *testing.T) {
		mockStore.EXPECT().
			Read(ctx, 500).
			Return(nil, errors.New("database connection error"))

		result, err := mockStore.Read(ctx, 500)
		assert.Error(t, err)
		assert.Nil(t, result)
		require.Equal(t, "database connection error", err.Error())
	})
}

func TestProductStore_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockProductStore(ctrl)
	ctx := context.Background()

	t.Run("successful update", func(t *testing.T) {
		product := &entities.Product{
			ID:          1,
			Name:        "Updated Product",
			Price:       99.99,
			Description: "Updated Description",
			CreatedAt:   time.Now().Add(-24 * time.Hour).Truncate(time.Microsecond),
			UpdatedAt:   time.Now().Truncate(time.Microsecond),
		}

		mockStore.EXPECT().
			Update(ctx, product).
			Return(nil)

		err := mockStore.Update(ctx, product)
		assert.NoError(t, err)
	})

	t.Run("product not found", func(t *testing.T) {
		nonExistentProduct := &entities.Product{
			ID:   999,
			Name: "Non-existent Product",
		}

		mockStore.EXPECT().
			Update(ctx, nonExistentProduct).
			Return(errors.New("product not found"))

		err := mockStore.Update(ctx, nonExistentProduct)
		assert.Error(t, err)
		require.Equal(t, "product not found", err.Error())
	})

	t.Run("database error on update", func(t *testing.T) {
		product := &entities.Product{
			ID:   1,
			Name: "Test Product",
		}

		mockStore.EXPECT().
			Update(ctx, product).
			Return(errors.New("update failed"))

		err := mockStore.Update(ctx, product)
		assert.Error(t, err)
		require.Equal(t, "update failed", err.Error())
	})
}

func TestProductStore_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockProductStore(ctrl)
	ctx := context.Background()

	t.Run("successful delete", func(t *testing.T) {
		expectedProduct := &entities.Product{
			ID:          1,
			Name:        "Test Product",
			Price:       99.99,
			Description: "Test Description",
			CreatedAt:   time.Now().Truncate(time.Microsecond),
		}

		mockStore.EXPECT().
			Delete(ctx, 1).
			Return(expectedProduct, nil)

		result, err := mockStore.Delete(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, expectedProduct, result)
	})

	t.Run("product not found", func(t *testing.T) {
		mockStore.EXPECT().
			Delete(ctx, 999).
			Return(nil, nil)

		result, err := mockStore.Delete(ctx, 999)
		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("database error on delete", func(t *testing.T) {
		mockStore.EXPECT().
			Delete(ctx, 500).
			Return(nil, errors.New("delete failed"))

		result, err := mockStore.Delete(ctx, 500)
		assert.Error(t, err)
		assert.Nil(t, result)
		require.Equal(t, "delete failed", err.Error())
	})
}
