APP_NAME=runlog-api

.PHONY: help

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