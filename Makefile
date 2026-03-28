.PHONY: all run build test clean lint wire docker-up docker-down help
.DEFAULT_GOAL := help

# ==============================================================================
# Environment Variables
# ==============================================================================
APP_NAME := gofiber-hax
BIN_DIR := bin
MAIN_FILE := main.go

# ==============================================================================
# Development Commands
# ==============================================================================

## run: รันเซิร์ฟเวอร์แบบ Development
run:
	@echo "Starting development server..."
	@go run $(MAIN_FILE) start

## build: คอมไพล์โปรเจกต์เป็นไฟล์ Binary
build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(BIN_DIR)/$(APP_NAME) $(MAIN_FILE)
	@echo "Build complete: $(BIN_DIR)/$(APP_NAME)"

## test: รัน Unit Test ทั้งหมด
test:
	@echo "Running tests..."
	@go test -v -cover ./...

## lint: ตรวจสอบโค้ดให้ถูกต้องตามมาตรฐาน Go (ต้องติดตั้ง golangci-lint ก่อน)
lint:
	@echo "Running linter..."
	@golangci-lint run ./...

## clean: ลบไฟล์ Binary รูปแบบเก่าทิ้ง
clean:
	@echo "🧹 Cleaning up..."
	@rm -rf $(BIN_DIR)
	@rm -rf keys/*

# ==============================================================================
# Docker Commands
# ==============================================================================

## docker-up: รันฐานข้อมูลทั้งหมดใน Docker (Mongo, MySQL)
docker-up:
	@echo "🐳 Starting Docker containers..."
	@docker-compose up -d

## docker-down: ปิดฐานข้อมูลทั้งหมด
docker-down:
	@echo "🐳 Stopping Docker containers..."
	@docker-compose down

# ==============================================================================
# Helper
# ==============================================================================

## help: โชว์คำสั่งทั้งหมดที่มีใน Makefile นี้
help:
	@echo "Usage: make [command]"
	@echo ""
	@echo "Commands:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
