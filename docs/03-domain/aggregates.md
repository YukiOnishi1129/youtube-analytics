# Aggregates and Invariants

Aggregates define transactional consistency boundaries and the invariants they must uphold.

## Ingestion Context

The Ingestion Service is responsible for:
- Ingesting data from YouTube
- Registering new videos as monitoring targets  
- Saving D0/3h/6h/12h/24h/48h/72h/168h snapshots
- Filtering videos by keywords
- Managing channel subscriptions

All internal entities use UUID v7 as primary keys, while YouTube external IDs are stored in separate fields.

### Common Value Objects

```go
type UUID = string                    // v7 assumed
type YouTubeVideoID = string
type YouTubeChannelID = string  
type CategoryID = int                 // YouTube videos.snippet.categoryId as int
type CheckpointHour = int             // 0,3,6,12,24,48,72,168
```

### Ingestion Context Aggregates

### A-1) YouTubeCategory

#### Purpose
Manage YouTube category reference data with admin control

#### Entity
```go
type YouTubeCategory struct {
    ID          CategoryID  // YouTube's category ID (e.g., 27, 28)
    Name        string      // Category name (e.g., "Education", "Science & Technology")
    Assignable  bool        // Whether videos can be assigned to this category
}
```

#### Invariants
- ID must be a valid YouTube category ID
- Name must not be empty
- Only assignable categories can be used for videos

#### Commands
```go
// Update category information
UpdateCategory(name string, assignable bool) error
```

### A-2) Genre

#### Purpose
Define collection targets with region, categories, and language settings

#### Entity
```go
type Genre struct {
    ID           UUID
    Code         string         // e.g., "engineering_jp", "engineering_en"
    Name         string         // e.g., "Engineering (JP)"
    Language     string         // e.g., "ja", "en"
    RegionCode   string         // e.g., "JP", "US"
    CategoryIDs  []CategoryID   // e.g., [27, 28] for Education & Science
    Enabled      bool
}
```

#### Invariants
- Code must be unique
- RegionCode must be valid YouTube region
- CategoryIDs must contain at least one valid category
- Only enabled genres are processed in batch collection

#### Commands
```go
// Enable/Disable genre
EnableGenre(id UUID) error
DisableGenre(id UUID) error
```

### A-3) Keyword

#### Purpose
Define inclusion/exclusion rules for video selection within a genre

#### Entity
```go
type FilterType string
const (
    Include FilterType = "include"
    Exclude FilterType = "exclude"
)

type Keyword struct {
    ID          UUID
    GenreID     UUID           // Associated genre
    Name        string
    FilterType  FilterType     // include / exclude
    Pattern     string         // Normalized regex pattern
    Enabled     bool
    Description *string
}
```

#### Domain Service
```go
// Pure function to build pattern from name
func BuildPattern(name string) (string, error)
```

#### Invariants
- Pattern != "" (must not be empty)
- FilterType ∈ {include, exclude}
- (GenreID, Name, FilterType) logical uniqueness (enforced at app/repo level)
- Disabled keywords are not used for judgment
- Pattern is auto-generated from Name
- Keywords belong to exactly one genre

#### Commands
```go
// Register or update a keyword for a genre
PutKeyword(genreID UUID, name string, filterType FilterType, description string) error

// Enable/Disable
EnableKeyword(id UUID) error
DisableKeyword(id UUID) error

// Remove (soft delete)
RemoveKeyword(id UUID) error
```

### A-4) Channel

#### Purpose
Manage subscription state and track subscriber trends. External ID is stored separately.

#### Entity
```go
type Channel struct {
    ID               UUID              // Internal PK
    YouTubeChannelID YouTubeChannelID  // External ID
    Title            string
    ThumbnailURL     string
    Subscribed       bool              // WebSub subscription target
}
```

#### Child Entity
```go
type ChannelSnapshot struct {
    ID                UUID              // Snapshot UUID (for aggregation/audit)
    ChannelID         UUID              // Internal FK
    MeasuredAt        time.Time
    SubscriptionCount int
}
```

#### Invariants
- YouTubeChannelID must not be empty
- ChannelSnapshot: (ChannelID, MeasuredAt) logical uniqueness (insert-only)
- Only Subscribed=true channels are WebSub renewal targets
- Channel ID is immutable once created

#### Commands
```go
// Subscription management
SubscribeChannel(channelId UUID) error
UnsubscribeChannel(channelId UUID) error
RenewSubscription(channelId UUID) error

// Record subscriber counts
RecordSubscriberCount(channelId UUID, measuredAt time.Time, count int64) error
```

### A-5) Video

#### Purpose
Track monitored videos with metadata and checkpoint snapshots. External ID is stored separately.

#### Entity
```go
type Video struct {
    ID               UUID              // Internal PK
    YouTubeVideoID   YouTubeVideoID    // External ID
    ChannelID        UUID              // Internal FK
    Title            string
    PublishedAt      time.Time
    CategoryID       CategoryID        // YouTube categoryId (e.g. 27,28)
}
```

#### Child Entity
```go
type VideoSnapshot struct {
    ID                UUID             // Snapshot UUID
    VideoID           UUID             // Internal FK
    CheckpointHour    CheckpointHour   // 0,3,6,12,24,48,72,168
    MeasuredAt        time.Time
    ViewCount         int64
    LikeCount         int64
    SubscriptionCount int64            // Channel subscriber count at time (copy)
}
```

