# Service Specifications

## authority-service (Auth & Authorization)

Location: `/services/authority-service`  
Responsibility: Identity Platform integration, ID token verification, profile/role management

Database
- Shared Postgres instance with schema-per-service policy (authority schema)
- Access only own schema; no cross-schema joins
- Repositories generated via sqlc (see `services/authority-service/sqlc.yaml`) — build with `-tags sqlc`

Configuration (mandatory)
- Postgres: provide either `DATABASE_URL` or the individual envs `DB_HOST, DB_PORT(=5432), DB_USER, DB_PASSWORD, DB_NAME, DB_SSLMODE(=disable)`; `OpenPostgres` constructs DSN when `DATABASE_URL` is absent
- Identity Platform: `FIREBASE_API_KEY`
- OIDC verifier: `OIDC_ISSUER`, `OIDC_AUDIENCE`

Adapters (reference: authority-service)
- DB adapter: Postgres via pgx + sqlc (`internal/adapter/gateway/postgres`)
- Token verifier: OIDC using go-oidc (`internal/adapter/gateway/firebase/verifier.go`)
- gRPC auth: Unary interceptor injecting claims (`internal/driver/security/oidc_interceptor.go`)

### gRPC Methods (MVP)

| Method | Purpose | Auth | Request/Response |
|--------|---------|------|------------------|
| GetAccount | Get own profile | USER_ID_TOKEN | {id_token} → {account} |
| SignUp | Register with email/password | PUBLIC | {email,password} → {account,id_token,refresh_token} |
| SignIn | Login with email/password | PUBLIC | {email,password} → {id_token,refresh_token} |
| SignOut | Logout (revoke refresh token) | USER_ID_TOKEN | {refresh_token} → {} |
| ResetPassword | Send reset email | PUBLIC | {email} → {} |

## ingestion-service (Collection & Storage)

Location: `/services/ingestion-service`  
Responsibility: Video collection from YouTube, WebSub receiver, snapshot storage, keyword filter management

### Authentication Policies

- **PUBLIC**: Unauthenticated allowed (WebSub Verify/Notify for Hub→Server S2S)
- **USER_ID_TOKEN**: Identity Platform ID token (user operations)
- **SERVICE_OIDC**: Cloud Tasks/Scheduler OIDC (aud=service URL)
- Operation: Cloud Run disables unauthenticated calls, Tasks/Scheduler SA gets Invoker role
- WebSub is PUBLIC, internal APIs use SERVICE_OIDC middleware validation

### gRPC API (User/Admin UI)

Service: `ingestion.v1.IngestionService`

#### Keyword Management

| Method | Purpose | Auth | Request/Response |
|--------|---------|------|------------------|
| ListChannels | List monitored channels | User login required | query → [ChannelListItem] |
| SetChannelSubscription | Set channel subscription | User login required | {channel_id,subscribed} → void |
| ListKeywords | List keywords | User login required | → [Keyword] |
| CreateKeyword | Create keyword | User login required | {name,filter_type,description} → Keyword |
| UpdateKeyword | Update keyword | User login required | {id,name,filter_type,enabled,description} → Keyword |
| DeleteKeyword | Delete keyword | User login required | {id} → void |
| InsertSnapshot | Save snapshot (internal) | Service | {video_id,checkpoint_hour} → void |

### HTTP APIs (Hub → Server)

#### WebSub Endpoints

| Endpoint | Purpose | Auth | Request → Response |
|----------|---------|------|-------------------|
| GET /websub/youtube | Verify subscription | PUBLIC | Query: hub.mode, hub.topic, hub.challenge, hub.lease_seconds → Body: hub.challenge |
| POST /websub/youtube | Receive notification | PUBLIC + Secret | Body: Atom XML → 204 No Content |

**Processing:**
1. Extract videoId from Atom XML → `ApplyWebSubNotification(videoId)`
2. Get D0 snapshot → `ScheduleSnapshots(+3,+6,+12,+24,+48,+72,+168)`
3. Idempotency: D0 uses `UNIQUE(video_id, cp=0)` with `ON CONFLICT DO NOTHING`
4. Always return 2xx for Hub retry handling

### HTTP APIs (Cloud Tasks)

| Endpoint | Purpose | Auth | Request → Response |
|----------|---------|------|-------------------|
| POST /tasks/snapshot | Save checkpoint snapshot | SERVICE_OIDC | {video_id, checkpoint_hour} → 204 |

**Details:**
- Headers: `Authorization: Bearer <OIDC>` (Tasks-provided, aud=service URL)
- Processing: YouTube API → `video_snapshots.Insert` (idempotent)
- Idempotency: `UNIQUE(video_id, checkpoint_hour)` prevents duplicates
- TaskID: `snap:{video_id}:{checkpoint_hour}` (deterministic)
- Retry: 5xx only, exponential backoff

### HTTP APIs (Cloud Scheduler)

| Endpoint | Purpose | Auth | Request → Response |
|----------|---------|------|-------------------|
| POST /admin/collect-trending | Collect trending videos | SERVICE_OIDC | {region?, category_ids?, pages?} → {collected, adopted} |
| POST /admin/renew-subscriptions | Renew WebSub leases | SERVICE_OIDC | {} → {renewed} |
| GET /warm | Keep instance warm | SERVICE_OIDC | {} → "ok" |

**Collect Trending Details:**
- Default: `{region: "JP", category_ids: [27, 28], pages: 1}`
- Process: YouTube trending → Keyword filter → Register if matched
- Idempotency: `youtube_video_id UNIQUE` ignores duplicates

### Admin/Debug Endpoints (Optional)

| Endpoint | Purpose | Auth | Request → Response |
|----------|---------|------|-------------------|
| POST /admin/apply-websub-notification | Manual D0 retry | SERVICE_OIDC | {youtube_video_id} → 204 |

### Idempotency & Constraints

- **VideoSnapshot**: `UNIQUE(video_id, checkpoint_hour)` + `INSERT ... ON CONFLICT DO NOTHING`
- **Videos**: `youtube_video_id UNIQUE` (absorbs trending duplicates)
- **Channels**: `youtube_channel_id UNIQUE` (absorbs subscription duplicates)
- **Tasks**: TaskID = `snap:{video}:{cp}` (prevents duplicate enqueue)
- **Rate limit**: Consider token bucket on `/tasks/snapshot` via interceptor
- **Retry**: YouTube API handled by adapter with timeout/retry/backoff

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
