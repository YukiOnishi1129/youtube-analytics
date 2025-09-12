# Ingestion Service Makefile

.PHONY: help sqlc sqlc-clean build-postgres migrate-up migrate-down generate-http

help:
	@echo "Targets:"
	@echo "  sqlc          Generate typed queries (sqlc generate)"
	@echo "  sqlc-clean    Remove generated sql (sqlcgen)"
	@echo "  build-postgres Build with Postgres+sqlc tags"
	@echo "  migrate-up    Run DB migrations (golang-migrate)"
	@echo "  migrate-down  Rollback one migration"
	@echo "  generate-http Generate HTTP types from TypeSpec"

# Generate sqlc code into internal/adapter/gateway/postgres/sqlcgen
sqlc:
	@echo "==> sqlc generate"
	@command -v sqlc >/dev/null 2>&1 || { echo "sqlc not found. Install with: brew install sqlc"; exit 1; }
	@sqlc generate

sqlc-clean:
	@echo "==> cleaning sqlc generated code"
	rm -rf internal/adapter/gateway/postgres/sqlcgen

# Convenience: ensure Postgres driver and sqlc code are included
build-postgres:
	@echo "==> go build with tags 'postgres sqlc'"
	GOFLAGS="-tags 'postgres sqlc'" go build ./...

# Migration commands (requires golang-migrate installed)
MIGRATIONS_DIR ?= internal/driver/datastore/migrations
DATABASE_URL   ?= $(DATABASE_URL)

migrate-up:
	@echo "==> migrate up"
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

migrate-down:
	@echo "==> migrate down 1"
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down 1

# Generate HTTP types from TypeSpec
generate-http:
	@echo "==> Generating HTTP types from TypeSpec"
	@bash scripts/generate-http-types.sh