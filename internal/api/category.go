package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Vladimir-Cha/product_accounting_service/internal/entities"
	"github.com/Vladimir-Cha/product_accounting_service/internal/errors"
	"github.com/Vladimir-Cha/product_accounting_service/internal/postgres"
	"github.com/labstack/echo/v4"
)

type CategoryHandler struct {
	store *postgres.CategoryStore
}

func NewCategoryHandler(store *postgres.CategoryStore) *CategoryHandler {
	return &CategoryHandler{store: store}
}

// хэндлер для создания записи в БД
func (h *CategoryHandler) CreateCategory(c echo.Context) error {
	var p entities.Category
	if err := c.Bind(&p); err != nil {
		log.Printf("Invalid request body for category")
		return c.JSON(errors.ErrBadRequest.Code, errors.ErrBadRequest.WithMap())
	}

	if err := h.store.CreateCat(c.Request().Context(), &p); err != nil {
		log.Printf("Failed to create category")
		return c.JSON(errors.ErrDatabase.Code, errors.ErrBadRequest.WithMap())
	}
	return c.JSON(http.StatusCreated, p)
}

// хэндлер для получения записи по id
func (h *CategoryHandler) GetCategory(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Invalid ID format for category")
		return c.JSON(errors.ErrBadRequest.Code, errors.ErrBadRequest.WithDetails(map[string]any{
			"field": "id",
			"value": c.Param("id"),
		}).WithMap())
	}

	category, err := h.store.ReadCat(c.Request().Context(), id)
	if err != nil {
		log.Printf("Category not found: %v", err)
		return c.JSON(errors.ErrDatabase.Code, errors.ErrDatabase.WithError(err).WithMap())
	}

	return c.JSON(http.StatusOK, category)
}

// хэндлер для обновления записи по id
func (h *CategoryHandler) UpdateCategory(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Invalid category ID format for category")
		return c.JSON(errors.ErrBadRequest.Code,
			errors.ErrBadRequest.WithDetails(map[string]any{
				"field": "id",
				"value": c.Param("id"),
			}).WithMap())
	}

	var input entities.Category

	if err := c.Bind(&input); err != nil {
		log.Printf("Invalid request body for category")
		return c.JSON(errors.ErrBadRequest.Code,
			errors.ErrBadRequest.WithError(err).WithMap())
	}

	// проверка на совпадение id в URL и id в теле запроса
	if input.ID != 0 && input.ID != id {
		return c.JSON(errors.ErrBadRequest.Code,
			errors.ErrBadRequest.WithDetails("ID in URL and body mismatch").WithMap())
	}
	input.ID = id // Устанавливаем ID из URL

	if err := c.Validate(input); err != nil {
		return c.JSON(errors.ErrValidation.Code,
			errors.ErrValidation.WithDetails(err.Error()).WithMap())
	}

	err = h.store.UpdateCat(c.Request().Context(), &input)
	if err != nil {
		log.Printf("Update category failed: %v", err)
		return c.JSON(errors.ErrDatabase.Code, errors.ErrDatabase.WithError(err).WithMap())
	}
	return c.JSON(http.StatusOK, input)
}

// хэндлер для удаления записи по id
func (h *CategoryHandler) DeleteCategory(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Invalid ID format for category")
		return c.JSON(errors.ErrBadRequest.Code,
			errors.ErrBadRequest.WithDetails(map[string]any{
				"field": "id",
				"value": c.Param("id"),
			}).WithMap())
	}

	deletedCategory, err := h.store.DeleteCat(c.Request().Context(), id)
	if err != nil {
		log.Printf("Delete failed for category %d: %v", id, err)
		return c.JSON(errors.ErrDatabase.Code, errors.ErrDatabase.WithError(err).WithMap())
	}

	if deletedCategory == nil {
		return c.JSON(errors.ErrNotFound.Code,
			errors.ErrNotFound.WithDetails(map[string]any{
				"id": id,
			}).WithMap())
	}

	log.Printf("Product deleted: ID=%d, Name=%s", deletedCategory.ID, deletedCategory.Name)
	return c.JSON(http.StatusOK, deletedCategory)
}
