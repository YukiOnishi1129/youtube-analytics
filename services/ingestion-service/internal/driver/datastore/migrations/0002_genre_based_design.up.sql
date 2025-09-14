-- Migration to support genre-based video collection

-- Create genres table
CREATE TABLE IF NOT EXISTS ingestion.genres (
  id           uuid PRIMARY KEY,
  code         varchar(50) UNIQUE NOT NULL,
  name         varchar(100) NOT NULL,
  language     varchar(10) NOT NULL,
  region_code  varchar(10) NOT NULL,
  category_ids integer[] NOT NULL,
  enabled      boolean NOT NULL DEFAULT true,
  created_at   timestamptz NOT NULL DEFAULT now(),
  updated_at   timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX idx_genres_enabled ON ingestion.genres(enabled);

-- Create youtube_categories table
CREATE TABLE IF NOT EXISTS ingestion.youtube_categories (
  id         integer PRIMARY KEY,
  name       varchar(100) NOT NULL,
  assignable boolean NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

-- Update channels table with new fields
ALTER TABLE ingestion.channels
  ADD COLUMN IF NOT EXISTS description text,
  ADD COLUMN IF NOT EXISTS country varchar(10),
  ADD COLUMN IF NOT EXISTS view_count bigint,
  ADD COLUMN IF NOT EXISTS subscription_count bigint,
  ADD COLUMN IF NOT EXISTS video_count bigint;

-- Update channel_snapshots table
ALTER TABLE ingestion.channel_snapshots
  ADD COLUMN IF NOT EXISTS view_count bigint,
  ADD COLUMN IF NOT EXISTS video_count bigint,
  ADD COLUMN IF NOT EXISTS updated_at timestamptz NOT NULL DEFAULT now();

-- Add unique constraint to channel_snapshots
ALTER TABLE ingestion.channel_snapshots 
  ADD CONSTRAINT channel_snapshots_channel_measured_unique UNIQUE (channel_id, measured_at);

-- Remove thumbnail_url and video_url from videos table (not in new schema)
ALTER TABLE ingestion.videos
  DROP COLUMN IF EXISTS thumbnail_url,
  DROP COLUMN IF EXISTS video_url;

-- Update keywords table to add genre_id and target_field
ALTER TABLE ingestion.keywords
  ADD COLUMN IF NOT EXISTS genre_id uuid,
  ADD COLUMN IF NOT EXISTS target_field varchar(20) NOT NULL DEFAULT 'TITLE';

-- Add foreign key constraint after adding column
ALTER TABLE ingestion.keywords
  ADD CONSTRAINT keywords_genre_id_fkey FOREIGN KEY (genre_id) REFERENCES ingestion.genres(id) ON DELETE CASCADE;

-- Update unique constraint on keywords to include genre_id
DROP INDEX IF EXISTS ingestion.keywords_name_unique;
CREATE UNIQUE INDEX keywords_genre_name_filter_unique ON ingestion.keywords(genre_id, name, filter_type) WHERE deleted_at IS NULL;

-- Create new index for keywords
CREATE INDEX idx_keywords_genre_enabled ON ingestion.keywords(genre_id, enabled);
CREATE INDEX idx_keywords_deleted ON ingestion.keywords(deleted_at);

-- Create video_genres junction table
CREATE TABLE IF NOT EXISTS ingestion.video_genres (
  id         uuid PRIMARY KEY,
  video_id   uuid NOT NULL REFERENCES ingestion.videos(id) ON DELETE CASCADE,
  genre_id   uuid NOT NULL REFERENCES ingestion.genres(id) ON DELETE CASCADE,
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX idx_video_genres_video ON ingestion.video_genres(video_id);
CREATE INDEX idx_video_genres_genre ON ingestion.video_genres(genre_id);
CREATE UNIQUE INDEX video_genres_video_genre_unique ON ingestion.video_genres(video_id, genre_id);

-- Create audit_logs table
CREATE TABLE IF NOT EXISTS ingestion.audit_logs (
  id            uuid PRIMARY KEY,
  actor_id      uuid NOT NULL,
  actor_email   varchar(255) NOT NULL,
  action        varchar(100) NOT NULL,
  resource_type varchar(50) NOT NULL,
  resource_id   varchar(100),
  old_values    jsonb,
  new_values    jsonb,
  ip_address    inet,
  user_agent    text,
  created_at    timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX idx_audit_logs_actor ON ingestion.audit_logs(actor_id, created_at DESC);
CREATE INDEX idx_audit_logs_resource ON ingestion.audit_logs(resource_type, resource_id);
CREATE INDEX idx_audit_logs_created ON ingestion.audit_logs(created_at DESC);

-- Create batch_jobs table
CREATE TABLE IF NOT EXISTS ingestion.batch_jobs (
  id            uuid PRIMARY KEY,
  job_type      varchar(50) NOT NULL,
  status        varchar(20) NOT NULL,
  parameters    jsonb,
  started_at    timestamptz,
  completed_at  timestamptz,
  error_message text,
  statistics    jsonb,
  created_at    timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX idx_batch_jobs_type_status ON ingestion.batch_jobs(job_type, status);
CREATE INDEX idx_batch_jobs_created ON ingestion.batch_jobs(created_at DESC);

-- Update video_snapshots to fix column names
ALTER TABLE ingestion.video_snapshots
  RENAME COLUMN views_count TO view_count;
ALTER TABLE ingestion.video_snapshots
  RENAME COLUMN likes_count TO like_count;

-- Add updated_at column to video_snapshots
ALTER TABLE ingestion.video_snapshots
  ADD COLUMN IF NOT EXISTS updated_at timestamptz NOT NULL DEFAULT now();