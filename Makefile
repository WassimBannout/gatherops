APP_NAME := gatherops
MIGRATIONS_PATH ?= migrations
MIGRATE_DOWN_STEPS ?= 1

-include .env
.EXPORT_ALL_VARIABLES:

.PHONY: run test test-integration vet lint openapi-check docker-up docker-down migrate-up migrate-down

run:
	go run ./cmd/api

test:
	go test ./...

test-integration:
	GATHEROPS_INTEGRATION_TESTS=1 go test -count=1 ./test/integration/...

vet:
	go vet ./...

lint: vet

openapi-check:
	test -f docs/openapi.yaml

docker-up:
	docker compose up -d postgres

docker-down:
	docker compose down

migrate-up:
	go run ./cmd/migrate -direction up -path $(MIGRATIONS_PATH)

migrate-down:
	go run ./cmd/migrate -direction down -path $(MIGRATIONS_PATH) -steps $(MIGRATE_DOWN_STEPS)
