APP_NAME := simple-commerce
BIN_DIR := bin
GO := go

ifeq ($(OS),Windows_NT)
	BINARY := $(BIN_DIR)/$(APP_NAME).exe
	MKDIR_BIN := powershell -NoProfile -Command "New-Item -ItemType Directory -Force -Path '$(BIN_DIR)' | Out-Null"
	CLEAN_BIN := powershell -NoProfile -Command "if (Test-Path '$(BIN_DIR)') { Remove-Item -Recurse -Force -LiteralPath '$(BIN_DIR)' }"
else
	BINARY := $(BIN_DIR)/$(APP_NAME)
	MKDIR_BIN := mkdir -p $(BIN_DIR)
	CLEAN_BIN := rm -rf $(BIN_DIR)
endif

.PHONY: run build test tidy fmt vet compose-config docker-up docker-down clean

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

clean:
	$(CLEAN_BIN)
