# Key Workflows (Technical Implementation Perspective)

## Competitor Tracking (New Videos → Snapshots)

1. WebSub (YouTube Hub) → POST /yt/websub (yt-svc)
2. Extract videoId → Save D0 with videos.list
3. Register +6h/+24h/+48h/+72h to Cloud Tasks
   - **TaskID**: Deterministically generated as `snap:{videoId}:{checkpoint}`
   - Same TaskIDs are automatically deduplicated
4. ETA → POST /snapshot (Fetch single video/UPSERT, delete remaining tasks if slowing down)

### Idempotency Guarantee Mechanism

**Deterministic TaskID Generation**
```go
taskID := fmt.Sprintf("snap:%s:%d", videoID, checkpointHour)
```

**Duplicate Prevention at DB Layer**
```sql
INSERT INTO video_snapshots (video_id, checkpoint_hour, ...)
VALUES ($1, $2, ...)
ON CONFLICT (video_id, checkpoint_hour) 
DO UPDATE SET updated_at = NOW()
RETURNING created_at = updated_at AS is_new;
```

- `is_new = true`: Newly inserted
- `is_new = false`: Existing record found (processing skipped due to idempotency)

## New Discovery (Trending → Subscription Promotion)

1. Scheduler (1-3 times/day) → POST /admin/collect-trending
2. videos.list?chart=mostPopular&regionCode=JP&videoCategoryId=27/28&maxResults=50 (1-2 pages)
3. Apply keyword OR filter to title/description/tags
4. For new channelIds, get subscriberCount with channels.list(part=statistics)
5. Promote to WebSub subscription if "subscriber ratio/growth rate" criteria are met

## Rankings & Themes

- From recent snapshots and subscriber counts: rel_views / momentum / like_rate → z-score → param_score
- From title/description n-gram/TF-IDF → frequent theme rankings (7-day/14-day)

## Operations Support (For Administrators)

### Manual Collection Execution
- When you want to "collect trending videos right now," you can execute immediately from the admin panel

### Threshold Adjustment
- You can edit conditions such as "don't include videos with less than 1000 views in rankings"
