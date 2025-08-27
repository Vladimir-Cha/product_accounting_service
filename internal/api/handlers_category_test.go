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

// Валидатор для Echo
type structValidator struct {
	validator *validator.Validate
}

func (v *structValidator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func TestCategoryHandler_CreateCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockCategoryStore(ctrl)
	handler := NewCategoryHandler(mockStore)
	e := echo.New()

	inputCategory := entities.Category{
		Name:        "Test Category",
		Description: "Test Description",
	}

	expectedCategory := entities.Category{
		ID:          1,
		Name:        "Test Category",
		Description: "Test Description",
	}

	// Настраиваем ожидания
	mockStore.EXPECT().
		CreateCat(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, category *entities.Category) error {
			// Симулируем присвоение ID
			category.ID = 1
			return nil
		})

	body, _ := json.Marshal(inputCategory)
	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Выполняем
	err := handler.CreateCategory(c)

	// Проверяем
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response entities.Category
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, expectedCategory.ID, response.ID)
	assert.Equal(t, expectedCategory.Name, response.Name)
}

func TestCategoryHandler_CreateCategory_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockCategoryStore(ctrl)
	handler := NewCategoryHandler(mockStore)
	e := echo.New()

	// Невалидный JSON
	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader([]byte("invalid json")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.CreateCategory(c)

	assert.NoError(t, err)
	assert.Equal(t, errors.ErrBadRequest.Code, rec.Code)
}

func TestCategoryHandler_GetCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockCategoryStore(ctrl)
	handler := NewCategoryHandler(mockStore)
	e := echo.New()

	expectedCategory := &entities.Category{
		ID:          1,
		Name:        "Test Category",
		Description: "Test Description",
		CreatedAt:   time.Now().Truncate(time.Second), // Убираем наносекунды для сравнения
	}

	// Настраиваем ожидания
	mockStore.EXPECT().
		ReadCat(gomock.Any(), 1).
		Return(expectedCategory, nil)

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
	assert.Equal(t, expectedCategory.ID, response.ID)
	assert.Equal(t, expectedCategory.Name, response.Name)
}

func TestCategoryHandler_GetCategory_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockCategoryStore(ctrl)
	handler := NewCategoryHandler(mockStore)
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/categories/invalid", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid")

	err := handler.GetCategory(c)

	assert.NoError(t, err)
	assert.Equal(t, errors.ErrBadRequest.Code, rec.Code)
}

func TestCategoryHandler_UpdateCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockCategoryStore(ctrl)
	handler := NewCategoryHandler(mockStore)
	e := echo.New()
	e.Validator = &structValidator{validator: validator.New()}

	updateData := entities.Category{
		ID:          1,
		Name:        "Updated Category",
		Description: "Updated Description",
	}

	// Настраиваем ожидания
	mockStore.EXPECT().
		UpdateCat(gomock.Any(), &updateData).
		Return(nil)

	body, _ := json.Marshal(updateData)
	req := httptest.NewRequest(http.MethodPut, "/categories/1", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := handler.UpdateCategory(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCategoryHandler_DeleteCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockCategoryStore(ctrl)
	handler := NewCategoryHandler(mockStore)
	e := echo.New()

	expectedCategory := &entities.Category{
		ID:          1,
		Name:        "Test Category",
		Description: "Test Description",
		CreatedAt:   time.Now().Truncate(time.Second),
	}

	// Настраиваем ожидания
	mockStore.EXPECT().
		DeleteCat(gomock.Any(), 1).
		Return(expectedCategory, nil)

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
	assert.Equal(t, expectedCategory.ID, response.ID)
	assert.Equal(t, expectedCategory.Name, response.Name)
}
