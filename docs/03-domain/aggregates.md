# Aggregates and Invariants

Aggregates define transactional consistency boundaries and the invariants they must uphold.

## Ingestion Context Aggregates

### A-1) FilterKeyword

#### Identifier
- KeywordId (UUID v7)

#### Attributes
- name: Display name (e.g., "React", "Next.js")
- filter_type: Filter type (include | exclude)
- pattern: Regular expression pattern
- enabled: Enabled flag
- description: Description

#### Invariants
- filter_type ∈ {include, exclude}
- Disabled keywords are not used for judgment
- pattern is auto-generated from name
- Deletion is soft-delete (set deleted_at)

#### Commands
```go
// Register or update a keyword
PutKeyword(name string, filterType FilterType, description string) error
// pattern is auto-generated as: pattern = BuildPattern(name)

// Enable/Disable
EnableKeyword(id KeywordId) error
DisableKeyword(id KeywordId) error

// Remove (soft delete)
RemoveKeyword(id KeywordId) error
```

#### Purpose
Apply consistent inclusion/exclusion rules during ingestion

### A-2) Channel

#### Identifier
- channel_id (YouTube Channel ID)

#### Attributes
- title: Channel name
- thumbnail_url: Thumbnail URL
- subscribed: WebSub subscription state
- subscription_expires_at: Subscription expiration

#### Child Entities
- ChannelSnapshot: Time series of subscriber counts
  - measured_at: Measurement datetime
  - subscription_count: Subscriber count

#### Invariants
- channel_id is immutable (from YouTube)
- subscribed is a single boolean flag
- Expired subscriptions are renewed automatically or manually

#### Commands
```go
// Subscription management
SubscribeChannel(channelId string) error
UnsubscribeChannel(channelId string) error
RenewSubscription(channelId string) error

// Record subscriber counts
RecordSubscriberCount(channelId string, measuredAt time.Time, count int64) error
```

#### Purpose
Maintain subscription continuity and subscriber trends (reference for video snapshot measurements)

### A-3) Video

#### Identifier
- video_id (YouTube Video ID)

#### Attributes
- channel_id: Owning channel
- title: Video title
- published_at: Publication datetime
- thumbnail_url: Thumbnail URL
- video_url: Video URL
- youtube_category_id: YouTube category

#### Child Entities
- VideoSnapshot: Snapshot at each checkpoint
  - checkpoint_hour: Checkpoint ∈ {0,3,6,12,24,48,72,168}
  - views_count: Views
  - likes_count: Likes
  - subscription_count: Channel subscriber count at that time (embedded)
  - measured_at: Measurement datetime
  - source: Data source ('websub'|'task'|'manual')

#### Invariants
- (video_id, checkpoint_hour) is insert-only (at most once)
- checkpoint_hour increases in ascending order
- subscription_count embeds the channel subscriber count at that time

#### Commands
```go
// New registration from trending
RegisterVideoFromTrending(meta VideoMeta) error

// WebSub notification handling (finalize D0 and schedule follow-ups)
ApplyWebSubNotification(videoId string) error

// Add snapshot
AddSnapshot(videoId string, checkpoint int, counts SnapshotCounts) error
```

#### Event
```go
SnapshotAdded {
    VideoId: string
    CheckpointHour: int
    MeasuredAt: time.Time
}
```

#### Purpose
Collect and accumulate data across the video's lifecycle

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

### C-1) Account

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
  - provider: Provider kind ('google'|'password'|'github')
  - provider_uid: Provider-side UID

#### Invariants
- email is unique
- (account_id, provider) is unique
- Deleting an account cascades to related identities
- Multiple providers can be linked

#### Queries
```go
// Fetch profile
GetMe() (*Account, []Identity, error)

// Create/Update account
CreateOrUpdateAccount(claims TokenClaims) error

// Link provider
LinkProvider(accountId string, provider string, providerUid string) error
```

#### Purpose
User and profile management with multi-provider support

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
