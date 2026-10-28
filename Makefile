.PHONY: dev run-backend dev-backend run-frontend dev-frontend build clean k8s-deploy k8s-restart k8s-status k8s-logs compose-up compose-up-gateway compose-down compose-logs compose-ps

PORT ?= 8089

dev:
	@echo "=================================================="
	@echo "🚀 Starting Enterprise Architecture Sizing Suite..."
	@echo "📡 Backend (Go):     http://localhost:$(PORT)"
	@echo "💻 Frontend (Nuxt):  http://localhost:5454"
	@echo "=================================================="
	@(trap 'kill 0' SIGINT; cd backend && PORT=$(PORT) go run ./cmd/api & cd frontend && NUXT_PUBLIC_API_BASE=http://localhost:$(PORT) pnpm dev & wait)

run-backend: dev-backend

dev-backend:
	@echo "🚀 Running Go Sizing Backend on :$(PORT)..."
	@cd backend && PORT=$(PORT) go run ./cmd/api

run-frontend: dev-frontend

dev-frontend:
	@echo "💻 Running Nuxt Sizing Frontend on :5454..."
	@cd frontend && NUXT_PUBLIC_API_BASE=http://localhost:$(PORT) pnpm dev

build:
	@echo "📦 Building Golang backend..."
	@cd backend && go build -o bin/server ./cmd/api
	@echo "📦 Building Frontend..."
	@cd frontend && pnpm build
	@echo "✨ Build completed successfully!"

clean:
	@rm -rf backend/bin frontend/.output frontend/.nuxt

# ---- Kubernetes & Deployment (k3s / vps-malang) ----
REGISTRY_LOCAL_PORT ?= 5005
K8S_NS ?= sizing

k8s-deploy: ## Build images (linux/amd64), push ke local registry k8s, dan restart deployment
	@echo "🔌 Membuka port-forward ke registry cluster di port $(REGISTRY_LOCAL_PORT)..."
	@kubectl port-forward -n infra svc/docker-registry $(REGISTRY_LOCAL_PORT):5000 >/dev/null 2>&1 & PF_PID=$$!; \
	trap 'kill $$PF_PID 2>/dev/null || true' EXIT; \
	sleep 2; \
	echo "📦 Building & pushing backend..."; \
	docker build --platform linux/amd64 -t localhost:$(REGISTRY_LOCAL_PORT)/sizing-backend:latest -f backend/Dockerfile backend/ && \
	docker push localhost:$(REGISTRY_LOCAL_PORT)/sizing-backend:latest && \
	echo "📦 Building & pushing frontend..."; \
	docker build --platform linux/amd64 -t localhost:$(REGISTRY_LOCAL_PORT)/sizing-frontend:latest -f frontend/Dockerfile frontend/ && \
	docker push localhost:$(REGISTRY_LOCAL_PORT)/sizing-frontend:latest && \
	echo "🚀 Applying k3s overlay manifests..."; \
	kubectl apply -k deploy/kustomize/overlays/k3s && \
	echo "🔄 Rolling restart deployments..."; \
	kubectl rollout restart deployment/sizing-backend deployment/sizing-frontend -n $(K8S_NS) && \
	kubectl rollout status deployment/sizing-backend -n $(K8S_NS) --timeout=90s && \
	kubectl rollout status deployment/sizing-frontend -n $(K8S_NS) --timeout=90s && \
	echo "✅ Deployment selesai dan berhasil running!"

k8s-restart: ## Restart deployment sizing di Kubernetes
	@kubectl rollout restart deployment/sizing-backend deployment/sizing-frontend -n $(K8S_NS)
	@kubectl rollout status deployment/sizing-backend -n $(K8S_NS)
	@kubectl rollout status deployment/sizing-frontend -n $(K8S_NS)

k8s-status: ## Cek status pod, service, dan ingress sizing
	@kubectl get pods,svc,ingress -n $(K8S_NS) -o wide

k8s-logs: ## Cek log backend & frontend
	@echo "=== BACKEND LOGS ==="
	@kubectl logs deployment/sizing-backend -n $(K8S_NS) --tail=25
	@echo ""
	@echo "=== FRONTEND LOGS ==="
	@kubectl logs deployment/sizing-frontend -n $(K8S_NS) --tail=25

# ---- Docker Compose ----
compose-up: ## Jalankan seluruh stack (db, backend, frontend) via docker compose
	docker compose up -d --build

compose-up-gateway: ## Jalankan stack lengkap dengan Nginx reverse proxy di port 80
	docker compose --profile gateway up -d --build

compose-down: ## Matikan seluruh container docker compose
	docker compose down

compose-logs: ## Pantau log container docker compose
	docker compose logs -f

compose-ps: ## Cek status container docker compose
	docker compose ps
