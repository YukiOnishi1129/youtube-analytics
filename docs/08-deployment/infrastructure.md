# Deployment & Resources

## Cloud Run

- minInstances=0 (optimizing for free tier), min=1 if needed (~$10/month)
- WebSub only minInstances=1 (requires always-on)
- Unauthenticated calls disabled (grant Invoker permissions via IAM)

## YouTube API

- Primarily videos.list / channels.list → plenty of room within 10k unit/day quota

## Other

- Cloud Tasks: nearly possible to operate within free tier
- Neon (free tier) + sqlc/migrate (golang-migrate) operation

## Infrastructure & Deployment Details

**Dockerfile**: Placed directly under each service (production + development)

**Cloud Run**:
- minInstances=0 (WebSub only 1)
- Unauthenticated calls disabled (grant Invoker permissions via IAM)

**Secrets**: Environment injection via Secret Manager

**CI/CD (GitHub Actions)**: Build/deploy workflow for each service

**Observability**: OpenTelemetry, Prometheus, Zap structured logging

## Database Migrations (golang-migrate)

- Tool: https://github.com/golang-migrate/migrate
- Policy: schema-per-service, no cross-schema joins
- Authority migrations: `services/authority-service/internal/driver/datastore/migrations`

Example commands:
```
migrate -path services/authority-service/internal/driver/datastore/migrations \
  -database "$DATABASE_URL" up

migrate -path services/authority-service/internal/driver/datastore/migrations \
  -database "$DATABASE_URL" down 1
```

Or use the service Makefile:

```
cd services/authority-service
DATABASE_URL=postgres://user:pass@localhost:5432/app?sslmode=disable make migrate-up
DATABASE_URL=... make migrate-down
```

## Local Development

- Prefer Docker Compose for Postgres during local development.
- Each service connects to the same instance but uses a dedicated schema.
- Build flags control wiring:
  - `-tags 'postgres sqlc'` enables pgx driver and postgres repository (sqlc)
  - otherwise, the service refuses to start if repositories are not wired

## CI Tests

- Use `testcontainers-go` to run Postgres for adapter/integration tests.
- Apply migrations in test setup using golang-migrate before running tests.
