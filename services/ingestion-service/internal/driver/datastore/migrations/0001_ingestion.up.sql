-- Up migration: create ingestion schema and tables
CREATE SCHEMA IF NOT EXISTS ingestion;

-- Channels table
CREATE TABLE IF NOT EXISTS ingestion.channels (
  id                  uuid PRIMARY KEY,
  youtube_channel_id  text NOT NULL,
  title               text NOT NULL,
  thumbnail_url       text NOT NULL,
  subscribed          boolean DEFAULT false,
  created_at          timestamptz DEFAULT now(),
  updated_at          timestamptz,
  deleted_at          timestamptz
);
CREATE UNIQUE INDEX IF NOT EXISTS channels_youtube_channel_id_unique ON ingestion.channels(youtube_channel_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS channels_subscribed_idx ON ingestion.channels(subscribed) WHERE deleted_at IS NULL;

-- Videos table
CREATE TABLE IF NOT EXISTS ingestion.videos (
  id                  uuid PRIMARY KEY,
  youtube_video_id    text NOT NULL,
  channel_id          uuid NOT NULL REFERENCES ingestion.channels(id) ON DELETE CASCADE,
  youtube_channel_id  text NOT NULL,
  title               text NOT NULL,
  published_at        timestamptz NOT NULL,
  category_id         integer NOT NULL,
  thumbnail_url       text NOT NULL,
  video_url           text NOT NULL,
  created_at          timestamptz DEFAULT now(),
  updated_at          timestamptz,
  deleted_at          timestamptz
);
CREATE UNIQUE INDEX IF NOT EXISTS videos_youtube_video_id_unique ON ingestion.videos(youtube_video_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS videos_channel_id_idx ON ingestion.videos(channel_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS videos_published_at_idx ON ingestion.videos(published_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS videos_category_id_idx ON ingestion.videos(category_id) WHERE deleted_at IS NULL;

-- Keywords table for filtering
CREATE TABLE IF NOT EXISTS ingestion.keywords (
  id           uuid PRIMARY KEY,
  name         text NOT NULL,
  filter_type  text NOT NULL CHECK (filter_type IN ('include', 'exclude')),
  pattern      text NOT NULL,
  enabled      boolean DEFAULT true,
  description  text,
  created_at   timestamptz DEFAULT now(),
  updated_at   timestamptz,
  deleted_at   timestamptz
);
CREATE UNIQUE INDEX IF NOT EXISTS keywords_name_unique ON ingestion.keywords(name) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS keywords_filter_type_idx ON ingestion.keywords(filter_type) WHERE deleted_at IS NULL AND enabled = true;
CREATE INDEX IF NOT EXISTS keywords_enabled_idx ON ingestion.keywords(enabled) WHERE deleted_at IS NULL;

-- Channel snapshots table
CREATE TABLE IF NOT EXISTS ingestion.channel_snapshots (
  id                  uuid PRIMARY KEY,
  channel_id          uuid NOT NULL REFERENCES ingestion.channels(id) ON DELETE CASCADE,
  measured_at         timestamptz NOT NULL,
  subscription_count  integer NOT NULL,
  created_at          timestamptz DEFAULT now()
);
CREATE INDEX IF NOT EXISTS channel_snapshots_channel_id_idx ON ingestion.channel_snapshots(channel_id);
CREATE INDEX IF NOT EXISTS channel_snapshots_measured_at_idx ON ingestion.channel_snapshots(measured_at DESC);

-- Video snapshots table
CREATE TABLE IF NOT EXISTS ingestion.video_snapshots (
  id                  uuid PRIMARY KEY,
  video_id            uuid NOT NULL REFERENCES ingestion.videos(id) ON DELETE CASCADE,
  checkpoint_hour     integer NOT NULL CHECK (checkpoint_hour IN (0, 3, 6, 12, 24, 48, 72, 168)),
  measured_at         timestamptz NOT NULL,
  views_count         bigint NOT NULL DEFAULT 0,
  likes_count         bigint NOT NULL DEFAULT 0,
  subscription_count  bigint NOT NULL DEFAULT 0,
  source              text NOT NULL CHECK (source IN ('websub', 'task', 'manual')),
  created_at          timestamptz DEFAULT now()
);
CREATE UNIQUE INDEX IF NOT EXISTS video_snapshots_video_checkpoint_unique ON ingestion.video_snapshots(video_id, checkpoint_hour);
CREATE INDEX IF NOT EXISTS video_snapshots_video_id_idx ON ingestion.video_snapshots(video_id);
CREATE INDEX IF NOT EXISTS video_snapshots_checkpoint_hour_idx ON ingestion.video_snapshots(checkpoint_hour);
CREATE INDEX IF NOT EXISTS video_snapshots_measured_at_idx ON ingestion.video_snapshots(measured_at DESC);
CREATE INDEX IF NOT EXISTS video_snapshots_source_idx ON ingestion.video_snapshots(source);

-- Snapshot tasks table (for scheduled snapshots)
CREATE TABLE IF NOT EXISTS ingestion.snapshot_tasks (
  video_id         uuid NOT NULL REFERENCES ingestion.videos(id) ON DELETE CASCADE,
  checkpoint_hour  integer NOT NULL CHECK (checkpoint_hour IN (0, 3, 6, 12, 24, 48, 72, 168)),
  scheduled_at     timestamptz NOT NULL,
  PRIMARY KEY (video_id, checkpoint_hour)
);
CREATE INDEX IF NOT EXISTS snapshot_tasks_scheduled_at_idx ON ingestion.snapshot_tasks(scheduled_at);
CREATE INDEX IF NOT EXISTS snapshot_tasks_checkpoint_hour_idx ON ingestion.snapshot_tasks(checkpoint_hour);