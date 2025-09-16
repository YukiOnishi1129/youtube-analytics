# System Overview

## Architecture Overview

```
youtube-analytics/
├─ proto/                        # .proto definitions (managed by buf)
├─ services/
│   ├─ go.work                   # Go workspace
│   ├─ pkg/
│   │   ├─ identityauth/         # Shared auth utilities
│   │   └─ pb/                   # buf generate output for Go
│   ├─ ingestion-service/        # Video ingestion, WebSub, keywords
│   ├─ analytics-service/        # Precompute metrics, provide rankings
│   └─ authority-service/        # Identity Platform integration, profiles
└─ web/
    └─ client/                   # Next.js frontend

[ Next.js (App Router, shadcn/ui, TanStack Query) ]
  └─ (gRPC client)
      ├─ authority-service (gRPC)
      ├─ ingestion-service (gRPC) 
      └─ analytics-service (gRPC)

[ ingestion-service (Go, Cloud Run, gRPC + HTTP) ]
  ├─ gRPC: Keywords/Channels management
  ├─ HTTP: WebSub receiver, Snapshot, Admin API
  └─ Neon(Postgres): channels/videos/video_snapshots

[ analytics-service (Go, Cloud Run, gRPC) ]
  ├─ gRPC: Rankings, history, video details
  └─ Neon(Postgres): video_metrics/ranking_snapshots

[ authority-service (Go, Cloud Run, gRPC) ]
  ├─ gRPC: Auth and profiles
  └─ Neon(Postgres): accounts/identities

[ Cloud Tasks ] → ingestion-service /snapshot
[ Cloud Scheduler ] → ingestion-service /admin/*
[ Identity Platform ] Email/Password authentication
```

## Microservices

### ingestion-service
- Responsibility: Data collection and storage
- Tech: Go, gRPC + HTTP, Clean Architecture
- Key features:
  - Receive YouTube WebSub
  - Collect trending videos
  - Acquire snapshots
  - Manage keywords

### analytics-service
- Responsibility: Data analysis and serving
- Tech: Go, gRPC, Clean Architecture
- Key features:
  - Metric computation
  - Ranking generation
  - History management
  - Theme extraction

### authority-service
- Responsibility: Authentication and authorization
- Tech: Go, gRPC, Identity Platform integration
- Key features:
  - ID token verification
  - Profile management
  - Role management

## Design Principles

- Clean Architecture + DDD: Business logic centered
- Microservices: Split by Bounded Context
- Internal gRPC: Type-safe internal APIs
- External HTTP: WebSub, Cloud Tasks and other external events
- Idempotency: Deterministic Task IDs + DB constraints

## Project Layout

```
youtube-analytics/
├── proto/                    # Protocol Buffers definitions
│   ├── analytics/
│   ├── ingestion/
│   └── authority/
├── services/
│   ├── go.work              # Go workspace
│   ├── pkg/
│   │   ├── identityauth/    # Shared auth package
│   │   └── pb/              # Generated Go code
│   ├── ingestion-service/
│   ├── analytics-service/
│   └── authority-service/
├── web/
│   └── client/              # Next.js frontend
└── db/
    └── migrations/          # Database migrations
```
