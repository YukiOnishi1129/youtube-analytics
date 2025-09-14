-- Down migration to revert genre-based design changes

-- Drop new tables
DROP TABLE IF EXISTS ingestion.batch_jobs;
DROP TABLE IF EXISTS ingestion.audit_logs;
DROP TABLE IF EXISTS ingestion.video_genres;
DROP TABLE IF EXISTS ingestion.youtube_categories;
DROP TABLE IF EXISTS ingestion.genres;

-- Revert video_snapshots column names
ALTER TABLE ingestion.video_snapshots
  RENAME COLUMN view_count TO views_count;
ALTER TABLE ingestion.video_snapshots
  RENAME COLUMN like_count TO likes_count;

-- Remove updated_at from video_snapshots
ALTER TABLE ingestion.video_snapshots
  DROP COLUMN IF EXISTS updated_at;

-- Revert keywords table changes
ALTER TABLE ingestion.keywords
  DROP CONSTRAINT IF EXISTS keywords_genre_id_fkey;
ALTER TABLE ingestion.keywords
  DROP COLUMN IF EXISTS genre_id,
  DROP COLUMN IF EXISTS target_field;

-- Restore original keywords unique index
DROP INDEX IF EXISTS ingestion.keywords_genre_name_filter_unique;
DROP INDEX IF EXISTS ingestion.idx_keywords_genre_enabled;
DROP INDEX IF EXISTS ingestion.idx_keywords_deleted;
CREATE UNIQUE INDEX keywords_name_unique ON ingestion.keywords(name) WHERE deleted_at IS NULL;

-- Restore videos table columns
ALTER TABLE ingestion.videos
  ADD COLUMN IF NOT EXISTS thumbnail_url text NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS video_url text NOT NULL DEFAULT '';

-- Remove constraint from channel_snapshots
ALTER TABLE ingestion.channel_snapshots
  DROP CONSTRAINT IF EXISTS channel_snapshots_channel_measured_unique;

-- Revert channel_snapshots table changes
ALTER TABLE ingestion.channel_snapshots
  DROP COLUMN IF EXISTS view_count,
  DROP COLUMN IF EXISTS video_count,
  DROP COLUMN IF EXISTS updated_at;

-- Revert channels table changes
ALTER TABLE ingestion.channels
  DROP COLUMN IF EXISTS description,
  DROP COLUMN IF EXISTS country,
  DROP COLUMN IF EXISTS view_count,
  DROP COLUMN IF EXISTS subscription_count,
  DROP COLUMN IF EXISTS video_count;

-- Drop new indexes
DROP INDEX IF EXISTS ingestion.idx_genres_enabled;