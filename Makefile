include .env
export $(shell sed 's/=.*//' .env)

MIGRATION_DIR = ./migrations
DB_URL = $(DATABASE_URL)

# Create a new migration
.PHONY: create-migration
create-migration:
	@if [ -z "$(name)" ]; then \
		echo "Error: Please specify the migration name using 'make create-migration name=your_migration'"; \
		exit 1; \
	fi
	migrate create -dir $(MIGRATION_DIR) -ext sql $(name)

# Apply migrations
.PHONY: migrate-up
migrate-up:
	migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" up

# Rollback last migration
.PHONY: migrate-down
migrate-down:
	migrate -database "$(DB_URL)" -path $(MIGRATION_DIR) down 1

# Reset database (dangerous)
.PHONY: migrate-reset
migrate-reset:
	migrate -database "$(DB_URL)" -path $(MIGRATION_DIR) down
	migrate -database "$(DB_URL)" -path $(MIGRATION_DIR) up

# Show migration status
.PHONY: migrate-status
migrate-status:
	migrate -database "$(DB_URL)" -path $(MIGRATION_DIR) version

# Compile sqlc queries
.PHONY: sqlc
sqlc:
	sqlc generate

# Build the application
.PHONY: build
build:
	go build -o ./bin/controlriver ./cmd/api_server/main.go

# Run the application
.PHONY: run-build
run-build:
	./bin/controlriver

.PHONY: tsp-compile
tsp-compile:
	cd ./docs/ && tsp compile .
