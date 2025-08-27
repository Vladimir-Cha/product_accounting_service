package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Vladimir-Cha/product_accounting_service/internal/entities"
	"github.com/Vladimir-Cha/product_accounting_service/internal/errors"
	"github.com/Vladimir-Cha/product_accounting_service/internal/mocks"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestProductHandler_CreateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProduct := mocks.NewMockProductStore(ctrl)
	handler := NewProductHandler(mockProduct)
	e := echo.New()

	inputProduct := entities.Product{
		Name:        "Test Product",
		Price:       99.99,
		Description: "Test Description",
	}

	expectedProduct := entities.Product{
		ID:          1,
		Name:        "Test Product",
		Price:       99.99,
		Description: "Test Description",
	}

	// Настраиваем ожидания
	mockProduct.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, category *entities.Product) error {
			// Симулируем присвоение ID
			category.ID = 1
			return nil
		})

	body, _ := json.Marshal(inputProduct)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Выполняем
	err := handler.CreateProduct(c)

	// Проверяем
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response entities.Product
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, expectedProduct.ID, response.ID)
	assert.Equal(t, expectedProduct.Name, response.Name)
}

func TestProductHandler_CreateProduct_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockProductStore(ctrl)
	handler := NewProductHandler(mockStore)
	e := echo.New()

	// Невалидный JSON
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader([]byte("invalid json")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.CreateProduct(c)

	assert.NoError(t, err)
	assert.Equal(t, errors.ErrBadRequest.Code, rec.Code)
}

func TestProductHandler_GetProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockProductStore(ctrl)
	handler := NewProductHandler(mockStore)
	e := echo.New()

	expectedProduct := &entities.Product{
		ID:          1,
		Name:        "Test Product",
		Price:       99.99,
		Description: "Test Description",
		CreatedAt:   time.Now().Truncate(time.Second), // Убираем наносекунды для сравнения
	}

	// Настраиваем ожидания
	mockStore.EXPECT().
		Read(gomock.Any(), 1).
		Return(expectedProduct, nil)

	req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := handler.GetProduct(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response entities.Product
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, expectedProduct.ID, response.ID)
	assert.Equal(t, expectedProduct.Name, response.Name)
}

func TestProductHandler_GetProduct_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockProductStore(ctrl)
	handler := NewProductHandler(mockStore)
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/products/invalid", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid")

	err := handler.GetProduct(c)

	assert.NoError(t, err)
	assert.Equal(t, errors.ErrBadRequest.Code, rec.Code)
}

func TestProductHandler_UpdateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockProductStore(ctrl)
	handler := NewProductHandler(mockStore)
	e := echo.New()
	e.Validator = &structValidator{validator: validator.New()}

	updateData := entities.Product{
		ID:          1,
		Name:        "Updated Product",
		Price:       99.99,
		Description: "Updated Description",
	}

	// Настраиваем ожидания
	mockStore.EXPECT().
		Update(gomock.Any(), &updateData).
		Return(nil)

	body, _ := json.Marshal(updateData)
	req := httptest.NewRequest(http.MethodPut, "/products/1", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := handler.UpdateProduct(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestProductHandler_DeleteProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockProductStore(ctrl)
	handler := NewProductHandler(mockStore)
	e := echo.New()

	expectedProduct := &entities.Product{
		ID:          1,
		Name:        "Test Product",
		Price:       99.99,
		Description: "Test Description",
		CreatedAt:   time.Now().Truncate(time.Second),
	}

	// Настраиваем ожидания
	mockStore.EXPECT().
		Delete(gomock.Any(), 1).
		Return(expectedProduct, nil)

	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := handler.DeleteProduct(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response entities.Product
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, expectedProduct.ID, response.ID)
	assert.Equal(t, expectedProduct.Name, response.Name)
}
