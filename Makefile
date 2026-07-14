APP_NAME := simple-commerce
BIN_DIR := bin
GO := go
AIR := air
DOCKER_COMPOSE := docker compose
K6 := k6
K6_BASE_URL ?= http://localhost:8080
K6_PRODUCT_LIMIT ?= 100
K6_DURATION ?= 3m

ifeq ($(OS),Windows_NT)
	BINARY := $(BIN_DIR)/$(APP_NAME).exe
	MKDIR_BIN := powershell -NoProfile -Command "New-Item -ItemType Directory -Force -Path '$(BIN_DIR)' | Out-Null"
	CLEAN_BIN := powershell -NoProfile -Command "if (Test-Path '$(BIN_DIR)') { Remove-Item -Recurse -Force '$(BIN_DIR)' }"
else
	BINARY := $(BIN_DIR)/$(APP_NAME)
	MKDIR_BIN := mkdir -p $(BIN_DIR)
	CLEAN_BIN := rm -rf $(BIN_DIR)
endif

.PHONY: help run dev build build-linux build-windows generate test tidy fmt vet deps docker-build docker-up-d docker-up docker-down docker-logs docker-ps k6-smoke k6-load-v1-100 k6-load-v1-300 k6-load-v1-500 k6-load-v1-1000 k6-load-v2-100 k6-load-v2-300 k6-load-v2-500 k6-load-v2-1000 clean

help:
	@echo "Available commands:"
	@echo "  make run              Run the app with go run"
	@echo "  make dev              Run the app with Air hot reload"
	@echo "  make build            Build local binary into ./bin"
	@echo "  make build-linux      Build Linux amd64 binary"
	@echo "  make build-windows    Build Windows amd64 binary"
	@echo "  make generate         Generate Go API contracts from api.yaml"
	@echo "  make test             Run all Go tests"
	@echo "  make tidy             Run go mod tidy"
	@echo "  make fmt              Format Go files"
	@echo "  make vet              Run go vet"
	@echo "  make docker-up        Start services with rebuild"
	@echo "  make docker-down      Stop services"
	@echo "  make docker-logs      Follow service logs"
	@echo "  make docker-ps        Show service status"
	@echo "  make k6-smoke         Run k6 smoke test"
	@echo "  make k6-load-v1-100   Run v1 catalog browsing load test with 100 VUs"
	@echo "  make k6-load-v1-300   Run v1 catalog browsing load test with 300 VUs"
	@echo "  make k6-load-v1-500   Run v1 catalog browsing load test with 500 VUs"
	@echo "  make k6-load-v1-1000  Run v1 catalog browsing load test with 1000 VUs"
	@echo "  make k6-load-v2-100   Run v2 catalog browsing load test with 100 VUs"
	@echo "  make k6-load-v2-300   Run v2 catalog browsing load test with 300 VUs"
	@echo "  make k6-load-v2-500   Run v2 catalog browsing load test with 500 VUs"
	@echo "  make k6-load-v2-1000  Run v2 catalog browsing load test with 1000 VUs"
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

generate:
	$(GO) run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1 --config oapi-codegen.catalog-v2.yaml -o internal/api/generated/openapi.gen.go api.yaml

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

k6-load-v1-100:
	$(K6) run --env BASE_URL=$(K6_BASE_URL) --env PRODUCT_LIST_ENDPOINT=/api/v1/product --env PRODUCT_DETAIL_ENDPOINT=/api/v1/product --env PRODUCT_LIMIT=$(K6_PRODUCT_LIMIT) --vus 100 --duration $(K6_DURATION) tests/k6/catalog-browsing.js

k6-load-v1-300:
	$(K6) run --env BASE_URL=$(K6_BASE_URL) --env PRODUCT_LIST_ENDPOINT=/api/v1/product --env PRODUCT_DETAIL_ENDPOINT=/api/v1/product --env PRODUCT_LIMIT=$(K6_PRODUCT_LIMIT) --vus 300 --duration $(K6_DURATION) tests/k6/catalog-browsing.js

k6-load-v1-500:
	$(K6) run --env BASE_URL=$(K6_BASE_URL) --env PRODUCT_LIST_ENDPOINT=/api/v1/product --env PRODUCT_DETAIL_ENDPOINT=/api/v1/product --env PRODUCT_LIMIT=$(K6_PRODUCT_LIMIT) --vus 500 --duration $(K6_DURATION) tests/k6/catalog-browsing.js

k6-load-v1-1000:
	$(K6) run --env BASE_URL=$(K6_BASE_URL) --env PRODUCT_LIST_ENDPOINT=/api/v1/product --env PRODUCT_DETAIL_ENDPOINT=/api/v1/product --env PRODUCT_LIMIT=$(K6_PRODUCT_LIMIT) --vus 1000 --duration $(K6_DURATION) tests/k6/catalog-browsing.js

k6-load-v2-100:
	$(K6) run --env BASE_URL=$(K6_BASE_URL) --env PRODUCT_LIST_ENDPOINT=/api/v2/product --env PRODUCT_DETAIL_ENDPOINT=/api/v2/product --env PRODUCT_LIMIT=$(K6_PRODUCT_LIMIT) --vus 100 --duration $(K6_DURATION) tests/k6/catalog-browsing.js

k6-load-v2-300:
	$(K6) run --env BASE_URL=$(K6_BASE_URL) --env PRODUCT_LIST_ENDPOINT=/api/v2/product --env PRODUCT_DETAIL_ENDPOINT=/api/v2/product --env PRODUCT_LIMIT=$(K6_PRODUCT_LIMIT) --vus 300 --duration $(K6_DURATION) tests/k6/catalog-browsing.js

k6-load-v2-500:
	$(K6) run --env BASE_URL=$(K6_BASE_URL) --env PRODUCT_LIST_ENDPOINT=/api/v2/product --env PRODUCT_DETAIL_ENDPOINT=/api/v2/product --env PRODUCT_LIMIT=$(K6_PRODUCT_LIMIT) --vus 500 --duration $(K6_DURATION) tests/k6/catalog-browsing.js

k6-load-v2-1000:
	$(K6) run --env BASE_URL=$(K6_BASE_URL) --env PRODUCT_LIST_ENDPOINT=/api/v2/product --env PRODUCT_DETAIL_ENDPOINT=/api/v2/product --env PRODUCT_LIMIT=$(K6_PRODUCT_LIMIT) --vus 1000 --duration $(K6_DURATION) tests/k6/catalog-browsing.js

clean:
	$(CLEAN_BIN)
