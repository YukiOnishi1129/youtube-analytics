# Ingestion Service Makefile

.PHONY: help sqlc sqlc-clean build-postgres migrate-up migrate-down generate-http seed

help:
	@echo "== Database & Code Generation =="
	@echo "  sqlc              Generate typed queries (sqlc generate)"
	@echo "  sqlc-clean        Remove generated sql (sqlcgen)"
	@echo "  build-postgres    Build with Postgres+sqlc tags"
	@echo "  migrate-up        Run DB migrations (golang-migrate)"
	@echo "  migrate-down      Rollback one migration"
	@echo "  generate-http     Generate HTTP types from TypeSpec"
	@echo "  seed              Run database seeds"
	@echo ""
	@echo "== Batch Processing Commands =="
	@echo "  batch-trending    Collect trending videos for all enabled genres"
	@echo "  batch-schedule-snapshots  Schedule snapshot tasks for recent videos"
	@echo "  batch-websub-renewal  Renew expiring WebSub subscriptions"
	@echo "  batch-rankings    Generate daily rankings"
	@echo "  batch-daily       Run all batches in sequence"
	@echo ""
	@echo "== Batch Options =="
	@echo "  DRY_RUN=true      Run in dry-run mode (no changes)"
	@echo "  GENRE_ID=xxx      Target specific genre"
	@echo "  HOURS=24          Time window in hours"
	@echo "  DAYS=7            Time window in days"
	@echo "  TOP=10            Number of top videos in ranking"

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

# Run database seeds (all seeds are production-safe master data)
seed:
	@echo "==> Running all database seeds"
	@go run cmd/seeder/main.go -target=all

seed-categories:
	@echo "==> Running YouTube categories seeds"
	@go run cmd/seeder/main.go -target=youtube_categories

seed-genres:
	@echo "==> Running genres seeds"
	@go run cmd/seeder/main.go -target=genres

seed-keywords:
	@echo "==> Running keyword seeds"
	@go run cmd/seeder/main.go -target=keywords


seed-dry-run:
	@echo "==> Showing seed SQL (dry run)"
	@go run cmd/seeder/main.go -target=all -dry-run

# Build seeder binary
build-seeder:
	@echo "==> Building seeder binary"
	@go build -o bin/seeder cmd/seeder/main.go

# Build seeder Docker image
docker-build-seeder:
	@echo "==> Building seeder Docker image"
	@docker build -f Dockerfile.seeder -t ingestion-seeder:latest .

# ===========================================
# Batch Processing Commands
# ===========================================

# Trending video collection
.PHONY: batch-trending
batch-trending:
	go run ./cmd/batch/trending/main.go $(if $(DRY_RUN),-dry-run)

.PHONY: batch-trending-genre
batch-trending-genre:
	@if [ -z "$(GENRE_ID)" ]; then echo "GENRE_ID is required"; exit 1; fi
	go run ./cmd/batch/trending/main.go -genre $(GENRE_ID) $(if $(DRY_RUN),-dry-run)

# Snapshot scheduling
.PHONY: batch-schedule-snapshots
batch-schedule-snapshots:
	go run ./cmd/batch/schedule-snapshots/main.go $(if $(HOURS),-hours $(HOURS)) $(if $(DRY_RUN),-dry-run)

# WebSub renewal
.PHONY: batch-websub-renewal
batch-websub-renewal:
	go run ./cmd/batch/websub-renewal/main.go $(if $(DAYS),-days $(DAYS)) $(if $(DRY_RUN),-dry-run)

# Rankings generation
.PHONY: batch-rankings
batch-rankings:
	go run ./cmd/batch/rankings/main.go $(if $(GENRE_ID),-genre $(GENRE_ID)) $(if $(CHECKPOINT),-checkpoint $(CHECKPOINT)) $(if $(TOP),-top $(TOP)) $(if $(DRY_RUN),-dry-run)

# Daily batch sequence (typically run by cron/scheduler)
.PHONY: batch-daily
batch-daily:
	@echo "Running daily batch sequence..."
	@echo "1. Collecting trending videos..."
	@$(MAKE) batch-trending
	@echo "2. Scheduling snapshots..."
	@$(MAKE) batch-schedule-snapshots
	@echo "3. Generating rankings..."
	@$(MAKE) batch-rankings
	@echo "4. Renewing WebSub subscriptions..."
	@$(MAKE) batch-websub-renewal
	@echo "Daily batch sequence completed!"

# Build all batch commands
.PHONY: build-batch
build-batch:
	go build -o bin/batch-trending ./cmd/batch/trending
	go build -o bin/batch-schedule-snapshots ./cmd/batch/schedule-snapshots
	go build -o bin/batch-websub-renewal ./cmd/batch/websub-renewal
	go build -o bin/batch-rankings ./cmd/batch/rankings
	@echo "All batch commands built to ./bin/"