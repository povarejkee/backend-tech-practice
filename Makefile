include .env
export

export PROJECT_ROOT=$(shell pwd)

env-up:
	@docker compose up -d app-postgres

env-down:
	@docker compose down app-postgres

env-cleanup:
	@read -p "Do you want to clear all the volume files? Your data could be lost. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down app-postgres port-forwarder && \
		rm -rf out/pgdata && \
		echo "Files has been cleared"; \
	else  \
		echo "Volume cleaning has been cancelled"; \
	fi

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down -d port-forwarder

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "'seq' param is required! Sample: make migrate-create seq=test"; \
		exit 1; \
	fi; \

	docker compose run --rm app-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"
	
migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "'action' param is required! Sample: make migrate-action action=down"; \
		exit 1; \
	fi; \

	docker-compose run --rm app-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@app-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

app-start:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run cmd/app/main.go