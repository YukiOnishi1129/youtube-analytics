# Database Seeder

Database seeder for YouTube Analytics ingestion service master data.

## Master Data

All seeds are production-ready master data:

- **youtube_categories**: Official YouTube API categories (27: Education, 28: Science & Technology, etc.)
- **genres**: Japanese Engineering genre only (combining categories 27 & 28)
- **keywords**: Search patterns for Japanese engineering content

## Usage

```bash
# Run all seeds
make seed

# Run specific seeds
make seed-categories
make seed-genres  
make seed-keywords

# Preview SQL without executing
make seed-dry-run
```

## Seed Order

Seeds are executed in dependency order:
1. youtube_categories
2. genres (depends on categories)
3. keywords (depends on genres)

## Initial Launch Configuration

For the initial launch, we're focusing on:
- **Single Genre**: Japanese Engineering (エンジニア)
- **Categories**: 27 (Education) + 28 (Science & Technology)
- **Keywords**: ~20 patterns covering programming languages, frameworks, and Japanese tech terms

## Docker

```bash
# Build Docker image
make docker-build-seeder

# Run with Docker
docker run --rm \
  -e DATABASE_URL="postgres://user:pass@host:5432/db" \
  ingestion-seeder:latest \
  -target=all
```

## Kubernetes Job

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: ingestion-seeder
spec:
  template:
    spec:
      containers:
      - name: seeder
        image: ingestion-seeder:latest
        args: ["-target=all"]
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: ingestion-db-secret
              key: url
      restartPolicy: Never
```

## Cloud Run Job

```bash
# Deploy as Cloud Run Job
gcloud run jobs create ingestion-seeder \
  --image=gcr.io/your-project/ingestion-seeder:latest \
  --args="-target=all" \
  --set-env-vars="DATABASE_URL=$DATABASE_URL" \
  --region=us-central1

# Execute the job
gcloud run jobs execute ingestion-seeder --region=us-central1
```

## Environment Variables

- `DATABASE_URL`: PostgreSQL connection string
- Or individual DB components:
  - `DB_HOST`
  - `DB_PORT`
  - `DB_USER`
  - `DB_PASSWORD`
  - `DB_NAME`
  - `DB_SSLMODE`

## Idempotency

All seeds use `ON CONFLICT DO UPDATE` to ensure they can be run multiple times safely.