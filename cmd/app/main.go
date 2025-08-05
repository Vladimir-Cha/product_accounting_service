package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Vladimir-Cha/product_accounting_service/internal/api"
	"github.com/Vladimir-Cha/product_accounting_service/internal/postgres"
	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type customValidator struct {
	validator *validator.Validate
}

func (v *customValidator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func main() {
	//читаем .env
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found: %v", err)
	}
	// Инициализация pgxpool
	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal(err)
	}

	if err := api.GooseMigrationsWithPool(context.Background(), pool); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migrations applied successfully")

	productStore := postgres.NewProductStore(pool)
	productHandler := api.NewProductHandler(productStore)
	categoryStore := postgres.NewCategoryStore(pool)
	categoryHandler := api.NewCategoryHandler(categoryStore)

	e := echo.New()
	e.Validator = &customValidator{validator: validator.New()}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/products", productHandler.CreateProduct)
	e.GET("/products/:id", productHandler.GetProduct)
	e.PUT("/products/:id", productHandler.UpdateProduct)
	e.DELETE("/products/:id", productHandler.DeleteProduct)
	e.POST("/categories", categoryHandler.CreateCategory)
	e.GET("/categories/:id", categoryHandler.GetCategory)
	e.PUT("/categories/:id", categoryHandler.UpdateCategory)
	e.DELETE("/categories/:id", categoryHandler.DeleteCategory)

	// Запуск сервера
	if err := e.Start(":8080"); err != http.ErrServerClosed {
		log.Fatal(err)
	}

}
