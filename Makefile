APP_NAME=runlog-api

.PHONY: help

DB_CONTAINER=runlog-postgres
DB_USER=runlog
DB_NAME=runlog_db
MIGRATIONS_DIR=migrations

help:
	@echo "Available commands:"
	@echo "  make run"
	@echo "  make test"
	@echo "  make fmt"
	@echo "  make tidy"
	@echo "  make docker-build"
	@echo "  make docker-run"
	@echo "  make db-up"
	@echo "  make db-down"

run:
	go run ./cmd/api

fmt:
	go fmt ./...

test:
	go test ./...

tidy:
	go mod tidy

docker-build:
	docker build -t $(APP_NAME) .

docker-run:
	docker run \
	-p 8080:8080 \
	-e APP_PORT=8080 \
	-e DATABASE_URL="postgres://runlog:runlog_password@host.docker.internal:5432/runlog_db?sslmode=disable" \
	-e JWT_SECRET="super_secret_change_later" \
	-e JWT_EXPIRY_HOURS=24 \
	$(APP_NAME)

db-up:
	docker compose up -d

db-down:
	docker compose down

db-logs:
	docker compose logs -f

clean:
	go clean

lint:
	golangci-lint run

deps:
	go mod download

migrate-up:
	for file in $(MIGRATIONS_DIR)/*.sql; do \
		echo "Applying $$file"; \
		docker exec -i $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME) < $$file; \
	done

db-shell:
	docker exec -it $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME)

db-reset:
	docker compose down -v
	docker compose up -d
	sleep 3
	make migrate-up

vet:
	go vet ./...

build:
	go build -o bin/$(APP_NAME) ./cmd/api

ci: fmt vet test build docker-build