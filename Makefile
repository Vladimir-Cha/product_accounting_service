#Переменные
APP_NAME = product_accounting_service
DOCKER_IMAGE = product_accounting_service
GO_FILES = $(shell find . -type f -name '*.go')
PORT = 8080
DB_PORT = 5432

.PHONY: all build run test clean docker-build docker-run compose-up compose-down

all: build

#Сборка проекта
build:
	@echo "Building..."
	@./bin/$(APP_NAME)

#Запуск
run: build
	@echo "Starting..."
	@./bin/$(APP_NAME)

# Docker
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .

docker-run: docker-build
	@echo "Running Docker container..."
	@docker run -d -p $(PORT):$(PORT) --name $(APP_NAME) $(DOCKER_IMAGE)

docker-stop:
	@docker stop $(APP_NAME)

# Docker Compose
compose-up:
	@echo "Starting services with Docker Compose..."
	@docker compose up -d --build

compose-down:
	@echo "Stopping services..."
	@docker compose down

tidy:
	@go mod tidy -v

#Удаление созданного бинарника
clean:
	@echo "Cleaning..."
	@rm -rf bin/