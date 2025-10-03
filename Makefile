.PHONY: help build run dev test clean migrate-up migrate-down migrate-status migrate-reset db-up db-down

APP_NAME := customable-corporate-site-api
MAIN_FILE := cmd/server/main.go
MIGRATE_FILE := cmd/migrate/main.go

help:
	@echo "Makefile commands:"
	@echo "  build           - Build the application"
	@echo "  run             - Run the application"
	@echo "  dev             - Run the application in development mode with live reload"
	@echo "  test            - Run tests"
	@echo "  clean           - Clean up build artifacts"
	@echo "  migrate-up      - Apply all pending migrations"
	@echo "  migrate-down    - Rollback the last migration"
	@echo "  migrate-status  - Show current migration status"
	@echo "  migrate-reset   - Reset all migrations (down then up)"
	@echo "  db-up           - Start the PostgreSQL database using Docker"
	@echo "  db-down         - Stop and remove the PostgreSQL database Docker container"

build:
	@echo "Building the application..."
	@mkdir -p bin
	go build -o bin/$(APP_NAME) $(MAIN_FILE)
	@echo "Build completed: bin/$(APP_NAME)"

run: build
	@echo "Starting application..."
	@bin/$(APP_NAME)

dev:
	@echo "Starting application in development mode with live reload..."
	@if command -v air >/dev/null 2>&1; then \
		air -c .air.toml; \
	else \
		echo "Air is not installed. Run: go install github.com/cosmtrek/air@latest"; \
		go run $(MAIN_FILE); \
	fi

test:
	@echo "Running tests..."
	go test -v ./...

clean:
	@echo "Cleaning up build artifacts..."
	@rm -rf bin/
	@rm -rf build/
	@echo "Clean completed."

db-up:
	@echo "Starting PostgreSQL database using Docker..."
	docker-compose up -d postgres
	@echo "PostgreSQL database started."

db-down:
	@echo "Stopping PostgreSQL database..."
	docker-compose down
	@echo "PostgreSQL database stopped."

db-status:
	@echo "Checking PostgreSQL database status..."
	docker-compose ps
	@echo "PostgreSQL database status checked."

migrate-build:
	@echo "Building the migration tool..."
	@mkdir -p bin
	go build -o bin/migrate $(MIGRATE_FILE)
	@echo "Migration tool built: bin/migrate"

migrate-up: migrate-build
	@echo "Applying all pending migrations..."
	@bin/migrate -up

migrate-down: migrate-build
	@echo "Rolling back the last migration..."
	@bin/migrate -down

migrate-status: migrate-build
	@echo "Checking migration status..."
	@bin/migrate -status

migrate-reset: migrate-build
	@echo "Resetting all migrations (down then up)..."
	@bin/migrate -reset

deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies updated."

	