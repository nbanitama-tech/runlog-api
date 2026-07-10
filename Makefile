APP_NAME=runlog-api

.PHONY: \
	help \
	run fmt tidy vet \
	test test-cover test-cover-func test-cover-html \
	build clean ci lint \
	swagger \
	db-up db-down db-logs db-shell db-reset \
	migrate-up migrate-down migrate-force migrate-create \
	docker-build docker-run docker-tag docker-push \
	test-db-drop test-db-create test-db-migrate \
	test test-integration test-all

DB_CONTAINER=runlog-postgres
DB_USER=runlog
DB_NAME=runlog_db
MIGRATIONS_DIR=migrations
DB_URL=postgres://runlog:runlog_password@localhost:5432/runlog_db?sslmode=
TEST_DB_URL := postgres://runlog:runlog_password@localhost:5432/runlog_test_db?sslmode=disable



#########################################
# Help command
#########################################

help:
	@echo "Available commands:"
	@echo "Development:"
	@echo "  make run"
	@echo "  make test"
	@echo "  make fmt"
	@echo "  make tidy"
	@echo "  make vet"
	@echo "  make ci"
	@echo ""
	@echo "Coverage:"
	@echo "  make test-cover"
	@echo "  make test-cover-html"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build"
	@echo "  make docker-run"
	@echo "  make docker-tag"
	@echo "  make docker-push"
	@echo ""
	@echo "Database:"
	@echo "  make db-up"
	@echo "  make db-down"
	@echo "  make db-logs"
	@echo "  make db-shell"
	@echo "  make db-reset"
	@echo "  make migrate-up"
	@echo "  make migrate-down"
	@echo "  make migrate-force"
	@echo "  make migrate-create"
	@echo "Swagger:"
	@echo "  make swagger"


##########################################
# Development commands
##########################################
run:
	go run ./cmd/api

fmt:
	go fmt ./...

test:
	go test ./...

tidy:
	go mod tidy

deps:
	go mod download

clean:
	go clean

vet:
	go vet ./...

build:
	go build -o bin/$(APP_NAME) ./cmd/api

lint:
	golangci-lint run

ci: fmt lint test build docker-build

###########################################
# Docker commands
###########################################
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

############################################
# Database commands
############################################
db-up:
	docker compose up -d

db-down:
	docker compose down

db-logs:
	docker compose logs -f

migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

migrate-force:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force $(VERSION)

migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(NAME)

db-shell:
	docker exec -it $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME)

db-reset:
	docker compose down -v
	docker compose up -d
	sleep 3
	make migrate-up

############################################
# Testing commands
############################################
test-cover:
	go test ./... -cover

test-cover-html:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

test-db-create:
	docker exec $(DB_CONTAINER) psql -U $(DB_USER) -d postgres \
	-c "CREATE DATABASE runlog_test_db;" || true

test-db-drop:
	docker exec $(DB_CONTAINER) psql -U $(DB_USER) -d postgres \
	-c "DROP DATABASE IF EXISTS runlog_test_db WITH (FORCE);"

test-db-migrate:
	migrate -path $(MIGRATIONS_DIR) -database "$(TEST_DB_URL)" up

test-integration: test-db-drop test-db-create test-db-migrate
	RUN_INTEGRATION_TESTS=true TEST_DATABASE_URL="$(TEST_DB_URL)" go test ./internal/repository -v

test-all: test test-integration

############################################
# Swagger commands
############################################
swagger:
	swag init -g cmd/api/main.go --parseDependency --parseInternal