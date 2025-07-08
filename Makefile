include .env
default: run
APP_NAME=go-boilerplate

# ==============================================================================
# Docker Commands
# ==============================================================================
VERSION=1.0.0
DOCKER_IMAGE=$(APP_NAME):$(VERSION)

.PHONY: compose-up
compose-up: # Start the Docker containers
	@echo "=====> Starting Docker containers"
	@docker compose up -d --remove-orphans

.PHONY: compose-stop
compose-stop: # Stop the Docker containers
	@echo "=====> Stopping Docker containers"
	@docker compose stop

.PHONY: compose-down
compose-down: # down the Docker containers
	@echo "=====> Removing Docker containers"
	@docker compose down

.PHONY: docker-build
docker-build: # Build the Docker image
	@echo "=====> Building Docker image"
	@docker build --no-cache -t $(DOCKER_IMAGE) .

# ==============================================================================
# Mocks
# ==============================================================================
GO_MODULE_PATH := $(shell go list -m)

.PHONY: mock
mock:
	@echo "=====> Generating services mock"
	@rm -rf __mocks
	@mkdir -p __mocks
	@for dir in internal/domain/*/; do \
		domain=$$(basename $$dir); \
		echo "Generating mock for $$domain service"; \
		mockgen -destination=__mocks/domain/$${domain}/service.go -package=$${domain}mock $(GO_MODULE_PATH)/internal/domain/$$domain Service; \
	done

.PHONY: install-mockgen
install-mockgen:
	@echo "=====> Installing mockgen"
	@go install go.uber.org/mock/mockgen@latest

# ==============================================================================
# Infra
# ==============================================================================
.PHONY: redis
redis: # Access the Redis container
	@echo "=====> Accessing Redis container"
	@docker exec -it ${APP_NAME}-redis redis-cli

.PHONY: psql
psql: # Access the PostgreSQL container
	@echo "=====> Accessing PostgreSQL container"
	@docker exec -it ${APP_NAME}-postgres psql -U $(DB_USER) -d $(DB_NAME)

# ==============================================================================
# Go
# ==============================================================================
.PHONY: test
test: # Run the Go tests with coverage
	@echo "=====> Running tests with coverage"
	go test -v -cover ./...

.PHONY: tidy
tidy: # Run go mod tidy
	@echo "=====> Running go mod tidy"
	go mod tidy

.PHONY: run
run: # Execute the Go server
	@go run cmd/api/main.go

# ==============================================================================
# Migrations
# ==============================================================================
MIGRATE_CMD = docker run -it --rm --network host --volume $(PWD)/internal/infra/database:/db migrate/migrate

.PHONY: migrate
migrate: # Add a new migration
	@echo "=====> Adding a new migration"
	@if [ -z "$(name)" ]; then echo "Migration name is required"; exit 1; fi
	@$(MIGRATE_CMD) create -ext sql -dir /db/migrations $(name)

.PHONY: migrate-up
migrate-up: # Apply all pending migrations
	@echo "=====> Applying all pending migrations"
	@$(MIGRATE_CMD) -path=/db/migrations -database "$(DB_URL)" up

.PHONY: migrate-down
migrate-down: # Revert all applied migrations
	@echo "=====> Reverting all applied migrations"
	@$(MIGRATE_CMD) -path=/db/migrations -database "$(DB_URL)" down

migrate-next: # Apply the last pending migration
	@$(MIGRATE_CMD) -path=/db/migrations -database "$(DB_URL)" up 1
.PHONY: migrate-next

migrate-prev: # Revert the last applied migration
	@$(MIGRATE_CMD) -path=/db/migrations -database "$(DB_URL)" down 1
.PHONY: migrate-prev
