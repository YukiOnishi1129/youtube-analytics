# Clean Architecture + DDD Details

## Design Principles

### Clean Architecture Layers

```
┌──────────────────────────────────────┐
│           Frameworks & Drivers        │  Outer (concrete)
│  ┌──────────────────────────────┐    │
│  │     Interface Adapters        │    │
│  │  ┌──────────────────────┐    │    │
│  │  │  Application Business  │    │    │
│  │  │  ┌──────────────┐     │    │    │
│  │  │  │ Enterprise    │     │    │    │  Inner (abstract)
│  │  │  │ Business      │     │    │    │
│  │  │  └──────────────┘     │    │    │
│  │  └──────────────────────┘    │    │
│  └──────────────────────────────┐    │
└──────────────────────────────────────┘
```

### Dependency Rule

Dependency always points inward
- domain ← usecase ← adapter ← driver
- Outer layers may depend on inner layers
- Inner layers do not know outer layers

### Layer Responsibilities

#### Domain Layer (innermost)
- Entity: Core business concepts (Video, Channel, Keyword, Account)
- Value Object: Immutable values (VideoID, ChannelID, Checkpoint)
- Domain Service: Business logic spanning multiple entities
  - MetricCalculator: Growth rate, Wilson lower bound, LPS
  - RankingScorer: z-score normalization, param_score
- Note: Repository interfaces are not defined in Domain; they are owned by Use Case as Output Ports.

#### Use Case Layer
- Interactor/Application Service: Implement use cases
  - CollectTrendingVideos
  - CalculateVideoMetrics
  - GenerateRanking
- Ports (under `internal/port/`)
  - `input/`: Public interfaces of use cases
  - `output/gateway/`: DB/external API abstraction (e.g., VideoRepository, AccountRepository, TokenVerifier, Clock)
  - `output/presenter/`: Output boundary

#### Adapter Layer
- Controller: gRPC/HTTP handlers
- Presenter: Response shaping
- Gateway: Repository and external API implementations
  - PostgresVideoRepository
  - YouTubeAPIClient
  - CloudTasksClient

#### Driver Layer (outermost)
- Config: Configuration management
- Transport: gRPC/HTTP server bootstrap
- Datastore: DB connection management
- Security: AuthN/AuthZ middleware
- Observability: Metrics, logging, tracing
- Health: Health checks

## Mapping DDD Concepts

### Bounded Contexts

Each microservice forms a Bounded Context:

```
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│   Ingestion     │  │   Analytics     │  │   Authority     │
│    Context      │  │    Context      │  │    Context      │
│                 │  │                 │  │                 │
│  - Video        │  │  - Ranking      │  │  - Account      │
│  - Channel      │  │  - Metric       │  │  - Identity     │
│  - Keyword      │  │  - History      │  │  - Role         │
└─────────────────┘  └─────────────────┘  └─────────────────┘
        ↓                    ↓                    ↓
     gRPC Comm           gRPC Comm           gRPC Comm
```

### Key Aggregates

#### Video Aggregate
```go
type Video struct {
    ID          VideoID
    ChannelID   ChannelID  
    Title       string
    PublishedAt time.Time
    Snapshots   []VideoSnapshot  // entity within aggregate
}

type VideoSnapshot struct {
    VideoID         VideoID
    CheckpointHour  Checkpoint
    ViewsCount      int64
    LikesCount      int64
    MeasuredAt      time.Time
}
```

#### Channel Aggregate
```go
type Channel struct {
    ID              ChannelID
    Title           string
    Subscribed      bool
    SubscriberCount int64
}
```

#### Keyword Aggregate
```go
type Keyword struct {
    ID          KeywordID
    Name        string
    FilterType  FilterType
    Pattern     string
    Enabled     bool
}
```

### Domain Service Examples

```go
// domain/service/metric_calculator.go
type MetricCalculator interface {
    CalculateGrowthRate(baseline, current int64, hours int) float64
    CalculateWilsonLowerBound(likes, views int64) float64
    CalculateLPS(likes, subscribers int64) float64
}

// domain/service/ranking_scorer.go  
type RankingScorer interface {
    CalculateZScore(values []float64) []float64
    CalculateParamScore(momentum, relViews, quality float64) float64
}
```

## Implementation Patterns

### Repository Pattern
```go
// port/output/gateway/video_repository.go
type VideoRepository interface {
    Save(ctx context.Context, video *domain.Video) error
    FindByID(ctx context.Context, id domain.VideoID) (*domain.Video, error)
    FindByChannelID(ctx context.Context, channelID domain.ChannelID) ([]*domain.Video, error)
}

// adapter/gateway/postgres_video_repository.go
type PostgresVideoRepository struct {
    db *sql.DB
}

func (r *PostgresVideoRepository) Save(ctx context.Context, video *domain.Video) error {
    // ...
    return nil
}
```

// Authority (Account) specific output ports owned by UseCase
```go
// internal/port/output/gateway/repositories.go (excerpt)
type AccountRepository interface {
    FindByID(ctx context.Context, id string) (*domain.Account, error)
    FindByEmail(ctx context.Context, email string) (*domain.Account, error)
    Save(ctx context.Context, a *domain.Account) error
}

// internal/port/output/gateway/repositories.go (excerpt)
type IdentityRepository interface {
    ListByAccount(ctx context.Context, accountID string) ([]domain.Identity, error)
    FindByProvider(ctx context.Context, provider domain.Provider, providerUID string) (*domain.Account, error)
    Save(ctx context.Context, accountID string, id domain.Identity) error
    Delete(ctx context.Context, accountID string, provider domain.Provider) error
}

// internal/port/output/gateway/repositories.go (excerpt)
type RoleRepository interface {
    ListByAccount(ctx context.Context, accountID string) ([]domain.Role, error)
    Assign(ctx context.Context, accountID string, role domain.Role) error
    Revoke(ctx context.Context, accountID string, role domain.Role) error
}
```
