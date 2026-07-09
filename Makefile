# Generate SQL code from queries
sqlc:
	- sqlc generate -f ./configs/sqlc.yaml

# Create a new database migration
# Usage:
# make migration-create NAME=create_users_table
migration-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME parameter is required."; \
		echo "Usage: make migration-create NAME=create_users_table"; \
		exit 1; \
	fi
	@TIMESTAMP=$$(date +"%Y%m%d%H%M%S"); \
	UP_FILE="database/schema/$${TIMESTAMP}_$(NAME).up.sql"; \
	DOWN_FILE="database/schema/$${TIMESTAMP}_$(NAME).down.sql"; \
	mkdir -p database/schema; \
	echo "-- +migrate Up" > $$UP_FILE; \
	echo "" >> $$UP_FILE; \
	echo "-- Write your migration here" >> $$UP_FILE; \
	echo "" >> $$UP_FILE; \
	echo "-- +migrate Down" > $$DOWN_FILE; \
	echo "" >> $$DOWN_FILE; \
	echo "-- Write your rollback here" >> $$DOWN_FILE; \
	echo ""; \
	echo "Migration created successfully:"; \
	echo "  $$UP_FILE"; \
	echo "  $$DOWN_FILE";

DB_URL=postgres://postgres:postgres@localhost:5433/restaurant_platform?sslmode=disable

migrate-up:
	migrate -path database/schema -database "$(DB_URL)" up

migrate-down:
	migrate -path database/schema -database "$(DB_URL)" down 1

migrate-status:
	migrate -path database/schema -database "$(DB_URL)" version

	# Start all infrastructure (DB + Temporal)
infra-up:
	docker start restaurant_platform_db
	docker compose -f docker-compose.temporal.yml up -d

# Stop all infrastructure
infra-down:
	docker stop restaurant_platform_db
	docker compose -f docker-compose.temporal.yml down