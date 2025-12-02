# Environment files
ENV_FILE=.env
TEST_ENV_FILE=.env.test

# Docker compose command and project names
DC=docker compose
DEV_PROJECT=bike-map-dev
TEST_PROJECT=bike-map-test

.PHONY: help dev-build dev-start dev-stop dev-restart dev-fl dev-bl dev-fr dev-br \
        test-build test-start test-stop test-restart test-fl test-bl test-fr test-br

help:
	@echo "Available targets: dev-build dev-start dev-stop dev-restart dev-fl dev-bl dev-fr dev-br"
	@echo "                 test-build test-start test-stop test-restart test-fl test-bl test-fr test-br"

# Dev targets
dev-build:
	$(DC) --env-file $(ENV_FILE) -f docker-compose.yml -p $(DEV_PROJECT) build --no-cache

dev-start:
	$(DC) --env-file $(ENV_FILE) -f docker-compose.yml -p $(DEV_PROJECT) up -d

dev-stop:
	$(DC) -p $(DEV_PROJECT) down

dev-restart: dev-stop dev-build dev-start

dev-fl:
	$(DC) -p $(DEV_PROJECT) logs dev-frontend

dev-bl:
	$(DC) -p $(DEV_PROJECT) logs dev-backend

dev-fr:
	$(DC) -p $(DEV_PROJECT) restart dev-frontend

dev-br:
	$(DC) -p $(DEV_PROJECT) restart dev-backend

# Test targets
test-build:
	$(DC) --env-file $(TEST_ENV_FILE) -f docker-compose.test.yml -p $(TEST_PROJECT) build --no-cache

test-start:
	$(DC) --env-file $(TEST_ENV_FILE) -f docker-compose.test.yml -p $(TEST_PROJECT) up --force-recreate -d

test-stop:
	$(DC) -p $(TEST_PROJECT) down

test-restart: test-stop test-build test-start

test-fl:
	$(DC) -p $(TEST_PROJECT) logs test-frontend

test-bl:
	$(DC) -p $(TEST_PROJECT) logs test-backend

test-fr:
	$(DC) -p $(TEST_PROJECT) restart test-frontend

test-br:
	$(DC) -p $(TEST_PROJECT) restart test-backend
