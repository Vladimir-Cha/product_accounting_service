package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/Vladimir-Cha/product_accounting_service/internal/api"
	"github.com/Vladimir-Cha/product_accounting_service/internal/config"
	"github.com/Vladimir-Cha/product_accounting_service/internal/database"
	"github.com/Vladimir-Cha/product_accounting_service/internal/postgres"
	"github.com/Vladimir-Cha/product_accounting_service/internal/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.Load()

	repo, err := database.New(context.Background(), config.DBConfig{
		DBUrl:      cfg.DBUrl,
		DBMaxConns: cfg.DBMaxConns,
		DBMinConns: cfg.DBMinConns,
		DBConnLife: cfg.DBConnLife,
	})
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer repo.Close()

	if err := repo.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// опциональные миграции
	runMigrations := true
	if runMigrations {
		if err := postgres.GooseMigrationsWithPool(context.Background(), repo.Pool()); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		log.Println("Migrations applied successfully")
	}

	productStore := postgres.NewProductStore(repo.Pool())
	productHandler := api.NewProductHandler(productStore)
	categoryStore := postgres.NewCategoryStore(repo.Pool())
	categoryHandler := api.NewCategoryHandler(categoryStore)

	e := echo.New()
	e.Validator = validator.New()
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

	// запуск сервера
	addr := ":" + strconv.Itoa(cfg.ServerPort)
	if err := e.Start(addr); err != http.ErrServerClosed {
		log.Fatalf("Starting server error: %v", err)
	}
}
