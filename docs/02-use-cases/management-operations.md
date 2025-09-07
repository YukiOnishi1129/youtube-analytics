# Management & Operations Use Cases

## Manually Run Popular Collection

- Actor: Admin
- Process: Immediately collect popular videos in categories 27/28 and apply keyword filters
- Output: Report of counts for included/excluded

### Detailed Flow

1. Admin triggers execution
   - Click “Collect Now” on admin dashboard
   - Or call admin API endpoint directly

2. Fetch data from YouTube API
   - Fetch popular list for category 27 (Education)
   - Fetch popular list for category 28 (Science & Technology)
   - Up to 200 items per category

3. Filtering
   - Evaluate video titles with keyword filters
   - If exclusion filter matches: skip
   - If inclusion filter matches: include
   - If none matches: skip

4. Process included videos
   - Save video info into `videos` table
   - Save channel info into `channels` table (if not present)
   - Enable WebSub subscription for the channel
   - Save 0-hour snapshot for the video
   - Schedule future snapshot tasks

5. Result report
   - On completion, show:
     - Total videos fetched
     - Videos included
     - Videos excluded
     - New channels
     - Error count (if any)

### Timing

- Can be executed at any time
- Safe even if overlapping with the automated run (idempotent)
- Mind API rate limits (daily quota)

## Adjust Thresholds

- Actor: Admin
- Process: Edit rules like “exclude if views < 1000 at 24h”
- Output: Reflected in subsequent rankings

### Detailed Flow

1. Show current settings
   - Thresholds per checkpoint
   - Minimum views (`min_views_count`)
   - Minimum subscribers (`min_subscription_count`)
   - Minimum likes (`min_likes_count`) (optional)

2. Edit thresholds
   - Choose checkpoint: 3h/6h/12h/24h/48h/72h/168h
   - Enter new thresholds
   - Validation:
     - No negative values
     - Recommend higher values for later checkpoints

3. Impact simulation (optional)
   - Show impact if new thresholds are applied
   - Predicted excluded video count
   - Predicted count remaining in rankings

4. Save settings
   - Persist to env or config file
   - Record change history (who changed what and when)

5. Immediate application
   - Recompute existing `video_metrics_checkpoint`
   - Update `exclude_from_ranking`
   - Take effect on the next ranking request

### Default Thresholds (recommended)

| Checkpoint | Min Views | Min Subscribers |
|------------|-----------|-----------------|
| 3h         | 100       | 100             |
| 6h         | 300       | 100             |
| 12h        | 500       | 100             |
| 24h        | 1000      | 100             |
| 48h        | 2000      | 100             |
| 72h        | 3000      | 100             |
| 168h       | 5000      | 100             |

## Re-fetch Snapshots

- Actor: Admin
- Input: Video IDs, checkpoints
- Process: Re-fetch data for specified videos at specified times
- Output: Updated snapshots and metrics

### Detailed Flow

1. Select targets
   - Enter video IDs (multiple)
   - Select checkpoints (multiple)
   - Or choose “re-fetch all”

2. Fetch from YouTube API
   - Get current views and likes
   - Get current channel subscriber count

3. Update snapshots
   - Update `video_snapshots` (UPDATE or UPSERT)
   - Record updated timestamp

4. Recompute metrics
   - Run RecomputeMetrics
   - Recalculate growth, relative views, quality, heat
   - Update `exclude_from_ranking`

5. Confirm results
   - Show before/after values
   - Show change rates
   - Show errors if any

## Manage WebSub Subscriptions

- Actor: Admin
- Input: Channel ID, subscription state
- Process: Enable/disable/renew WebSub subscriptions
- Output: Updated subscription state

### Detailed Flow

1. Check current subscriptions
   - List subscribed channels
   - Show subscription expiry
   - Warn about near-expiry channels

2. Enable subscription
   - Specify channel ID
   - Send subscribe request to YouTube Hub
   - Register callback URL
   - Run verification

3. Renew subscription (extend lease)
   - Renew subscriptions nearing expiry
   - Typically automatic, but can be manual
   - Up to 10 days extension

4. Disable subscription
   - Stop monitoring the channel
   - Send unsubscribe request to YouTube Hub
   - Set `subscribed=false` in `channels`

5. Error handling
   - Retry on subscription failure
   - Show error responses from Hub
   - Log and alert
