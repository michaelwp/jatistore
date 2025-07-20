APP_NAME=jatistore
BIN_DIR=bin
SWAGGER_DIR=docs

.PHONY: all build run swag migrate-up migrate-down tidy clean lint pre-commit install-hooks

all: build

build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) main.go

run:
	go run main.go

swag:
	swag init --parseDependency --parseInternal --output $(SWAGGER_DIR)

migrate-up:
	migrate -path internal/database/migrations -database "$$DATABASE_URL" up

migrate-down:
	migrate -path internal/database/migrations -database "$$DATABASE_URL" down

tidy:
	go mod tidy

lint:
	golangci-lint run

pre-commit:
	@echo "Running pre-commit checks..."
	@make lint
	@echo "Pre-commit checks passed!"

install-hooks:
	@echo "Installing git hooks..."
	@chmod +x .git/hooks/pre-commit
	@echo "Git hooks installed successfully!"

clean:
	rm -rf $(BIN_DIR) $(SWAGGER_DIR) 