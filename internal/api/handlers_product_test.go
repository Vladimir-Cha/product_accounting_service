package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Vladimir-Cha/product_accounting_service/internal/entities"
	"github.com/Vladimir-Cha/product_accounting_service/internal/errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Мок для ProductStore
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

func TestProductHandler_CreateProduct(t *testing.T) {
	// Настройка Echo
	e := echo.New()
	store := new(MockProductStore)
	handler := NewProductHandler(store)

	// успешное создание категории
	product := entities.Product{ID: 1, Name: "Test Product"}
	store.On("Create", mock.Anything, &product).Return(nil)

	body, _ := json.Marshal(product)
	req := httptest.NewRequest(http.MethodPost, "/product", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.CreateProduct(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response entities.Product
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, product, response)

	// некорректный JSON
	req = httptest.NewRequest(http.MethodPost, "/product", bytes.NewReader([]byte("invalid json")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	err = handler.CreateProduct(c)
	assert.NoError(t, err)
	assert.Equal(t, errors.ErrBadRequest.Code, rec.Code)
}

func TestProductHandler_GetProduct(t *testing.T) {
	e := echo.New()
	store := new(MockProductStore)
	handler := NewProductHandler(store)

	// успешное получение категории
	product := &entities.Product{ID: 1, Name: "Test Product"}
	store.On("Read", mock.Anything, 1).Return(product, nil)

	req := httptest.NewRequest(http.MethodGet, "/product/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := handler.GetProduct(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response entities.Product
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, *product, response)

	// некорректный ID
	req = httptest.NewRequest(http.MethodGet, "/product/invalid", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid")

	err = handler.GetProduct(c)
	assert.NoError(t, err)
	assert.Equal(t, errors.ErrBadRequest.Code, rec.Code)
}
