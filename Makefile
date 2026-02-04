include .env
export

MIGRATION_PATH=cmd/migrate/migrations

# Run the application
run:
	@echo "Running the API locally..."
	go run cmd/api/main.go

docker-run:
	@echo "Starting Docker environment..."
	docker compose up --build -d

# may want to run our application locally while developing (must do make run)
docker-run-db:
	@echo "Starting Docker database..."
	docker compose up db -d

# Stop all containers
docker-down:
	@echo "Stopping Docker environment..."
	docker compose down

# Stop containers AND delete the database volume (Fresh Start)
docker-clean:
	@echo "Stopping Docker and removing volumes..."
	docker compose down -v

docker-logs:
	@echo "Following logs..."
	docker compose logs -f api

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