# runlog-api

Runlog API is a Go REST API for storing and managing running activity logs. It provides user registration, login with JWT authentication, profile lookup, and authenticated CRUD endpoints for activities.

## Tech Stack

- Go 1.25
- Gin HTTP framework
- PostgreSQL
- pgx database driver
- JWT authentication
- Docker and Docker Compose for local infrastructure

## Project Structure

```text
cmd/api                 Application entrypoint and route registration
internal/config         Environment and PostgreSQL configuration
internal/handler        HTTP handlers
internal/model          Domain models
internal/repository     PostgreSQL data access
internal/usecase        Application business logic
migrations              SQL schema migrations
pkg/auth                JWT helpers
pkg/dto                 Request and response DTOs
pkg/errors              Shared domain errors
pkg/logging             Structured logger setup
pkg/middleware          Auth, CORS, request ID, and request logging middleware
pkg/response            Shared JSON response helpers
```

## Prerequisites

- Go 1.25 or newer
- Docker and Docker Compose
- `make`

## Environment

Create a `.env` file in the project root:

```env
APP_PORT=8080
DATABASE_URL=postgres://runlog:runlog_password@localhost:5432/runlog_db?sslmode=disable
JWT_SECRET=super_secret_change_later
JWT_EXPIRY_HOURS=24
CORS_ALLOW_ORIGINS=http://localhost:3000,http://localhost:5173
```

| Variable | Default | Description |
| --- | --- | --- |
| `APP_PORT` | `8080` | HTTP server port |
| `DATABASE_URL` | empty | PostgreSQL connection string |
| `JWT_SECRET` | `dev_secret` | Secret used to sign JWT tokens |
| `JWT_EXPIRY_HOURS` | `24` | JWT expiration time in hours |
| `CORS_ALLOW_ORIGINS` | `http://localhost:3000,http://localhost:5173` | Comma-separated allowed frontend origins |

## Local Development

Start PostgreSQL:

```sh
make db-up
```

Apply database migrations:

```sh
make migrate-up
```

Run the API:

```sh
make run
```

The API will be available at:

```text
http://localhost:8080
```

Check the service:

```sh
curl http://localhost:8080/health
```

## Make Commands

| Command | Description |
| --- | --- |
| `make help` | Show available commands |
| `make run` | Run the API with `go run ./cmd/api` |
| `make test` | Run all Go tests |
| `make fmt` | Format Go code |
| `make tidy` | Tidy Go modules |
| `make deps` | Download Go modules |
| `make lint` | Run `golangci-lint` |
| `make docker-build` | Build the API Docker image |
| `make docker-run` | Run the API Docker image |
| `make db-up` | Start PostgreSQL with Docker Compose |
| `make db-down` | Stop Docker Compose services |
| `make db-logs` | Follow Docker Compose logs |
| `make db-shell` | Open a PostgreSQL shell |
| `make migrate-up` | Apply all SQL migrations |
| `make db-reset` | Recreate the database volume and apply migrations |
| `make clean` | Run `go clean` |

## API Endpoints

### Public

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/health` | Health check |
| `POST` | `/api/v1/users/register` | Register a user |
| `POST` | `/api/v1/users/login` | Login and receive a JWT |

### Authenticated

Authenticated routes require this header:

```text
Authorization: Bearer <token>
```

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/api/v1/users/profile` | Get the authenticated user profile |
| `POST` | `/api/v1/activities` | Create an activity |
| `GET` | `/api/v1/activities` | List authenticated user's activities |
| `GET` | `/api/v1/activities/:id` | Get one activity |
| `PUT` | `/api/v1/activities/:id` | Update one activity |
| `DELETE` | `/api/v1/activities/:id` | Delete one activity |

## Request Examples

Register:

```sh
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Run Logger",
    "email": "runner@example.com",
    "password": "password123"
  }'
```

Login:

```sh
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "runner@example.com",
    "password": "password123"
  }'
```

Create an activity:

```sh
curl -X POST http://localhost:8080/api/v1/activities \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "title": "Morning run",
    "sport_type": "running",
    "distance_km": 5.25,
    "duration_seconds": 1800,
    "elevation_gain_m": 35,
    "activity_date": "2026-07-08",
    "notes": "Easy pace"
  }'
```

List activities:

```sh
curl http://localhost:8080/api/v1/activities \
  -H "Authorization: Bearer <token>"
```

Update an activity:

```sh
curl -X PUT http://localhost:8080/api/v1/activities/<activity_id> \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "title": "Tempo run",
    "sport_type": "running",
    "distance_km": 6.00,
    "duration_seconds": 1920,
    "elevation_gain_m": 40,
    "activity_date": "2026-07-08",
    "notes": "Faster finish"
  }'
```

Delete an activity:

```sh
curl -X DELETE http://localhost:8080/api/v1/activities/<activity_id> \
  -H "Authorization: Bearer <token>"
```

## Response Format

Successful responses return data under a `data` key:

```json
{
  "data": {
    "id": "uuid",
    "title": "Morning run"
  }
}
```

Errors return an `error` key:

```json
{
  "error": "invalid request body"
}
```

Delete requests return `204 No Content` when successful.

## Database

The local PostgreSQL service is defined in `docker-compose.yml`:

- container: `runlog-postgres`
- database: `runlog_db`
- user: `runlog`
- password: `runlog_password`
- port: `5432`

To open a database shell:

```sh
make db-shell
```

To reset local database state:

```sh
make db-reset
```

## Docker

Build the API image:

```sh
make docker-build
```

Run the API image:

```sh
make docker-run
```
