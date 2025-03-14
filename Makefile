.PHONY: run build clean test swagger seed test-api

# Default target
all: build

# Run the application
run:
	go run cmd/api/main.go

# Build the application
build:
	go build -o bin/kontena cmd/api/main.go

# Clean build artifacts
clean:
	rm -rf bin/

# Run tests
test:
	go test -v ./...

# Generate Swagger documentation
swagger:
	swag init -g cmd/api/main.go

# Create database
createdb:
	psql -U postgres -c "CREATE DATABASE kontena;"

# Drop database
dropdb:
	psql -U postgres -c "DROP DATABASE IF EXISTS kontena;"

# Reset database (drop and create)
resetdb: dropdb createdb

# Install dependencies
deps:
	go mod download

# Run linter
lint:
	go vet ./...

# Seed the database with test data
seed:
	go run scripts/seed.go

# Test the API endpoints
test-api:
	./scripts/test_api.sh

# Help
help:
	@echo "Available targets:"
	@echo "  run       - Run the application"
	@echo "  build     - Build the application"
	@echo "  clean     - Clean build artifacts"
	@echo "  test      - Run tests"
	@echo "  swagger   - Generate Swagger documentation"
	@echo "  createdb  - Create database"
	@echo "  dropdb    - Drop database"
	@echo "  resetdb   - Reset database (drop and create)"
	@echo "  deps      - Install dependencies"
	@echo "  lint      - Run linter"
	@echo "  seed      - Seed the database with test data"
	@echo "  test-api  - Test the API endpoints"
	@echo "  help      - Show this help message" 