#### Domain Service
```go
// Schedule snapshots for +3/6/12/24/48/72/168h via Cloud Tasks
ScheduleSnapshots(videoID UUID) error
```

#### Invariants
- MeasuredAt >= Video.PublishedAt
- VideoSnapshot: (VideoID, CheckpointHour) logical uniqueness (insert-only)
- Snapshots are immutable (no updates, only inserts)
- CheckpointHour increases in ascending order
- Videos can be associated with multiple genres (M:N relationship)

#### Commands
```go
// New registration from trending (with genre associations)
RegisterVideoFromTrending(genreIDs []UUID, meta VideoMeta) error

// WebSub notification handling (finalize D0 and schedule follow-ups)
ApplyWebSubNotification(ytVideoID YouTubeVideoID) error

// Add snapshot (idempotent)
AddSnapshot(videoID UUID, checkpoint CheckpointHour, counts SnapshotCounts) error
```

#### Event
```go
SnapshotAdded {
    VideoId: UUID
    CheckpointHour: CheckpointHour
    MeasuredAt: time.Time
}
```

## Analytics Context Aggregates

### B-1) VideoMetrics

#### Identifier
- (video_id, checkpoint_hour)

#### Attributes (computed)
```go
type VideoMetrics struct {
    // Basics
    VideoId         string
    CheckpointHour  int
    PublishedAt     time.Time  // kept redundantly for ranking queries
    
    // Raw at X
    ViewsCount         int64
    LikesCount         int64
    SubscriptionCount  int64
    
    // Baseline at 0
    ViewsBaseline      int64
    LikesBaseline      int64
    SubscriptionBaseline int64
    
    // Growth 0→X
    ViewGrowthRatePerHour float64
    LikeGrowthRatePerHour float64
    
    // Point-in-time indicators
    ViewsPerSubscriptionRate       float64  // relative views
    WilsonLikeRateLowerBound       float64  // quality
    LikesPerSubscriptionShrunkRate float64  // heat
    
    // Flags
    ExcludeFromRanking bool  // low-sample judgment
}
```

#### Invariants
- Recompute only when Snapshot(0) and Snapshot(X) are both available
- Computed metrics are overwriteable (idempotent)
- ExcludeFromRanking is determined by threshold rules

#### Commands
```go
// Recompute metrics
RecomputeMetrics(videoId string, checkpointHour int) error
```

#### Purpose
Maintain a read model ready for rankings with consistent rules

### B-2) RankingSnapshot

#### Identifier
- snapshot_id (UUID v7)

#### Attributes
- snapshot_at: Saved at
- ranking_kind: Ranking kind
- checkpoint_hour: Checkpoint
- published_from: Published date range start
- published_to: Published date range end
- top_n: Number of items saved

#### Child Entities
- RankingSnapshotItem: Video info per rank
  - rank: Rank
  - video_id: Video ID
  - title: Title
  - published_at: Published at
  - main_metric: Primary metric value
  - views_count: Views
  - likes_count: Likes

#### Invariants
- A saved snapshot is immutable
- rank is contiguous from 1 to top_n

#### Commands
```go
// Save daily Top-N
CreateDailyTopN(
    rankingKind RankingKind,
    checkpointHour int,
    publishedRange DateRange,
    topN int,
) error
```

#### Purpose
Freeze and persist the ranking at that time (for review and CSV export)

## Authority Context Aggregates

### C-1) Account (Aggregate Root)

#### Identifier
- account_id (UUID v7)

#### Attributes
- email: Email address
- email_verified: Email verification state
- display_name: Display name
- photo_url: Profile image URL
- is_active: Account activation state
- last_login_at: Last login datetime

#### Child Entities
- Identity: Authentication provider info
  - provider: Provider kind ('password'|'google'|'github')
  - provider_uid: Provider-side UID
- Role: Permission set
  - name: 'admin' | 'user' | ... (system-defined)

#### Invariants
- email is unique among active accounts
- (account_id, provider) is unique
- At least one identity must exist per account
- Deactivated accounts cannot sign in
- Roles must be within system-defined set; duplicate assignment is not allowed

#### Commands (Domain Methods)
```go
// Lifecycle
Deactivate()
Reactivate()
VerifyEmail()
TouchLogin(now time.Time)

// Profile
UpdateProfile(displayName string, photoURL string)

// Identity management
LinkIdentity(provider Provider, providerUID string) error
UnlinkIdentity(provider Provider) error // cannot unlink last identity

// Role management
AssignRole(role Role) error
RevokeRole(role Role) error
```

#### Purpose
User management with multi-provider login and role assignment

## Relationships Between Aggregates

### Event Flow
```
Video.AddSnapshot()
  → SnapshotAdded Event
  → VideoMetrics.RecomputeMetrics()
```

### Data References
```
VideoMetrics ← VideoSnapshot (read-only)
VideoSnapshot ← ChannelSnapshot (subscriber count embedded)
RankingSnapshot ← VideoMetrics (at snapshot creation)
```

### Idempotency Guarantees
- VideoSnapshot: DB constraint on (video_id, checkpoint_hour)
- VideoMetrics: Recompute is idempotent (same input → same output)
- RankingSnapshot: snapshot_id uniqueness via UUID v7
