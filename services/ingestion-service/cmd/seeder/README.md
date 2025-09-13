# Ingestion Service Seeder

A command-line tool for seeding the ingestion service database with initial data.

## Usage

### Local Development

```bash
# Run all seeds
make seed

# Run only keyword seeds
make seed-keywords

# Preview SQL without executing (dry run)
make seed-dry-run

# Build seeder binary
make build-seeder
./bin/seeder -target=all
```

### Docker

```bash
# Build Docker image
make docker-build-seeder

# Run with Docker
docker run --rm \
  -e DATABASE_URL="postgres://user:pass@host:5432/db" \
  ingestion-seeder:latest \
  -target=all
```

### Kubernetes Job

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

### Cloud Run Job

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

## Options

- `-target`: Seed target (default: "all")
  - `all`: Run all seeds
  - `keywords`: Run only keyword seeds
- `-dry-run`: Show SQL without executing

## Seed Data

### Keywords

Default keywords for video filtering (Japanese programming content):
- Programming Languages with Japanese terms (JavaScript/ジャバスクリプト, Python/パイソン, Go/ゴー言語, etc.)
- Japanese programming terms (プログラミング, エンジニア, Web開発, アプリ開発)
- Cloud Platforms with Japanese terms (AWS/アマゾンウェブサービス, GCP/グーグルクラウド, etc.)
- DevOps Tools with Japanese terms (Docker/ドッカー, Kubernetes/クバネティス, etc.)
- Japanese tutorial terms (入門, 解説, ハンズオン)

## Environment Variables

- `DATABASE_URL`: PostgreSQL connection string
- Or individual components:
  - `DB_HOST`
  - `DB_PORT`
  - `DB_USER`
  - `DB_PASSWORD`
  - `DB_NAME`
  - `DB_SSLMODE`