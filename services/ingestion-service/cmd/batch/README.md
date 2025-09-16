# Batch Processing Commands

This directory contains batch processing commands for the YouTube Analytics ingestion service.

## Available Commands

### 1. Trending Video Collection (`trending`)
Collects trending videos from YouTube API and filters them by genre-specific keywords.

```bash
# Collect for all enabled genres
go run ./cmd/batch/trending/main.go

# Collect for specific genre
go run ./cmd/batch/trending/main.go -genre 550e8400-e29b-41d4-a716-446655440001

# Dry run mode
go run ./cmd/batch/trending/main.go -dry-run
```

### 2. Snapshot Scheduling (`schedule-snapshots`)
Schedules snapshot tasks for videos at checkpoints (3h, 6h, 12h, 24h, 48h, 72h, 7d).
Note: Actual snapshot creation and metrics calculation are handled by the task queue handler.

```bash
# Schedule for videos from last 24 hours
go run ./cmd/batch/schedule-snapshots/main.go

# Schedule for videos from last 48 hours
go run ./cmd/batch/schedule-snapshots/main.go -hours 48
```

### 3. WebSub Subscription Renewal (`websub-renewal`)
Renews expiring WebSub subscriptions for channel monitoring.

```bash
# Renew subscriptions expiring in 7 days
go run ./cmd/batch/websub-renewal/main.go

# Renew subscriptions expiring in 3 days
go run ./cmd/batch/websub-renewal/main.go -days 3
```

### 4. Rankings Generation (`rankings`)
Generates daily rankings based on video metrics.

```bash
# Generate rankings for all genres
go run ./cmd/batch/rankings/main.go

# Generate ranking for specific genre
go run ./cmd/batch/rankings/main.go -genre GENRE_ID

# Generate top 20 instead of top 10
go run ./cmd/batch/rankings/main.go -top 20

# Use different checkpoint (default: 24h)
go run ./cmd/batch/rankings/main.go -checkpoint 48
```

## Using the Makefile

Batch commands are integrated in the main Makefile:

```bash
# Run commands
make batch-trending
make batch-trending-genre GENRE_ID=xxx
make batch-schedule-snapshots HOURS=48
make batch-websub-renewal DAYS=3
make batch-rankings TOP=20
make batch-daily  # Runs all batches in sequence
```

## Docker Deployment

Build and run batch jobs in Docker:

```bash
# Build the batch image
docker build -f deployments/batch/Dockerfile -t youtube-analytics-batch .

# Run a specific batch job
docker run --env-file .env youtube-analytics-batch /app/bin/batch-trending

# Run with cron
docker run -d \
  --name youtube-analytics-cron \
  --env-file .env \
  -v $(pwd)/deployments/batch/crontab:/etc/crontabs/appuser \
  youtube-analytics-batch crond -f
```

## Environment Variables

All batch commands require the following environment variables:

- `DATABASE_URL`: PostgreSQL connection string
- `YOUTUBE_API_KEY`: YouTube Data API v3 key
- `PUBSUB_PROJECT_ID`: Google Cloud Pub/Sub project ID
- `CLOUDTASKS_PROJECT_ID`: Google Cloud Tasks project ID
- `CLOUDTASKS_LOCATION`: Cloud Tasks location (e.g., us-central1)
- `CLOUDTASKS_QUEUE_NAME`: Cloud Tasks queue name
- `CLOUDTASKS_SERVICE_URL`: URL for snapshot task handler
- `WEBSUB_CALLBACK_URL`: WebSub callback URL for subscriptions

## Scheduling

Recommended schedule (see `deployments/batch/crontab`):

- **Trending Collection**: Twice daily at 3:00 AM and 3:00 PM
- **Snapshot Scheduling**: Every hour (for recent videos)
- **Rankings Generation**: Daily at 6:00 AM
- **WebSub Renewal**: Daily at 1:00 AM

## Monitoring

Each batch command logs:
- Start time and parameters
- Progress updates
- Results summary (items processed, created, updated)
- Total execution time
- Errors with context

Use structured logging aggregation (e.g., Stackdriver) to monitor batch job health.