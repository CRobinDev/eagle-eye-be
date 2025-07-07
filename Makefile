# Makefile for EagleEye Project

# Load environment variables from .env if present
-include .env

# Variables
POSTGRES_URL := postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
MIGRATION_DIR := database/migration
SEEDER_DIR := database/seeder

SEED_CMD := docker compose exec -T db psql -U $(DB_USER) -d $(DB_NAME) <
MIGRATE_CMD := docker compose run --rm migrate -database "$(POSTGRES_URL)" -path $(MIGRATION_DIR)
CREATE_MIGRATE_CMD := docker compose run --rm migrate create -ext sql -dir $(MIGRATION_DIR) 


LOCAL_MIGRATE_CMD := migrate -database "$(POSTGRES_URL)" -path $(MIGRATION_DIR)
LOCAL_CREATE_MIGRATE_CMD := migrate create -ext sql -dir $(MIGRATION_DIR)

# Help
help:
	@echo "Available commands: "

	@echo "  migrate-up         Apply all migrations"
	@echo "  migrate-down       Rollback the most recent migration"
	@echo "  migrate-redo       Revert and reapply the most recent migration"
	@echo "  migrate-status     Show migration status"
	@echo "  migrate-version    Show current migration version"
	@echo "  migrate-force      Force database to a specific version (use 'make migrate-force version=123')"
	@echo "  migrate-create     Create a new migration (use 'make migrate-create name=your_name')"
	@echo "  seed-up            Seed database with up SQL (use 'make seed-up file=seed_file.sql')"
	@echo "  seed-down          Seed database with down SQL (use 'make seed-down file=seed_file.sql')"
	@echo "  run                Start all services and run migrations"
	@echo "  stop               Stop all services"
	@echo "  build              Build all Docker images"
	@echo "  logs               Show logs for all containers"

# Migration commands
migrate-local:
	$(LOCAL_MIGRATE_CMD) up

migrate-local-d:
	$(LOCAL_MIGRATE_CMD) down

migrate-flocal:
	$(LOCAL_MIGRATE_CMD) force $(version)

create-local:
	$(LOCAL_CREATE_MIGRATE_CMD) $(name)
	
migrate-up:
	$(MIGRATE_CMD) up

migrate-down:
	$(MIGRATE_CMD) down

migrate-redo:
	$(MIGRATE_CMD) redo

migrate-status:
	$(MIGRATE_CMD) status

migrate-version:
	$(MIGRATE_CMD) version

migrate-force:
	$(MIGRATE_CMD) force $(version)

migrate-create:
	$(CREATE_MIGRATE_CMD) $(name)

# Seeder commands
seed-up:
	$(SEED_CMD) $(SEEDER_DIR)/$(file)

seed-down:
	$(SEED_CMD) $(SEEDER_DIR)/$(file)

# Docker Compose commands
run:
	docker compose up -d $(ARGS)
	$(MAKE) migrate-up

stop:
	docker compose down $(ARGS)

build:
	docker compose build

logs:
	docker compose logs -f

gorun:
	go run cmd/app/main.go
.PHONY: help migrate-up migrate-down migrate-redo migrate-status migrate-version migrate-force migrate-create seed-up seed-down run stop build logs

