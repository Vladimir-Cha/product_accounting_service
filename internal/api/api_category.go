package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Vladimir-Cha/product_accounting_service/internal/entities"
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
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.store.CreateCat(c.Request().Context(), &p); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to create category",
		})
	}
	return c.JSON(http.StatusCreated, p)
}

// хэндлер для получения записи по id
func (h *CategoryHandler) GetCategory(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid ID format",
		})
	}

	category, err := h.store.ReadCat(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Category not found",
		})
	}

	return c.JSON(http.StatusOK, category)
}

// хэндлер для обновления записи по id
func (h *CategoryHandler) UpdateCategory(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid ID format",
		})
	}

	var input entities.Category

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid request body",
		})
	}

	// проверка на совпадение id в URL и id в теле запроса
	if input.ID != 0 && input.ID != id {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "ID in URL and body mismatch",
		})
	}
	input.ID = id // Устанавливаем ID из URL

	if err := c.Validate(input); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	err = h.store.UpdateCat(c.Request().Context(), &input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to update category",
		})
	}
	return c.JSON(http.StatusOK, input)
}

// хэндлер для удаления записи по id
func (h *CategoryHandler) DeleteCategory(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid ID format",
		})
	}

	deletedCategory, err := h.store.DeleteCat(c.Request().Context(), id)
	if err != nil {
		log.Printf("Delete failed for category %d: %v", id, err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to delete category",
		})
	}

	if deletedCategory == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Category not found",
		})
	}

	log.Printf("Product deleted: ID=%d, Name=%s, Price=%v", deletedCategory.ID, deletedCategory.Name)
	return c.JSON(http.StatusOK, deletedCategory)
}
