include .env
export

MIGRATION_PATH=cmd/migrate/migrations

# Run the application
run:
	@echo "Running the API..."
	go run cmd/api/main.go

# Create a new migration file
# Usage: make migration name=create_some_table
migration:
	@echo "Creating migration files for ${name}..."
	migrate create -ext sql -dir $(MIGRATION_PATH) -seq $(name)

# Apply all available migrations (UP)
migrate-up:
	@echo "Running migrations UP..."
	migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" up

# Rollback the last migration (DOWN)
migrate-down:
	@echo "Rolling back last migration..."
	migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" down 1

# Force a specific version (Fix dirty state)
# Usage: make force version=1
force:
	@echo "Forcing migration version to ${version}..."
	migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" force $(version)