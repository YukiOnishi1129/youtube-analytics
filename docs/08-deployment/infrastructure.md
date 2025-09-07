# Deployment & Resources

## Cloud Run

- minInstances=0 (optimizing for free tier), min=1 if needed (~$10/month)
- WebSub only minInstances=1 (requires always-on)
- Unauthenticated calls disabled (grant Invoker permissions via IAM)

## YouTube API

- Primarily videos.list / channels.list → plenty of room within 10k unit/day quota

## Other

- Cloud Tasks: nearly possible to operate within free tier
- Neon (free tier) + sqlc/goose operation

## Infrastructure & Deployment Details

**Dockerfile**: Placed directly under each service (production + development)

**Cloud Run**:
- minInstances=0 (WebSub only 1)
- Unauthenticated calls disabled (grant Invoker permissions via IAM)

**Secrets**: Environment injection via Secret Manager

**CI/CD (GitHub Actions)**: Build/deploy workflow for each service

**Observability**: OpenTelemetry, Prometheus, Zap structured logging
