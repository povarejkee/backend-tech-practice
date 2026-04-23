include .env
export

export PROJECT_ROOT=$(shell pwd)

env-up:
	docker compose up -d app-postgres

env-down:
	docker compose down app-postgres

env-cleanup:
	@read -p "Do you want to clear all the volume files? Your data could be lost. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down app-postgres && \
		rm -rf out/pgdata && \
		echo "Files has been cleared"; \
	else  \
		echo "Volume cleaning has been cancelled"; \
	fi