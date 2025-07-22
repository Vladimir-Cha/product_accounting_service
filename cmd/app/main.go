package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Vladimir-Cha/product_accounting_service/internal/api"
	"github.com/Vladimir-Cha/product_accounting_service/internal/postgres"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func gooseMigrations(db *sql.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}
	return goose.Up(db, "migrations") // Применение миграций из папки migrations
}

func main() {
	// Параметры подключения (DSN)
	connStr := "user=postgres password=123 dbname=mydatabase host=localhost port=5432 sslmode=disable"

	// Открываем соединение
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Проверяем подключение
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	if err := gooseMigrations(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migrations applied successfully")

	productStore := postgres.NewProductStore(db)
	productHandler := api.NewProductHandler(productStore)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/products", productHandler.CreateProduct)
	e.GET("/products/:id", productHandler.GetProduct)

	// Запуск сервера
	if err := e.Start(":8080"); err != http.ErrServerClosed {
		log.Fatal(err)
	}

}
