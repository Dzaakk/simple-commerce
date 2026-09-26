APP_NAME := simple-commerce
BIN_DIR := bin
GO := go
K6 := k6
K6_BASE_URL ?= http://localhost:8080
K6_DURATION ?= 3m

ifeq ($(OS),Windows_NT)
	BINARY := $(BIN_DIR)/$(APP_NAME).exe
	MKDIR_BIN := powershell -NoProfile -Command "New-Item -ItemType Directory -Force -Path '$(BIN_DIR)' | Out-Null"
	CLEAN_BIN := powershell -NoProfile -Command "if (Test-Path '$(BIN_DIR)') { Remove-Item -Recurse -Force -LiteralPath '$(BIN_DIR)' }"
else
	BINARY := $(BIN_DIR)/$(APP_NAME)
	MKDIR_BIN := mkdir -p $(BIN_DIR)
	CLEAN_BIN := rm -rf $(BIN_DIR)
endif

.PHONY: run build test tidy fmt vet compose-config docker-up docker-down k6-smoke k6-v1 k6-v2 clean

run:
	$(GO) run .

build:
	$(MKDIR_BIN)
	$(GO) build -trimpath -o $(BINARY) .

test:
	$(GO) test ./...

tidy:
	$(GO) mod tidy

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

compose-config:
	docker compose config --quiet

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down

k6-smoke:
	$(K6) run -e BASE_URL=$(K6_BASE_URL) tests/k6/smoke.js

k6-v1:
	$(K6) run -e BASE_URL=$(K6_BASE_URL) -e ENDPOINT=/api/v1/product -e DURATION=$(K6_DURATION) tests/k6/catalog-browsing.js

k6-v2:
	$(K6) run -e BASE_URL=$(K6_BASE_URL) -e ENDPOINT=/api/v2/product -e DURATION=$(K6_DURATION) tests/k6/catalog-browsing.js

clean:
	$(CLEAN_BIN)
