# Schema Migration Guide

## Overview

This guide documents the schema changes from the original single-region design to the multi-genre support design.

## Key Changes

### 1. New Tables

#### genres
- Replaces hardcoded JP region and category logic
- Configurable per genre: region, language, categories
- Enables/disables collection targets dynamically

#### youtube_categories
- Reference table for YouTube category IDs
- Ensures referential integrity for video.category_id

#### video_genres
- Junction table for M:N relationship
- Videos can belong to multiple genres
- Example: A video about "Python programming in Japanese" could belong to both "Engineering (JP)" and "Programming Tutorial (JP)" genres

### 2. Modified Tables

#### keywords
- Added `genre_id` foreign key
- Keywords now belong to specific genres
- Allows different keyword sets per genre

#### videos
- Removed fields moved to domain logic or other services:
  - `youtube_channel_id` (redundant, use channel relation)
  - `description` (fetch from YouTube API when needed)
  - `tags[]` (not used in current logic)
  - `thumbnail_url` (fetch from YouTube API when needed)
  - `video_url` (can be constructed from youtube_video_id)
  - `duration` (fetch from YouTube API when needed)
  - `is_active` (not in domain model)
  - `crawled_at` (not in domain model)
- Removed direct `genre_id` (now uses video_genres junction)

#### video_snapshots
- Renamed for consistency:
  - `views_count` → `view_count`
  - `likes_count` → `like_count`
- Removed unused fields:
  - `dislike_count` (YouTube API deprecated)
  - `comment_count` (not used in metrics)
  - `channel_subscriber_count` (renamed to `subscription_count`)

### 3. Data Migration Steps

```sql
-- 1. Create new tables
CREATE TABLE genres (...);
CREATE TABLE youtube_categories (...);
CREATE TABLE video_genres (...);

-- 2. Seed youtube_categories from YouTube API
INSERT INTO youtube_categories (id, name, assignable) VALUES
(1, 'Film & Animation', true),
(2, 'Autos & Vehicles', true),
-- ... etc

-- 3. Create default genre for existing data
INSERT INTO genres (id, code, name, language, region_code, category_ids, enabled)
VALUES (
  gen_ulid(), 
  'engineering_jp',
  'Engineering (JP)',
  'ja',
  'JP',
  ARRAY[27, 28],
  true
);

-- 4. Add genre_id to keywords
ALTER TABLE keywords ADD COLUMN genre_id UUID;
UPDATE keywords SET genre_id = (SELECT id FROM genres WHERE code = 'engineering_jp');
ALTER TABLE keywords ALTER COLUMN genre_id SET NOT NULL;

-- 5. Migrate existing videos to video_genres
INSERT INTO video_genres (id, video_id, genre_id)
SELECT gen_ulid(), v.id, g.id
FROM videos v
CROSS JOIN genres g
WHERE g.code = 'engineering_jp';

-- 6. Clean up videos table
ALTER TABLE videos 
  DROP COLUMN youtube_channel_id,
  DROP COLUMN description,
  DROP COLUMN tags,
  DROP COLUMN thumbnail_url,
  DROP COLUMN video_url,
  DROP COLUMN duration,
  DROP COLUMN is_active,
  DROP COLUMN crawled_at;

-- 7. Rename video_snapshots columns
ALTER TABLE video_snapshots
  RENAME COLUMN views_count TO view_count,
  RENAME COLUMN likes_count TO like_count;

ALTER TABLE video_snapshots
  DROP COLUMN dislike_count,
  DROP COLUMN comment_count,
  RENAME COLUMN channel_subscriber_count TO subscription_count;
```

### 4. Application Code Changes

- Update repository interfaces to handle M:N video-genre relationships
- Modify trending collection to process by genre
- Update keyword filtering to be genre-specific
- Remove dependencies on dropped columns
- Update metric calculations for renamed columns

### 5. Benefits

1. **Multi-region Support**: Easy to add new regions/languages
2. **Flexible Categorization**: Videos can belong to multiple genres
3. **Dynamic Configuration**: Enable/disable genres without code changes
4. **Cleaner Schema**: Removed redundant/unused columns
5. **Better Separation**: Domain logic vs storage concerns