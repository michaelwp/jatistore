APP_NAME=jatistore
BIN_DIR=bin
SWAGGER_DIR=docs

.PHONY: all build run swag migrate-up migrate-down tidy clean

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

clean:
	rm -rf $(BIN_DIR) $(SWAGGER_DIR) 