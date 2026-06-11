# Run server with live-reload using Air tool
dev:
	air

# Automatically generate code from SQL queries (Run this command when edit .sql file)
sqlc:
	sqlc generate

# Clean up unnecessary libraries and reformat the go.mod file
vendor:
	go mod tidy
	go mod vendor

# Restart all background processes, read Dockerfile and rebuild the newest image before running (Restore a clean terminal)
docker-up:
	docker compose up --build -d
# Stop the system but keep the data in the DB in order to continue coding
docker-down:
	docker compose down
# Stop the system and clear all data in the DB to start over
docker-clear:
	docker compose down -v
# View real-time logs
docker-logs:
	docker compose logs -f
# Just build a test Dockerfile to check for compilation errors
docker-build:
	docker compose build
