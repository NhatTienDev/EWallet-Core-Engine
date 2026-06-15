# Run server with live-reload using Air tool
dev:
	air

# Automatically generate code from SQL queries (Run this command when edit .sql file)
sqlc:
	sqlc generate

# Clean up unnecessary libraries
tidy:
	go mod tidy

# Reformat the go.mod file
vendor:
	go mod vendor

# Restart all background processes, read Dockerfile and rebuild the newest image before running (Restore a clean terminal)
dev-docker-up:
	docker compose -f docker-compose.dev.yaml up --build -d

# Stop the system but keep the data in the DB in order to continue coding
dev-docker-down:
	docker compose -f docker-compose.dev.yaml down

# View real-time logs
dev-docker-logs:
	docker compose -f docker-compose.dev.yaml logs -f

# Stop the system and clear all data in the DB to start over
dev-docker-clear:
	docker compose -f docker-compose.dev.yaml down -v

prod-docker-up:
	docker compose -f docker-compose.prod.yaml up --build -d

prod-docker-down:
	docker compose -f docker-compose.prod.yaml down

prod-docker-logs:
	docker compose -f docker-compose.prod.yaml logs -f

# Just build a test Dockerfile to check for compilation errors
prod-docker-build:
	docker compose -f docker-compose.prod.yaml build
