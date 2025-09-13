# Backend Architecture — Final Plan (Clean Architecture + DDD)

## 1) Overall Strategy
- Clean Architecture (dependencies point inward)
  - domain (Entity/VO/Domain Service)
  - port/input (UseCase public interface)
  - port/output/gateway (DB/external API abstraction)
  - port/output/presenter (output boundary)
  - usecase (Interactor: Input Port implementation)
  - adapter (concrete controller/presenter/gateway)
  - driver (bootstrap/config/connectivity/security/observability)
  - cmd (Composition Root)
- DDD: Aggregates = keyword / video / metric / channel / account
- One aggregate = one package; split files by concept (e.g., keyword.go, keyword_service.go)

## 2) AuthN/Z (Identity Platform + Cloud Run)
- User APIs: each service validates ID tokens itself (go-oidc, JWKS cache)
- Internal APIs (Cloud Tasks/Scheduler): allow Cloud Run IAM only + double-check with OIDC
- Method policy: PUBLIC / USER_ID_TOKEN / SERVICE_OIDC (enforced via interceptor)
  - Implement OIDC verification and gRPC interceptors in each service under `internal/adapter` / `internal/driver/security` (standardize using authority-service as the template)

## 3) Directory Layout (Monorepo / go.work)

```
youtube-analytics/
├─ proto/                              # .proto (managed by buf)
├─ services/
│   ├─ go.work
│   ├─ pkg/
│   │   └─ pb/
│   ├─ ingestion-service/
│   ├─ analytics-service/
│   └─ authority-service/
└─ web/
```

Packages with the same name (e.g., keyword) can exist per layer. Use import aliases to distinguish them (e.g., domKeyword, ucKeyword, inKeyword).

## 4) Proto & Generated Artifacts
- Use buf generate
- Go → `services/pkg/pb`
- TS → `web/client/src/external/client/grpc`

## 5) Docker / CI
- Dockerfile per service (prod/dev)
- Cloud Run: minInstances=0 (1 if needed), disallow unauthenticated invocations, inject secrets via Secret Manager → env
- GitHub Actions: build/deploy per service (matrix possible)

## 6) Testing Strategy
- Unit (domain) / UseCase (swap ports) / Component (handler→usecase→repo)
- Contract (buf breaking) / Service E2E (few) / System E2E (staging smoke)
- Idempotency: repeated identical input yields a single effect

## 7) Batch & Idempotency
- Cloud Tasks → `/snapshot` (ingestion)
- TaskID = `snap:{videoId}:{cp}` / DB uses UNIQUE + upsert

## 8) Domain (Excerpt)
- Keyword: name, filterType, pattern, enabled / PatternBuilder.Build()
- Channel + ChannelSnapshot (subscriber trend)
- Video + VideoSnapshot (0/3/6/12/24/48/72/168h)
- Metric (read model per checkpoint: growth/ratio/Wilson/LPS/Exclude)
- History (freeze TopN)
- Account / Identity / Role (email unique; identity duplication prohibited; inactive cannot log in)

## 9) Implementation Rules (readability first)
- One domain = one package; split files by concept
- Enforce invariants in constructors; model state transitions as methods
- Reference across aggregates by ID
- Repository abstractions live in `port/output/gateway` (owned by UseCase)
- Use import aliases to avoid confusion with same-named packages

## 10) Configuration (minimum)
- OIDC: `IDP_ISSUER`, `IDP_AUDIENCE`, `JWKS_CACHE_TTL`
- Internal API: `TASKS_AUDIENCE`, `INTERNAL_SECRET`
- Metrics: `LIKES_PER_SUBSCRIPTION_SCALE=1000`, `LIKES_PER_SUBSCRIPTION_OFFSET=500`
