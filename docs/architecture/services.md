# Service Specifications

## authority-service (Auth & Authorization)

Location: `/services/authority-service`  
Responsibility: Identity Platform integration, ID token verification, profile/role management

### gRPC Methods

| Method | Purpose | Auth | Request/Response |
|--------|---------|------|------------------|
| GetMe | Get own profile | User login required | → {id,email,display_name,identities} |

## ingestion-service (Collection & Storage)

Location: `/services/ingestion-service`  
Responsibility: Video collection, WebSub receiver, snapshot storage, keyword management

### gRPC Methods

| Method | Purpose | Auth | Request/Response |
|--------|---------|------|------------------|
| ListChannels | List monitored channels | User login required | query → [ChannelListItem] |
| SetChannelSubscription | Set channel subscription | User login required | {channel_id,subscribed} → void |
| ListKeywords | List keywords | User login required | → [Keyword] |
| CreateKeyword | Create keyword | User login required | {name,filter_type,description} → Keyword |
| UpdateKeyword | Update keyword | User login required | {id,name,filter_type,enabled,description} → Keyword |
| DeleteKeyword | Delete keyword | User login required | {id} → void |
| InsertSnapshot | Save snapshot (internal) | Service | {video_id,checkpoint_hour} → void |

### HTTP APIs (External Events)

| API | Method/Path | Caller | Purpose | Protection |
|-----|-------------|--------|---------|-----------|
| WebSub Verify | GET /yt/websub | YouTube Hub | Subscription verification (return hub.challenge) | Public (lightly protected) |
| WebSub Notify | POST /yt/websub | YouTube Hub | New video notification | Public (signature verification) |
| Snapshot | POST /snapshot | Cloud Tasks | Fetch snapshot at ETA → UPSERT | OIDC/HMAC |
| Collect Trending | POST /admin/collect-trending | Cloud Scheduler | Collect categories 27/28 | OIDC/HMAC |
| Renew Subscriptions | POST /admin/renew-subscriptions | Cloud Scheduler | Refresh WebSub leases | OIDC/HMAC |
| Warm | GET /warm | Cloud Scheduler | Cold start mitigation | Optional |

## analytics-service (Analysis & Serving)

Location: `/services/analytics-service`  
Responsibility: Metric precomputation, ranking serving, history management

### gRPC Methods

| Method | Purpose | Auth | Request/Response |
|--------|---------|------|------------------|
| ListRanking | Get rankings | User login required | RankingQuery → [RankingItem] |
| ListChannelRanking | Rankings per channel | User login required | ChannelRankingQuery → [RankingItem] |
| GetVideoDetail | Video detail and timeline | User login required | {video_id} → {video,snapshots,metrics} |
| ListHistory | Get history list | User login required | {from,to,ranking_kind?,checkpoint_hour?} → [History] |
| GetHistoryItems | Get history snapshot items | User login required | {snapshot_id} → [HistoryItem] |

