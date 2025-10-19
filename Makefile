# Servers Filters Project Makefile

.PHONY: help setup test build run clean

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

setup: ## Setup the project
	@echo "Setting up the project..."
	cd backend && go mod download
	cd frontend && npm install

test: ## Run tests
	@echo "Running tests..."
	cd backend && go test ./...

build: ## Build the application
	@echo "Building application..."
	cd backend && go build -o servers-listing main.go
	cd frontend && npm run build

run: ## Run the application with Docker
	@echo "Starting application..."
	docker-compose up

run-dev: ## Run in development mode
	@echo "Starting in development mode..."
	docker-compose -f docker-compose.yml --env-file .env.development up

run-prod: ## Run in production mode
	@echo "Starting in production mode..."
	docker-compose -f docker-compose.yml --env-file .env.production up

clean: ## Clean up
	@echo "Cleaning up..."
	docker-compose down
	docker system prune -f

demo: ## Run demo script
	@echo "Running demo..."
	./demo.sh
