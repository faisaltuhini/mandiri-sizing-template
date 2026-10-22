.PHONY: dev run-backend dev-backend run-frontend dev-frontend build clean

PORT ?= 8089

dev:
	@echo "=================================================="
	@echo "🚀 Starting Enterprise Architecture Sizing Suite..."
	@echo "📡 Backend (Go):     http://localhost:$(PORT)"
	@echo "💻 Frontend (Nuxt):  http://localhost:5454"
	@echo "=================================================="
	@(trap 'kill 0' SIGINT; cd backend && PORT=$(PORT) go run ./cmd/api & cd frontend && pnpm dev & wait)

run-backend: dev-backend

dev-backend:
	@echo "🚀 Running Go Sizing Backend on :$(PORT)..."
	@cd backend && PORT=$(PORT) go run ./cmd/api

run-frontend: dev-frontend

dev-frontend:
	@echo "💻 Running Nuxt Sizing Frontend on :5454..."
	@cd frontend && pnpm dev

build:
	@echo "📦 Building Golang backend..."
	@cd backend && go build -o bin/server ./cmd/api
	@echo "📦 Building Frontend..."
	@cd frontend && pnpm build
	@echo "✨ Build completed successfully!"

clean:
	@rm -rf backend/bin frontend/.output frontend/.nuxt
