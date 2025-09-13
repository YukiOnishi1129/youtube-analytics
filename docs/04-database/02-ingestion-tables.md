# Ingestion Service Tables

## Overview

The ingestion service manages video collection, filtering, and monitoring with support for multiple genres (regions × languages × categories).

## Tables

### genres

Defines collection targets with region, category, and language settings.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID (v7) | PRIMARY KEY | Unique identifier |
| code | VARCHAR(50) | UNIQUE, NOT NULL | Genre code (e.g., "engineering_jp") |
| name | VARCHAR(100) | NOT NULL | Display name (e.g., "Engineering (JP)") |
| language | VARCHAR(10) | NOT NULL | Language code (e.g., "ja", "en") |
| region_code | VARCHAR(10) | NOT NULL | YouTube region (e.g., "JP", "US") |
| category_ids | INT[] | NOT NULL | YouTube category IDs (e.g., [27,28]) |
| enabled | BOOLEAN | NOT NULL DEFAULT true | Whether genre is active for collection |
| created_at | TIMESTAMP | NOT NULL DEFAULT NOW() | Creation timestamp |
| updated_at | TIMESTAMP | NOT NULL DEFAULT NOW() | Last update timestamp |

**Indexes:**
- `idx_genres_enabled` on (enabled)

### keywords

Filtering patterns for video selection within genres.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID (v7) | PRIMARY KEY | Unique identifier |
| genre_id | UUID | NOT NULL, FOREIGN KEY | Associated genre |
| name | VARCHAR(100) | NOT NULL | Keyword name |
| filter_type | VARCHAR(20) | NOT NULL CHECK IN ('INCLUDE','EXCLUDE') | Filter type |
| pattern | TEXT | NOT NULL | Regex pattern for matching |
| enabled | BOOLEAN | NOT NULL DEFAULT true | Whether keyword is active |
| description | TEXT | | Optional description |
| target_field | VARCHAR(20) | NOT NULL DEFAULT 'TITLE' | Field to match against |
| deleted_at | TIMESTAMP | | Soft delete timestamp |
| created_at | TIMESTAMP | NOT NULL DEFAULT NOW() | Creation timestamp |
| updated_at | TIMESTAMP | NOT NULL DEFAULT NOW() | Last update timestamp |

**Indexes:**
- `idx_keywords_genre_enabled` on (genre_id, enabled)
- `idx_keywords_deleted` on (deleted_at)

**Constraints:**
- UNIQUE on (genre_id, name, filter_type) WHERE deleted_at IS NULL

### youtube_categories

YouTube category master data managed by administrators.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INT | PRIMARY KEY | YouTube category ID |
| name | VARCHAR(100) | NOT NULL | Category name |
| assignable | BOOLEAN | NOT NULL | Whether videos can be assigned |
| created_at | TIMESTAMP | NOT NULL DEFAULT NOW() | Creation timestamp |
| updated_at | TIMESTAMP | NOT NULL DEFAULT NOW() | Last update timestamp |

### channels

YouTube channels being monitored.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID (v7) | PRIMARY KEY | Internal identifier |
| youtube_channel_id | VARCHAR(100) | UNIQUE, NOT NULL | YouTube channel ID |
| title | VARCHAR(255) | NOT NULL | Channel title |
| thumbnail_url | TEXT | | Channel thumbnail URL |
| description | TEXT | | Channel description |
| country | VARCHAR(10) | | Channel country |
| view_count | BIGINT | | Total channel views |
| subscription_count | BIGINT | | Current subscriber count |
| video_count | BIGINT | | Total videos |
| subscribed | BOOLEAN | NOT NULL DEFAULT false | WebSub subscription status |
| created_at | TIMESTAMP | NOT NULL DEFAULT NOW() | First seen timestamp |
| updated_at | TIMESTAMP | NOT NULL DEFAULT NOW() | Last update timestamp |

**Indexes:**
- `idx_channels_subscribed` on (subscribed)

### channel_snapshots

Historical channel metrics.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID (v7) | PRIMARY KEY | Snapshot identifier |
| channel_id | UUID | NOT NULL, FOREIGN KEY | Internal channel ID |
| measured_at | TIMESTAMP | NOT NULL | Measurement timestamp |
| subscription_count | BIGINT | NOT NULL | Subscriber count at time |
| view_count | BIGINT | | Total views at time |
| video_count | BIGINT | | Total videos at time |
| created_at | TIMESTAMP | NOT NULL DEFAULT NOW() | Creation timestamp |
| updated_at | TIMESTAMP | NOT NULL DEFAULT NOW() | Last update timestamp |

**Indexes:**
- `idx_channel_snapshots_channel_time` on (channel_id, measured_at DESC)

**Constraints:**
- UNIQUE on (channel_id, measured_at)

### videos

