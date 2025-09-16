# Use Cases Overview

This directory contains detailed use case documentation for the YouTube Analytics system.

## Use Case Categories

### 1. [Data Collection](./data-collection.md)
Automated processes for collecting YouTube video data:
- Automatic Trending Video Collection (per genre)
- New Video Monitoring via WebSub
- Periodic Snapshot Collection
- Growth Rate and Metrics Calculation

### 2. [Video Analytics](./video-analytics.md)
User-facing features for analyzing video performance:
- Trending Video Ranking Display
- Video Detail View
- Channel-based Analysis
- Historical Data Export

### 3. [Admin Management](./admin-management.md)
Administrative functions for system configuration:
- YouTube Category Management
- Genre Management (regions × languages × categories)
- Keyword Management (per genre filtering)
- Batch Operations Control
- System Monitoring

## Key Flows

### Genre-based Collection Flow
1. Admin creates/enables genres with specific region, language, and categories
2. Admin configures keywords for each genre (include/exclude patterns)
3. Batch process runs for each enabled genre:
   - Fetches trending videos from YouTube API (region + categories)
   - Filters videos using genre-specific keywords
   - Registers matched videos with genre associations
4. Videos are monitored at checkpoints (0/3/6/12/24/48/72/168 hours)

### Multi-region Support Example
- **Engineering (JP)**: region=JP, categories=[27,28], keywords=Japanese tech terms
- **Engineering (EN)**: region=US, categories=[27,28], keywords=English tech terms  
- **English Learning (JP)**: region=JP, categories=[27], keywords=English education terms

### Admin Portal vs User Portal
- **Admin Portal**: Configure genres, keywords, categories, monitor system
- **User Portal**: View rankings, analyze videos, export data

## Access Control
- **Public**: WebSub endpoints
- **User**: Video analytics, rankings, exports
- **Admin**: All configuration and monitoring functions
- **Service**: Internal APIs (Cloud Tasks, Scheduler)