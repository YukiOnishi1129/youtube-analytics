# Schedule & Queue

## Cloud Scheduler Configuration Example

```
*/5 * * * * → GET /warm (Next, auth-svc, yt-svc each 1 instance)
0 3,15 * * * → yt-svc POST /admin/collect-trending (27/28, pages=1..2)
0 5 * * * → yt-svc POST /admin/renew-subscriptions (WebSub resubscription)
```

## Cloud Tasks Configuration

- **Queue**: yt-snapshots
- **payload**: { videoId, checkpoint }
- **Task name**: snap:{videoId}:{cp}
- **ETA** calls /snapshot
- Retry/max attempts & DLQ configuration