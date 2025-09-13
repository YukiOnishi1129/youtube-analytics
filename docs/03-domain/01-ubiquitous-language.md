# Ubiquitous Language

Common vocabulary definitions in Domain-Driven Design. A glossary used consistently across the team with the same meaning.

## Core Concepts

### Collection & Monitoring

- Popular Video Collection: Import popular videos from YouTube categories (Education=27 / Science & Technology=28)
- Monitoring Target: Videos whose subsequent trajectory will be tracked
- WebSub: Mechanism to receive real-time new video notifications from YouTube
- Subscribe: Register a channel as a monitoring target and automatically track new uploads

### Measurement & Recording

- Snapshot: Measured values at specific times after publication (0/3/6/12/24/48/72/168h)
- Checkpoint: Measurement timings (3h, 6h, 12h, 24h, 48h, 72h, 7d)
- D0: Initial data at the moment of publication (0-hour point)
- ETA: Estimated Time of Arrival (scheduled execution time)

### Analysis & Metrics

- Growth Rate: Increase from 0→X hours (per hour)
  - Speed: Per-hour increase of views and likes
- Relative Views: Views per subscriber (Views/Sub)
- Quality: Lower bound of like rate (Wilson lower bound)
- Heat: Shrunk likes per subscriber indicator (LPS)
- Momentum: Indicator representing short-term momentum

### Filtering

- Filter Keyword: Rules to decide inclusion/exclusion
  - Include Rule: Adopt videos containing this keyword
  - Exclude Rule: Exclude videos containing this keyword
- Pattern: Regular expression generated from the keyword

### Rankings

- Ranking: A list ordered by Published Range × Checkpoint × Metric
- RankingKind:
  - `speed_views`: Initial speed of view growth
  - `speed_likes`: Initial speed of like growth
  - `relative_views`: Relative views
  - `quality`: Quality (Wilson lower bound)
  - `heat`: Heat (LPS)
- Ranking History: Frozen Top-N at a point in time
- Low-sample: Videos excluded from rankings due to insufficient data

### Channels & Videos

- Channel: Content publisher on YouTube
- Video: Individual piece of content posted to YouTube
- PublishedAt: Datetime when the video was published on YouTube
- Subscriber Count: Number of channel subscribers
- Views Count: Number of video views
- Likes Count: Number of likes on the video

### Auth & Authorization

- Account: User account in this system
- Identity: Linkage information with auth providers (Google/password, etc.)
- Role: Permission level (admin/editor/viewer)

## Abbreviations & Terms

| Abbrev | Full Name | Description |
|------|---------|------|
| CP | Checkpoint | Measurement timing |
| LPS | Likes Per Subscription (Shrunk) | Shrunk like rate |
| TF-IDF | Term Frequency–Inverse Document Frequency | Text importance measure |
| OIDC | OpenID Connect | Authentication protocol |
| JWKS | JSON Web Key Set | Public key set |
| MVP | Minimum Viable Product | Smallest viable product |

## Business Rule Terms

- Rapid-growth video: A video whose views surge in a short period
- Competitor channel: Other YouTube channels operating in the same genre
- Theme: Topic representing video content (e.g., React, Docker, career change)
- Script template: Structure patterns of high-performing videos (future work)

## Algorithm Terms

- z-score normalization: Statistical method that transforms data to mean 0, std dev 1
- Wilson lower bound: Lower bound of a binomial proportion confidence interval (like-rate confidence)
- n-gram: Sequence of n consecutive characters or words
- Normalization of variant spellings: Treat "React", "react", and Japanese transliterations as equivalent
