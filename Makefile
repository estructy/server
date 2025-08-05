export $(shell sed 's/=.*//' .env)

MIGRATION_DIR = ./migrations

# Create a new migration
.PHONY: create-migration
create-migration:
	@if [ -z "$(name)" ]; then \
		echo "Error: Please specify the migration name using 'make create-migration name=your_migration'"; \
		exit 1; \
	fi
	migrate create -dir $(MIGRATION_DIR) -ext sql $(name)

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
