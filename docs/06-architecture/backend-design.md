# Backend Architecture

## Design Philosophy

### Clean Architecture
- Dependency direction is inward: domain ← usecase ← adapter ← driver
- Layer responsibilities:
  - domain: Business rules (Entity, ValueObject, DomainService)
  - usecase: Application logic (Interactor, Application Service)
  - port/input: Public interfaces of use cases
  - port/output/gateway: Abstractions for DB/external APIs
  - port/output/presenter: Output boundary
  - adapter: Concrete implementations of controller/presenter/gateway
  - driver: Bootstrapping, config, security, observability
  - cmd: Composition root

### DDD (Domain-Driven Design)
- Ubiquitous Language: Map domain concepts directly to Entity/VO/Service
- Aggregates: Video, Channel, Keyword are primary aggregates
- Repository: Persistence abstraction defined under port/output/gateway
- Domain Services: Centralize metric calculations (growth, Wilson lower bound, LPS, etc.) in domain/services

## Microservice Split (Bounded Contexts)

- ingestion-service
  - Collect videos via YouTube WebSub/Trending
  - Persist snapshots (D0, 6h, 24h… via Cloud Tasks)
  - Manage keywords
- analytics-service
  - Precompute metrics from snapshots
  - Serve rankings (RankingKind, CheckpointHour)
  - Freeze and browse history
- authority-service
  - Integrate with Identity Platform
  - Verify ID tokens (JWT/JWKS)
  - Manage profiles/roles

Each service has its own DB schema and go.mod; services interact via gRPC.

## AuthN/Z Details

User API: Each service verifies OIDC (Identity Platform) locally
- Cache JWKS
- Validate iss/aud/exp/signature

Internal API (Cloud Tasks/Scheduler): Cloud Run IAM + OIDC
- Invoke with a service account
- Attach OIDC token (aud = service URL)
- Optionally double-verify via OIDC middleware in the app

Per-method policy: PUBLIC / USER_ID_TOKEN / SERVICE_OIDC
Shared implementation is centralized in `services/pkg/identityauth` and imported by all services

## Repository Layout

```
youtube-analytics/
├─ proto/                        # .proto definitions (buf)
├─ services/
│   ├─ go.work                   # workspace
│   ├─ pkg/
│   │   ├─ identityauth/         # shared auth
│   │   └─ pb/                   # buf generate (Go)
│   ├─ ingestion-service/        # go.mod, Dockerfile, internal/...
│   ├─ analytics-service/
│   └─ authority-service/
└─ web/
    └─ client/                   # Next.js frontend
```

## Docker & Deployment

- Dockerfile per service (prod and dev)
- Cloud Run: minInstances=0 (WebSub only = 1), disallow unauthenticated invocations
- CI/CD (GitHub Actions): Build/deploy jobs per service

## Testing Strategy

- Unit: Pure functions in domain
- UseCase: Swap port/output with fakes/mocks
- Component/Integration: handler → usecase → repo; DB via testcontainers-go
- Contract: buf breaking (proto schema checks)
- Service-level E2E: A few gRPC/REST tests per service surface
- System-level E2E: Limit to staging smoke tests

## Directory Structure (per service)

```
services/
├─ go.work
├─ pkg/
│   ├─ identityauth/   # shared OIDC verification
│   └─ pb/             # buf generate output (Go)
├─ ingestion-service/
│   ├─ go.mod
│   ├─ cmd/server/main.go
│   └─ internal/
│       ├─ domain/
│       ├─ port/input/
│       ├─ port/output/{gateway,presenter}/
│       ├─ usecase/
│       ├─ adapter/{controller,gateway,presenter}/
│       └─ driver/{config,transport,datastore,security,observability,health}
├─ analytics-service/
│   └─ ...
└─ authority-service/
    └─ ...
```

proto definitions live at repository root `/proto`.
buf generate outputs: Go → `/services/pkg/pb`, TS → `/web/client/src/external/client/grpc`.

## Batch Processing & Idempotency

- Cloud Tasks → `/snapshot`: TaskID = `snap:{videoId}:{cp}`
- DB: Use PK/UNIQUE + INSERT ... ON CONFLICT for safe retries
- Idempotency tests are required (multiple identical inputs → only one effective result)

## Detailed Testing

- Unit: Pure domain logic
- UseCase: Replace port/output with Fakes/Mocks
- Component/Integration: handler→usecase→repo with testcontainers-go for DB
- Contract: buf breaking checks for proto
- Service-level E2E: A small set via gRPC/REST
- System-level E2E: Staging smoke only; skip in CI

## Infra & Deployment

Dockerfile: Place at each service (prod + dev)
Cloud Run:
- minInstances=0 (1 only for WebSub)
- Disallow unauthenticated (grant IAM Invoker)

Secrets: Inject env via Secret Manager
CI/CD (GitHub Actions): Build/deploy workflows per service
Observability: OpenTelemetry, Prometheus, structured logging with Zap

## Frontend Boundary (proto/TS SDK)

proto definitions: `/proto`
buf generate:
- Go → `/services/pkg/pb` (internal import only)
- TS → `/web/client/src/external/client/grpc` (imported by Next.js)

Place TS artifacts under `external` and do not import them directly from features:
features → handlers (Server Actions) → services → client(grpc)
