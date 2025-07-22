package api

import (
	"net/http"
	"strconv"

	"github.com/Vladimir-Cha/product_accounting_service/internal/entities"
	"github.com/Vladimir-Cha/product_accounting_service/internal/postgres"
	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	store *postgres.ProductStore
}

func NewProductHandler(store *postgres.ProductStore) *ProductHandler {
	return &ProductHandler{store: store}
}

// хэндлер для создания записи в БД
func (h *ProductHandler) CreateProduct(c echo.Context) error {
	var p entities.Product
	if err := c.Bind(&p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.store.Create(c.Request().Context(), &p); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, p)
}

// хэндлер для получения записи по id
func (h *ProductHandler) GetProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}

	product, err := h.store.Read(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Product not found"})
	}

	return c.JSON(http.StatusOK, product)
}
