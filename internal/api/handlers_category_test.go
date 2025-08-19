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
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Мок для CategoryStore
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

// Валидатор для Echo
type structValidator struct {
	validator *validator.Validate
}

func (v *structValidator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func TestCategoryHandler_CreateCategory(t *testing.T) {
	// Настройка Echo
	e := echo.New()
	store := new(MockCategoryStore)
	handler := NewCategoryHandler(store)

	// успешное создание категории
	category := entities.Category{ID: 1, Name: "Test Category"}
	store.On("CreateCat", mock.Anything, &category).Return(nil)

	body, _ := json.Marshal(category)
	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.CreateCategory(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response entities.Category
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, category, response)

	// некорректный JSON
	req = httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader([]byte("invalid json")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	err = handler.CreateCategory(c)
	assert.NoError(t, err)
	assert.Equal(t, errors.ErrBadRequest.Code, rec.Code)
}

func TestCategoryHandler_GetCategory(t *testing.T) {
	e := echo.New()
	store := new(MockCategoryStore)
	handler := NewCategoryHandler(store)

	// успешное получение категории
	category := &entities.Category{ID: 1, Name: "Test Category"}
	store.On("ReadCat", mock.Anything, 1).Return(category, nil)

	req := httptest.NewRequest(http.MethodGet, "/categories/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := handler.GetCategory(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response entities.Category
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, *category, response)

	// некорректный ID
	req = httptest.NewRequest(http.MethodGet, "/categories/invalid", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid")

	err = handler.GetCategory(c)
	assert.NoError(t, err)
	assert.Equal(t, errors.ErrBadRequest.Code, rec.Code)
}

func TestCategoryHandler_UpdateCategory(t *testing.T) {
	e := echo.New()
	// Настраиваем валидатор
	e.Validator = &structValidator{validator: validator.New()}
	store := new(MockCategoryStore)
	handler := NewCategoryHandler(store)

	// Успешное обновление категории
	category := &entities.Category{
		ID:          1,
		Name:        "Updated Category",
		Description: "Updated Description",
		CreatedAt:   time.Now().Truncate(time.Microsecond).UTC(),
		UpdatedAt:   time.Now().Truncate(time.Microsecond).UTC(),
	}
	store.On("UpdateCat", mock.Anything, category).Return(nil)

	body, _ := json.Marshal(category)
	req := httptest.NewRequest(http.MethodPut, "/categories/1", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := handler.UpdateCategory(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response entities.Category
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, *category, response)

	// Некорректный ID в URL
	req = httptest.NewRequest(http.MethodPut, "/categories/invalid", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid")

	err = handler.UpdateCategory(c)
	assert.NoError(t, err)
	assert.Equal(t, errors.ErrBadRequest.Code, rec.Code)

	// Некорректный JSON
	req = httptest.NewRequest(http.MethodPut, "/categories/1", bytes.NewReader([]byte("invalid json")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err = handler.UpdateCategory(c)
	assert.NoError(t, err)
	assert.Equal(t, errors.ErrBadRequest.Code, rec.Code)

	// Несовпадение ID в URL и теле
	invalidCategory := &entities.Category{
		ID:          2,
		Name:        "Invalid Category",
		Description: "Invalid Description",
	}
	body, _ = json.Marshal(invalidCategory)
	req = httptest.NewRequest(http.MethodPut, "/categories/1", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err = handler.UpdateCategory(c)
	assert.NoError(t, err)
	assert.Equal(t, errors.ErrBadRequest.Code, rec.Code)

}

func TestCategoryHandler_DeleteCategory(t *testing.T) {
	e := echo.New()
	store := new(MockCategoryStore)
	handler := NewCategoryHandler(store)

	// Успешное удаление категории
	category := &entities.Category{
		ID:          1,
		Name:        "Test Category",
		Description: "Test Description",
		CreatedAt:   time.Now(),
	}
	store.On("DeleteCat", mock.Anything, 1).Return(category, nil)

	req := httptest.NewRequest(http.MethodDelete, "/categories/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := handler.DeleteCategory(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response entities.Category
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, category.ID, response.ID)
	assert.Equal(t, category.Name, response.Name)
	assert.Equal(t, category.Description, response.Description)
	assert.Equal(t, category.CategoryID, response.CategoryID)
	assert.WithinDuration(t, category.CreatedAt, response.CreatedAt, time.Millisecond)
	assert.WithinDuration(t, category.UpdatedAt, response.UpdatedAt, time.Millisecond)

}
