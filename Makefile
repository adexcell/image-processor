.PHONY: build run-api run-worker up down logs clean help


help: ## Показать справку
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the project
	go build -o bin/api ./cmd/api/main.go
	go build -o bin/worker ./cmd/worker/main.go

run-api: ## Run the API service locally
	go run ./cmd/api/main.go

run-worker: ## Run the Worker service locally
	go run ./cmd/worker/main.go

up: ## Start all services using docker-compose
	docker-compose up --build -d

down: ## Stop all services
	docker-compose down

logs: ## View logs from docker-compose
	docker-compose logs -f

clean: ## Remove binaries and temporary files
	rm -rf bin/ uploads/* processed/* logs/*