Monitored YouTube videos.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID (v7) | PRIMARY KEY | Internal identifier |
| youtube_video_id | VARCHAR(100) | UNIQUE, NOT NULL | YouTube video ID |
| channel_id | UUID | NOT NULL, FOREIGN KEY | Internal channel ID |
| title | TEXT | NOT NULL | Video title |
| published_at | TIMESTAMP | NOT NULL | Publication timestamp |
| category_id | INT | FOREIGN KEY | YouTube category ID |
| created_at | TIMESTAMP | NOT NULL DEFAULT NOW() | First seen timestamp |
| updated_at | TIMESTAMP | NOT NULL DEFAULT NOW() | Last update timestamp |

**Indexes:**
- `idx_videos_channel` on (channel_id)
- `idx_videos_published` on (published_at DESC)
- `idx_videos_category` on (category_id)

### video_genres

Junction table for video-genre many-to-many relationship.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID (v7) | PRIMARY KEY | Unique identifier |
| video_id | UUID | NOT NULL, FOREIGN KEY | Video ID |
| genre_id | UUID | NOT NULL, FOREIGN KEY | Genre ID |
| created_at | TIMESTAMP | NOT NULL DEFAULT NOW() | Association timestamp |

**Indexes:**
- `idx_video_genres_video` on (video_id)
- `idx_video_genres_genre` on (genre_id)

**Constraints:**
- UNIQUE on (video_id, genre_id)

### video_snapshots

Time-series video metrics at checkpoint hours.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID (v7) | PRIMARY KEY | Snapshot identifier |
| video_id | UUID | NOT NULL, FOREIGN KEY | Internal video ID |
| checkpoint_hour | INT | NOT NULL CHECK IN (0,3,6,12,24,48,72,168) | Hours after publication |
| measured_at | TIMESTAMP | NOT NULL | Actual measurement time |
| view_count | BIGINT | NOT NULL | View count |
| like_count | BIGINT | NOT NULL | Like count |
| subscription_count | BIGINT | NOT NULL | Channel subscribers at time |
| created_at | TIMESTAMP | NOT NULL DEFAULT NOW() | Creation timestamp |
| updated_at | TIMESTAMP | NOT NULL DEFAULT NOW() | Last update timestamp |

**Indexes:**
- `idx_video_snapshots_video_checkpoint` on (video_id, checkpoint_hour)
- `idx_video_snapshots_measured` on (measured_at DESC)

**Constraints:**
- UNIQUE on (video_id, checkpoint_hour)

## Data Flow

1. **Genre Configuration**: Admin enables genres with specific region/category/language settings
2. **Keyword Setup**: Admin adds inclusion/exclusion keywords for each genre
3. **Trending Collection**: Batch process fetches trending videos for enabled genres
4. **Filtering**: Videos are filtered by genre-specific keywords
5. **Registration**: Matched videos are registered with genre associations
6. **WebSub**: Channels are subscribed for real-time notifications
7. **Snapshots**: Video metrics collected at 0/3/6/12/24/48/72/168 hours

### audit_logs

Audit trail for administrative actions.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID (v7) | PRIMARY KEY | Log entry identifier |
| actor_id | UUID | NOT NULL | User who performed action |
| actor_email | VARCHAR(255) | NOT NULL | Actor email for reference |
| action | VARCHAR(100) | NOT NULL | Action performed |
| resource_type | VARCHAR(50) | NOT NULL | Type of resource affected |
| resource_id | VARCHAR(100) | | ID of affected resource |
| old_values | JSONB | | Previous values (for updates) |
| new_values | JSONB | | New values (for updates) |
| ip_address | INET | | Client IP address |
| user_agent | TEXT | | Client user agent |
| created_at | TIMESTAMP | NOT NULL DEFAULT NOW() | Action timestamp |

**Indexes:**
- `idx_audit_logs_actor` on (actor_id, created_at DESC)
- `idx_audit_logs_resource` on (resource_type, resource_id)
- `idx_audit_logs_created` on (created_at DESC)

### batch_jobs

Batch job execution history.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID (v7) | PRIMARY KEY | Job execution identifier |
| job_type | VARCHAR(50) | NOT NULL | Type of job (collect_trending, renew_subscriptions) |
| status | VARCHAR(20) | NOT NULL | Job status (pending, running, completed, failed) |
| parameters | JSONB | | Job parameters |
| started_at | TIMESTAMP | | Job start time |
| completed_at | TIMESTAMP | | Job completion time |
| error_message | TEXT | | Error details if failed |
| statistics | JSONB | | Job statistics (collected, adopted, etc.) |
| created_at | TIMESTAMP | NOT NULL DEFAULT NOW() | Job creation timestamp |

**Indexes:**
- `idx_batch_jobs_type_status` on (job_type, status)
- `idx_batch_jobs_created` on (created_at DESC)

## Migration from Single-Region Design

The new design supports multiple regions/languages by:
- Replacing hardcoded JP/categories with configurable genres
- Moving keywords under genres for targeted filtering
- Adding M:N video-genre relationships
- Keeping existing snapshot/channel structure intact
- Adding audit and job tracking tables