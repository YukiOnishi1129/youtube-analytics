# Security & Observability

## Security Details

### Authentication & Authorization Strategy

**User API Authentication**
- Each service verifies OIDC (Identity Platform) independently
- Fetch and cache public keys from JWKS endpoint
- JWT verification items:
  - `iss` (issuer): Verify Identity Platform project ID
  - `aud` (audience): Verify target service
  - `exp` (expiration): Verify validity period
  - Signature verification: Verify with RS256 algorithm

**Internal API Authentication (Cloud Tasks/Scheduler)**
- Cloud Run IAM + OIDC dual defense
- Service account-based authentication:
  1. Cloud Run: Configure to reject unauthenticated calls
  2. Grant Cloud Run Invoker permission to calling service account
  3. Attach OIDC token to request header
  4. Application-side OIDC token verification (double check)

**Method-level Policy Control**
```
PUBLIC: No authentication required (health checks, etc.)
USER_ID_TOKEN: Identity Platform ID token required
SERVICE_OIDC: Service-to-service communication OIDC token required
```

### Common Authentication Implementation

Consolidated in `services/pkg/identityauth` package:
- JWKS client implementation
- Token verification logic
- gRPC interceptor
- HTTP middleware

All microservices import and use this package

### Other Security Measures

- **WebSub notifications**: HMAC-SHA256 signature verification
- **Secrets management**: Centralized management with Google Secret Manager
- **TLS**: Cloud Run automatic TLS termination + consider mTLS for internal communication

## Observability

### Metrics Collection
- **OpenTelemetry**: Integration of traces, metrics, and logs
- **Prometheus format**: Published at `/metrics` endpoint
- Key metrics:
  - Request rate, latency, error rate
  - YouTube API quota usage rate
  - Cloud Tasks queue length, failure rate
  - WebSub notification reception count

### Logging
- **Structured logging**: JSON format logs via Zap (Go)
- Log levels: DEBUG, INFO, WARN, ERROR, FATAL
- Required fields:
  - `trace_id`: For request tracking
  - `service`: Service name
  - `method`: gRPC/HTTP method name
  - `latency`: Processing time

### Tracing
- **Distributed tracing**: OpenTelemetry + Cloud Trace
- Span design:
  - Per gRPC method
  - External API calls (YouTube API, Cloud Tasks)
  - DB operations

### Alert Configuration
- **Error Reporting**: Automatic error aggregation
- **Logs-based Metrics**: 
  - Cloud Tasks failure rate > 5%
  - WebSub errors 3 consecutive times
  - YouTube API quota usage rate > 80%
- **Notification targets**: Slack webhook, PagerDuty (future)

## Authentication Flow (Next.js × Auth.js × Identity Platform)

- **Client**: Sign in with IP Web SDK → get idToken → pass to Auth.js(Credentials)
- **Server**: Firebase Admin verifyIdToken → session.user = { uid, email, role }
- **RSC/Server Actions**: Reference via getServerSession() / auth()
- **gRPC calls**: Attach Authorization: Bearer <idToken> to metadata
- **Password reset**: IP OOB email → verify→confirm at /reset?oobCode=
