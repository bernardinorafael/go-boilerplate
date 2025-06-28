include .env
default: run
# Sets variable for common migration Docker command
MIGRATE_CMD = docker run -it --rm --network host --volume $(PWD)/internal/infra/database:/db migrate/migrate
APP_NAME=go-boilerplate
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

.PHONY: air
air: # Access the Air container
	@docker compose logs -f air

.PHONY: redis
redis: # Access the Redis container
	@echo "=====> Accessing Redis container"
	@docker exec -it ${APP_NAME}-redis redis-cli

.PHONY: psql
psql: # Access the PostgreSQL container
	@echo "=====> Accessing PostgreSQL container"
	@docker exec -it ${APP_NAME}-postgres psql -U $(DB_USER) -d $(DB_NAME)

.PHONY: tests
tests: # Run the Go tests with coverage
	@echo "=====> Running tests with coverage"
	go test -v -cover ./...

.PHONY: tidy
tidy: # Run go mod tidy
	@echo "=====> Running go mod tidy"
	go mod tidy

.PHONY: run
run: # Execute the Go server
	@go run cmd/api/main.go

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
