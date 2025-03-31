# Variables
GO = go
BINARY = greenlight
MIGRATE = migrate
DB_URL = $(GREENLIGHT_DB_DSN)
MIGRATIONS_DIR = ./migrations

# Phony targets (not actual files)
.PHONY: all build run test clean migrate-up migrate-down

# Default target: build the project
all: build

# Build the Go binary
build:
	$(GO) build -o $(BINARY) ./cmd/api # Adjust the path to your main file if different

# Run the Go project
run: build
	./$(BINARY) -pretty-logger $(ARGS)

run-limiter-off: build
	./$(BINARY) -pretty-logger -limiter-enabled=false

# Run tests
test:
	$(GO) test ./... -v

# Clean up generated files
clean:
	rm -f $(BINARY)

# Apply migrations (up)
migrate-up:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database $(DB_URL) up

# Roll back migrations (down)
migrate-down:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database $(DB_URL) down