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
- Entity: Core business concepts (Video, Channel, Keyword)
- Value Object: Immutable values (VideoID, ChannelID, Checkpoint)
- Domain Service: Business logic spanning multiple entities
  - MetricCalculator: Growth rate, Wilson lower bound, LPS
  - RankingScorer: z-score normalization, param_score
- Repository Interface: Persistence abstraction (interfaces only)

#### Use Case Layer
- Interactor/Application Service: Implement use cases
  - CollectTrendingVideos
  - CalculateVideoMetrics
  - GenerateRanking
- Input Port: Public interface of use cases
- Output Port: Interfaces to external systems
  - Gateway: DB/external API abstraction
  - Presenter: Output boundary

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
