.PHONY: run build test clean

APP_NAME = sizing-backend
PORT ?= 8089

run:
	@echo "Starting $(APP_NAME) on port $(PORT)..."
	PORT=$(PORT) go run ./cmd/api

build:
	@echo "Building binary $(APP_NAME)..."
	go build -o bin/$(APP_NAME) ./cmd/api

test:
	go test -v ./...

clean:
	rm -rf bin/
