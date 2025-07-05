# Makefile for go-zulip

.PHONY: test test-unit test-integration test-coverage lint fmt build clean help dev-setup dev-up dev-down dev-exec

# Docker environment
DEV_COMPOSE_FILES := -f docker-compose-dev-env.yml -f docker-compose-dev.yml -f docker-compose-zulip.yml
PROJECT_NAME := go-zulip
DEV_SERVICE := go-zulip

# Default target
help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Test targets
test: test-unit ## Run all tests (unit tests only by default)

test-unit: ## Run unit tests
	@if docker compose -p $(PROJECT_NAME) ps | grep -q $(DEV_SERVICE); then \
		docker compose -p $(PROJECT_NAME) exec $(DEV_SERVICE) go test -v -race -cover -short ./...; \
	else \
		echo "Development environment not running. Start with 'make dev-setup' first."; \
		go test -v -race -cover -short ./...; \
	fi

test-integration: ## Run integration tests (requires environment variables and running dev env)
	@if docker compose -p $(PROJECT_NAME) ps | grep -q $(DEV_SERVICE); then \
		docker compose -p $(PROJECT_NAME) exec $(DEV_SERVICE) go test -v -race -cover ./test/integration/...; \
	else \
		echo "Development environment not running. Start with 'make dev-setup' first."; \
		go test -v -race -cover ./test/integration/...; \
	fi

test-coverage: ## Run tests with coverage report
	@if docker compose -p $(PROJECT_NAME) ps | grep -q $(DEV_SERVICE); then \
		docker compose -p $(PROJECT_NAME) exec $(DEV_SERVICE) sh -c "go test -v -race -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html"; \
		echo "Coverage report generated: coverage.html"; \
	else \
		echo "Development environment not running. Start with 'make dev-setup' first."; \
		go test -v -race -coverprofile=coverage.out ./...; \
		go tool cover -html=coverage.out -o coverage.html; \
		echo "Coverage report generated: coverage.html"; \
	fi

# Code quality targets
lint: ## Run linters
	@if docker compose -p $(PROJECT_NAME) ps | grep -q $(DEV_SERVICE); then \
		docker compose -p $(PROJECT_NAME) exec $(DEV_SERVICE) go tool golangci-lint run; \
	else \
		echo "Development environment not running. Start with 'make dev-setup' first."; \
		go tool golangci-lint run; \
	fi

fmt: ## Format code
	@if docker compose -p $(PROJECT_NAME) ps | grep -q $(DEV_SERVICE); then \
		docker compose -p $(PROJECT_NAME) exec $(DEV_SERVICE) go tool gofumpt -w .; \
	else \
		echo "Development environment not running. Start with 'make dev-setup' first."; \
		go tool gofumpt -w .; \
	fi

# Build targets
build: ## Build the library (verify it compiles)
	@if docker compose -p $(PROJECT_NAME) ps | grep -q $(DEV_SERVICE); then \
		docker compose -p $(PROJECT_NAME) exec $(DEV_SERVICE) go build ./...; \
	else \
		echo "Development environment not running. Start with 'make dev-setup' first."; \
		go build ./...; \
	fi

tidy: ## Tidy go modules
	@if docker compose -p $(PROJECT_NAME) ps | grep -q $(DEV_SERVICE); then \
		docker compose -p $(PROJECT_NAME) exec $(DEV_SERVICE) go mod tidy; \
	else \
		echo "Development environment not running. Start with 'make dev-setup' first."; \
		go mod tidy; \
	fi

# Clean targets
clean: ## Clean generated files
	rm -f coverage.out coverage.html

# Development targets
dev-setup: ## Build development environment (Docker compose build)
	@echo "Building development environment..."
	@if [ ! -f docker-compose-dev-env.yml ]; then \
		echo "Copying dev environment example file..."; \
		cp docker-compose-dev-env.example.yml docker-compose-dev-env.yml; \
	fi
	@docker compose -p $(PROJECT_NAME) -f docker-compose-dev.yml -f docker-compose-zulip.yml build --pull

dev-up: ## Start development environment (Docker compose up)
	@echo "Starting development environment..."
	@docker compose -p $(PROJECT_NAME) $(DEV_COMPOSE_FILES) up --detach

dev-down: ## Stop development environment
	@echo "Stopping development environment..."
	@docker compose -p $(PROJECT_NAME) $(DEV_COMPOSE_FILES) down

dev-exec: ## Execute a command in the development container (usage: make dev-exec CMD="command")
	@docker compose -p $(PROJECT_NAME) exec $(DEV_SERVICE) $(CMD)

dev-logs: ## Show development environment logs
	@docker compose -p $(PROJECT_NAME) logs --follow --tail 100

# CI targets (for local development with Docker)
ci: tidy lint test-unit ## Run CI pipeline

# CI targets (for GitHub Actions - runs locally without Docker)
ci-local: ## Run CI pipeline locally without Docker (for GitHub Actions)
	go mod tidy
	go tool golangci-lint run  
	go test -v -race -cover -short ./...
	
verify: ci ## Verify the project is ready for commit
	@echo "Project verification passed!"
