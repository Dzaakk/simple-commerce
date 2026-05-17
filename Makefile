APP_NAME := simple-commerce
BIN_DIR := bin
GO := go
AIR := air
DOCKER_COMPOSE := docker compose

ifeq ($(OS),Windows_NT)
	BINARY := $(BIN_DIR)/$(APP_NAME).exe
	MKDIR_BIN := powershell -NoProfile -Command "New-Item -ItemType Directory -Force -Path '$(BIN_DIR)' | Out-Null"
	CLEAN_BIN := powershell -NoProfile -Command "if (Test-Path '$(BIN_DIR)') { Remove-Item -Recurse -Force '$(BIN_DIR)' }"
else
	BINARY := $(BIN_DIR)/$(APP_NAME)
	MKDIR_BIN := mkdir -p $(BIN_DIR)
	CLEAN_BIN := rm -rf $(BIN_DIR)
endif

.PHONY: help run dev build build-linux build-windows test tidy fmt vet deps docker-build docker-up docker-up-d docker-down docker-logs docker-ps k6-smoke clean

help:
	@echo "Available commands:"
	@echo "  make run              Run the app with go run"
	@echo "  make dev              Run the app with Air hot reload"
	@echo "  make build            Build local binary into ./bin"
	@echo "  make build-linux      Build Linux amd64 binary"
	@echo "  make build-windows    Build Windows amd64 binary"
	@echo "  make test             Run all Go tests"
	@echo "  make tidy             Run go mod tidy"
	@echo "  make fmt              Format Go files"
	@echo "  make vet              Run go vet"
	@echo "  make docker-up        Start services with rebuild"
	@echo "  make docker-down      Stop services"
	@echo "  make docker-logs      Follow service logs"
	@echo "  make docker-ps        Show service status"
	@echo "  make k6-smoke         Run k6 smoke test"
	@echo "  make clean            Remove build artifacts"

run:
	$(GO) run .

dev:
	$(AIR)

build:
	$(MKDIR_BIN)
	$(GO) build -trimpath -o $(BINARY) .

build-linux:
	$(MKDIR_BIN)
	GOOS=linux GOARCH=amd64 $(GO) build -trimpath -o $(BIN_DIR)/$(APP_NAME)-linux-amd64 .

build-windows:
	$(MKDIR_BIN)
	GOOS=windows GOARCH=amd64 $(GO) build -trimpath -o $(BIN_DIR)/$(APP_NAME)-windows-amd64.exe .

test:
	$(GO) test ./...

tidy:
	$(GO) mod tidy

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

docker-up:
	$(DOCKER_COMPOSE) up --build -d

docker-down:
	$(DOCKER_COMPOSE) down

docker-logs:
	$(DOCKER_COMPOSE) logs -f

docker-ps:
	$(DOCKER_COMPOSE) ps

k6-smoke:
	k6 run tests/k6/smoke.js

clean:
	$(CLEAN_BIN)
