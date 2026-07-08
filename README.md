# runlog-api

Simple GO Rest service using clean architecture design
The usecase :
- to store any running activities
- to show any running activities via Dashboard UI(TBD)

## Structure
- `cmd/api`: application entrypoint
- `internal/config`: configuration loading
- `internal/handler`: domain handler to receieve any request
- `internal/usecase`: application use cases

## Run

```sh
go run ./cmd/api
```

## Docker Compose

Start any application dependencies:
```sh
docker compose up -d
```

## Apply any data migration
```sh
docker exec -i runlog-postgres psql -U runlog -d runlog_db < migrations/001_create_users_table.sql
docker exec -i runlog-postgres psql -U runlog -d runlog_db < migrations/002_create_activities_table.sql

-- verify
docker exec -it runlog-postgres psql -U runlog -d runlog_db
```