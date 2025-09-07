# Frontend Domain Specifications

## Entities (ViewModels)

```typescript
type Checkpoint = 3|6|12|24|48|72|168;
type RankingKind = 'speed_views'|'speed_likes'|'relative_views'|'quality'|'heat';

type VideoVM = {
  videoId: string; title: string; channelId: string;
  thumbnailUrl: string; videoUrl: string; publishedAt: string;
};

type SnapshotPointVM = { 
  checkpointHour: Checkpoint; 
  viewsCount: number; 
  likesCount: number; 
};

type MetricPointVM = { 
  checkpointHour: Checkpoint; 
  viewGrowthRatePerHour: number; 
  likeGrowthRatePerHour: number; 
};

type RankingQueryVM = {
  publishedFrom: string; publishedTo: string;        // ISO(UTC)
  checkpointHour: Checkpoint; rankingKind: RankingKind;
  hideLowSample: boolean; category?: number;         // YouTube category id
  limit?: number; offset?: number;
};

type RankingItemVM = VideoVM & {
  checkpointHour: Checkpoint; mainMetric: number;
  viewsCount: number; likesCount: number;
};

type ChannelListItemVM = {
  channelId: string; title: string; thumbnailUrl?: string;
  subscribed: boolean; subscriberCount?: number; latestVideoPublishedAt?: string;
};

type ChannelListQueryVM = { 
  subscribed?: 'on'|'off'|'all'; 
  sort?: 'latest_video'|'subs'|'title'; 
  q?: string; 
  limit?: number; 
  offset?: number; 
};

type ChannelRankingQueryVM = RankingQueryVM & { channelId: string; };

type HistoryVM = { 
  snapshotId: string; snapshotAt: string; rankingKind: RankingKind; 
  checkpointHour: Checkpoint; publishedFrom: string; publishedTo: string; topN: number; 
};

type HistoryItemVM = { rank: number } & RankingItemVM;

type FilterType = 'include'|'exclude';
type KeywordVM = { 
  id: string; name: string; filterType: FilterType; 
  enabled: boolean; description?: string; 
};

type AccountVM = { 
  id: string; email: string; emailVerified: boolean; 
  displayName?: string; photoUrl?: string; 
  identities: { provider: 'google'|'password'|'github'; providerUid?: string; }[] 
};
```

## Service I/O (gRPC Contract Image)

**analytics-service**
- `ListRanking(RankingQueryVM) -> RankingItemVM[]`
- `ListChannelRanking(ChannelRankingQueryVM) -> RankingItemVM[]`
- `GetVideoDetail({videoId}) -> { video: VideoVM, snapshots: SnapshotPointVM[], metrics: MetricPointVM[] }`
- `ListHistory({from,to,rankingKind?,checkpointHour?}) -> HistoryVM[]`
- `GetHistoryItems({snapshotId}) -> HistoryItemVM[]`
- `ExportHistoryCSV({snapshotId}) -> file`

**ingestion-service**
- `ListChannels(ChannelListQueryVM) -> ChannelListItemVM[]`
- `SetChannelSubscription({channelId, subscribed}) -> void`
- `ListKeywords() -> KeywordVM[]`
- `CreateKeyword({name, filterType, description?}) -> KeywordVM`
- `UpdateKeyword({id, name?, filterType?, enabled?, description?}) -> KeywordVM`
- `DeleteKeyword({id}) -> void`

**authority-service**
- `GetMe() -> AccountVM`

## Screen-to-Service Call Mapping

- `/ranking` → analytics.ListRanking (Pass RankingQuery directly)
- `/videos/:id` → analytics.GetVideoDetail
- `/history` → analytics.ListHistory
- `/history/:id` → analytics.GetHistoryItems / analytics.ExportHistoryCSV
- `/channels` → ingestion.ListChannels / ingestion.SetChannelSubscription
- `/keywords` → ingestion.ListKeywords / Create|Update|DeleteKeyword
- `/mypage` → authority.GetMe

## Exception & State Handling (Common)

- **Validation**:
  - Published Range: from < to, convert to UTC
  - CP: in {3,6,12,24,48,72,168}
  - RankingKind: Fixed 5 types
- **Retry**: API 429/5xx with exponential backoff + toast display
- **Low-sample**: Frontend passes hideLowSample as-is. Exclusion logic is server-side