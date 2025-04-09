include .env
default: run
# Sets variable for common migration Docker command
MIGRATE_CMD = docker run -it --rm --network host --volume $(PWD)/internal/infra/database:/db migrate/migrate
APP_NAME=bankey-server
VERSION=1.0.0
DOCKER_IMAGE=$(APP_NAME):$(VERSION)

docker-build: # Build the Docker image
	@echo "=====> Building Docker image"
	@docker build --no-cache -t $(DOCKER_IMAGE) .
.PHONY: docker-build

air: # Access the Air container
	@docker compose logs -f air
.PHONY: air

psql: # Access the PostgreSQL container
	@echo "=====> Accessing PostgreSQL container"
	@docker exec -it $(DB_NAME) psql -U $(DB_USER) -d $(DB_NAME)
.PHONY: psql

run: # Execute the Go server
	@echo "=====> Running Go server"
	@go run cmd/api/main.go
.PHONY: run

migrate: # Add a new migration
	@echo "=====> Adding a new migration"
	@if [ -z "$(name)" ]; then echo "Migration name is required"; exit 1; fi
	@$(MIGRATE_CMD) create -ext sql -dir /db/migrations $(name)
.PHONY: migrate

migrate-up: # Apply all pending migrations
	@echo "=====> Applying all pending migrations"
	@$(MIGRATE_CMD) -path=/db/migrations -database "$(DB_URL)" up
.PHONY: migrate-up

migrate-down: # Revert all applied migrations
	@echo "=====> Reverting all applied migrations"
	@$(MIGRATE_CMD) -path=/db/migrations -database "$(DB_URL)" down
.PHONY: migrate-down