package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/Vladimir-Cha/product_accounting_service/internal/api"
	"github.com/Vladimir-Cha/product_accounting_service/internal/config"
	"github.com/Vladimir-Cha/product_accounting_service/internal/logger"
	"github.com/Vladimir-Cha/product_accounting_service/internal/storage"
	"github.com/Vladimir-Cha/product_accounting_service/internal/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	pool, err := storage.InitNew(ctx, *cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// опциональные миграции
	if cfg.RunMigrations {
		if err := storage.GooseMigrationsWithPool(ctx, pool); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		log.Println("Migrations applied successfully")
	}

	appLogger := logger.NewConsoleLogger()
	productStore := storage.NewProductStore(pool, appLogger)
	productHandler := api.NewProductHandler(productStore)
	categoryStore := storage.NewCategoryStore(pool, appLogger)
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
	go func() {
		addr := ":" + strconv.Itoa(cfg.Server.Port)
		log.Printf("Server starting on %s", addr)
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Starting server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server shutdown error:", err)
	}
	log.Println("Server shutdown completed gracefully")
}
