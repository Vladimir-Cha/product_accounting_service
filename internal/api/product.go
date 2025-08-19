package api

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/Vladimir-Cha/product_accounting_service/internal/entities"
	"github.com/Vladimir-Cha/product_accounting_service/internal/errors"
	"github.com/labstack/echo/v4"
)

type ProductStore interface {
	Create(ctx context.Context, category *entities.Product) error
	Read(ctx context.Context, id int) (*entities.Product, error)
	Update(ctx context.Context, category *entities.Product) error
	Delete(ctx context.Context, id int) (*entities.Product, error)
}

type ProductHandler struct {
	store ProductStore
}

func NewProductHandler(store ProductStore) *ProductHandler {
	return &ProductHandler{store: store}
}

// хэндлер для создания записи в БД
func (h *ProductHandler) CreateProduct(c echo.Context) error {
	var p entities.Product
	if err := c.Bind(&p); err != nil {
		log.Printf("Invalid request body for product")
		return c.JSON(errors.ErrBadRequest.Code, errors.ErrBadRequest.WithMap())
	}

	if err := h.store.Create(c.Request().Context(), &p); err != nil {
		log.Printf("Failed to create product")
		return c.JSON(errors.ErrBadRequest.Code, errors.ErrBadRequest.WithMap())
	}
	return c.JSON(http.StatusCreated, p)
}

// хэндлер для получения записи по id
func (h *ProductHandler) GetProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Invalid ID format for product")
		return c.JSON(errors.ErrBadRequest.Code, errors.ErrBadRequest.WithDetails(map[string]any{
			"field": "id",
			"value": c.Param("id"),
		}).WithMap())
	}

	product, err := h.store.Read(c.Request().Context(), id)
	if err != nil {
		log.Printf("Product not found: %v", err)
		return c.JSON(errors.ErrDatabase.Code, errors.ErrDatabase.WithError(err).WithMap())
	}

	return c.JSON(http.StatusOK, product)
}

// хэндлер для обновления записи по id
func (h *ProductHandler) UpdateProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Invalid ID format for product")
		return c.JSON(errors.ErrBadRequest.Code,
			errors.ErrBadRequest.WithDetails(map[string]any{
				"field": "id",
				"value": c.Param("id"),
			}).WithMap())
	}

	var input entities.Product

	if err := c.Bind(&input); err != nil {
		log.Printf("Invalid request body for product")
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
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	err = h.store.Update(c.Request().Context(), &input)
	if err != nil {
		log.Printf("Update product failed: %v", err)
		return c.JSON(errors.ErrBadRequest.Code, errors.ErrBadRequest.WithError(err).WithMap())
	}
	return c.JSON(http.StatusOK, input)
}

// хэндлер для удаления записи по id
func (h *ProductHandler) DeleteProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Invalid ID format for product")
		return c.JSON(errors.ErrBadRequest.Code,
			errors.ErrBadRequest.WithDetails(map[string]any{
				"field": "id",
				"value": c.Param("id"),
			}).WithMap())
	}

	deletedProduct, err := h.store.Delete(c.Request().Context(), id)
	if err != nil {
		log.Printf("Delete failed for product %d: %v", id, err)
		return c.JSON(errors.ErrDatabase.Code, errors.ErrDatabase.WithError(err).WithMap())
	}

	if deletedProduct == nil {
		return c.JSON(errors.ErrNotFound.Code,
			errors.ErrNotFound.WithDetails(map[string]any{
				"id": id,
			}).WithMap())
	}

	log.Printf("Product deleted: ID=%d, Name=%s, Price=%v", deletedProduct.ID, deletedProduct.Name, deletedProduct.Price)
	return c.JSON(http.StatusOK, deletedProduct)
}
