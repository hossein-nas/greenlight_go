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

confirm:
	@echo 'Are you sure? [y/N] '; read -r ans; [ "$${ans:-N}" = "y" ]

# Build the Go binary
build:
	$(GO) build -o $(BINARY) ./cmd/api # Adjust the path to your main file if different

# Run the Go project
run: build
	./$(BINARY) -pretty-logger $(ARGS)

run-limiter-off: build
	./$(BINARY) -pretty-logger -limiter-enabled=false

run-load-test-mode: build
	./$(BINARY) -limiter-enabled=false -db-max-open-conns=50 -db-max-idle-conns=50 -db-max-idle-time=20s -port=40000

# Run tests
test:
	$(GO) test ./... -v

# Clean up generated files
clean:
	rm -f $(BINARY)

# Apply migrations (up)
migration: confirm
	$(MIGRATE) create -seq -ext=.sql -dir $(MIGRATIONS_DIR) ${name}

migrate-up: confirm
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database $(DB_URL) up

# Roll back migrations (down)
migrate-down: confirm
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database $(DB_URL) down

migrate-down-1: confirm
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database $(DB_URL) down 1


.PHONY: audit
audit: vendor
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Runnding tests...'
	go test -race -vet=off ./...


.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor



current_time = $(shell date +"%Y-%m-%dT%H:%M:%S%z")
git_description = $(shell git describe --always --dirty)
linker_flags = '-s -X main.buildTime=${current_time} -X main.version=${git_description}'

.PHONY: build/api
build/api:
	@echo 'Bulding cmd/api...'
	go build -ldflags=${linker_flags} -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/linx_amd64/api ./cmd/api