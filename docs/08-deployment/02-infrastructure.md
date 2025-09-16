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

Example commands (using DATABASE_URL):
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
- DB configuration: either `DATABASE_URL` or the individual envs:
  - `DB_HOST`, `DB_PORT`(default 5432), `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`(default disable)
- Build flags control wiring:
  - `-tags 'sqlc'` enables sqlc-backed repositories
  - If not built with `-tags sqlc`, the driver returns an error to prevent accidental startup

## CI Tests

- Use `testcontainers-go` to run Postgres for adapter/integration tests.
- Apply migrations in test setup using golang-migrate before running tests.
