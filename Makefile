APP_NAME := simple-commerce
BIN_DIR := bin
GO := go
AIR := air
DOCKER_COMPOSE := docker compose

ifneq ($(filter v2,$(MAKECMDGOALS)),)
	CATALOG_API_VERSION := v2
else
	CATALOG_API_VERSION ?= v1
endif

K6_CATALOG_ENV := -e CATALOG_API_VERSION=$(CATALOG_API_VERSION)

ifeq ($(OS),Windows_NT)
	BINARY := $(BIN_DIR)/$(APP_NAME).exe
	MKDIR_BIN := powershell -NoProfile -Command "New-Item -ItemType Directory -Force -Path '$(BIN_DIR)' | Out-Null"
	CLEAN_BIN := powershell -NoProfile -Command "if (Test-Path '$(BIN_DIR)') { Remove-Item -Recurse -Force '$(BIN_DIR)' }"
else
	BINARY := $(BIN_DIR)/$(APP_NAME)
	MKDIR_BIN := mkdir -p $(BIN_DIR)
	CLEAN_BIN := rm -rf $(BIN_DIR)
endif

.PHONY: help run dev build build-linux build-windows test tidy fmt vet deps docker-build docker-up-d docker-up docker-down docker-logs docker-ps k6-smoke k6-load k6-load-100 k6-load-300 k6-load-500 k6-load-1000 k6-stress k6-spike v1 v2 clean

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
	@echo "  make k6-load          Run k6 load test with default profile"
	@echo "  make k6-load-100      Run k6 load test with 100 VUs"
	@echo "  make k6-load-300      Run k6 load test with 300 VUs"
	@echo "  make k6-load-500      Run k6 load test with 500 VUs"
	@echo "  make k6-load-1000     Run k6 load test with 1000 VUs"
	@echo "  make k6-stress        Run k6 stress test"
	@echo "  make k6-spike         Run k6 spike test"
	@echo "  Add v2 after a k6 target to run against catalog API v2"
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

k6-load:
	k6 run $(K6_CATALOG_ENV) tests/k6/load.js

k6-load-100:
	k6 run $(K6_CATALOG_ENV) -e LOAD_PROFILE=100 tests/k6/load.js

k6-load-300:
	k6 run $(K6_CATALOG_ENV) -e LOAD_PROFILE=300 tests/k6/load.js

k6-load-500:
	k6 run $(K6_CATALOG_ENV) -e LOAD_PROFILE=500 tests/k6/load.js

k6-load-1000:
	k6 run $(K6_CATALOG_ENV) -e LOAD_PROFILE=1000 tests/k6/load.js

k6-stress:
	k6 run $(K6_CATALOG_ENV) tests/k6/stress.js

k6-spike:
	k6 run $(K6_CATALOG_ENV) tests/k6/spike.js

v1 v2:
	@:

clean:
	$(CLEAN_BIN)
