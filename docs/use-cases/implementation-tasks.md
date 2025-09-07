# Implementation Tasks (MVP Priority Order)

1. **yt-svc HTTP**: /yt/websub (GET/POST)・/snapshot・/admin/collect-trending・/admin/renew-subscriptions
2. **Cloud Tasks**: D0/6/24/48/72h ETA triggers → /snapshot (UPSERT・delete on slowdown)
3. **yt-svc gRPC**: ListTrending / ListThemes / SubscribeChannel / ListChannels / GetVideo
4. **auth-svc gRPC**: VerifyToken / GetMe / UpdateMe
5. **Next.js**: Auth.js(Credentials)×IP integration, gRPC client, UI (trending/themes/monitoring)
6. **Scores/Themes**: Derived views & index optimization
7. **Scheduler**: CRON (collect/renew/warm)
8. **Monitoring/Notifications**: Task failures, WebSub errors