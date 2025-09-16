-- Up migration: create ingestion schema and tables
CREATE SCHEMA IF NOT EXISTS ingestion;

-- YouTube categories table
CREATE TABLE IF NOT EXISTS ingestion.youtube_categories (
  id         integer PRIMARY KEY,
  name       varchar(100) NOT NULL,
  assignable boolean NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

-- Genres table
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

-- Channels table
CREATE TABLE IF NOT EXISTS ingestion.channels (
  id                  uuid PRIMARY KEY,
  youtube_channel_id  text NOT NULL,
  title               text NOT NULL,
  description         text,
  country             varchar(10),
  view_count          bigint,
  subscription_count  bigint,
  video_count         bigint,
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
  created_at          timestamptz DEFAULT now(),
  updated_at          timestamptz,
  deleted_at          timestamptz
);
CREATE UNIQUE INDEX IF NOT EXISTS videos_youtube_video_id_unique ON ingestion.videos(youtube_video_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS videos_channel_id_idx ON ingestion.videos(channel_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS videos_published_at_idx ON ingestion.videos(published_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS videos_category_id_idx ON ingestion.videos(category_id) WHERE deleted_at IS NULL;

-- Video genres junction table
CREATE TABLE IF NOT EXISTS ingestion.video_genres (
  id         uuid PRIMARY KEY,
  video_id   uuid NOT NULL REFERENCES ingestion.videos(id) ON DELETE CASCADE,
  genre_id   uuid NOT NULL REFERENCES ingestion.genres(id) ON DELETE CASCADE,
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX idx_video_genres_video ON ingestion.video_genres(video_id);
CREATE INDEX idx_video_genres_genre ON ingestion.video_genres(genre_id);
CREATE UNIQUE INDEX video_genres_video_genre_unique ON ingestion.video_genres(video_id, genre_id);

-- Keyword groups table (for organizing keywords)
CREATE TABLE IF NOT EXISTS ingestion.keyword_groups (
  id           uuid PRIMARY KEY,
  genre_id     uuid NOT NULL REFERENCES ingestion.genres(id) ON DELETE CASCADE,
  name         text NOT NULL,
  filter_type  text NOT NULL CHECK (filter_type IN ('include', 'exclude')),
  target_field varchar(20) NOT NULL DEFAULT 'title',
  enabled      boolean DEFAULT true,
  description  text,
  created_at   timestamptz DEFAULT now(),
  updated_at   timestamptz,
  deleted_at   timestamptz
);
CREATE UNIQUE INDEX IF NOT EXISTS keyword_groups_genre_name_filter_unique ON ingestion.keyword_groups(genre_id, name, filter_type) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS keyword_groups_filter_type_idx ON ingestion.keyword_groups(filter_type) WHERE deleted_at IS NULL AND enabled = true;
CREATE INDEX IF NOT EXISTS keyword_groups_enabled_idx ON ingestion.keyword_groups(enabled) WHERE deleted_at IS NULL;
CREATE INDEX idx_keyword_groups_genre_enabled ON ingestion.keyword_groups(genre_id, enabled);

-- Keyword items table (individual keywords)
CREATE TABLE IF NOT EXISTS ingestion.keyword_items (
  id                uuid PRIMARY KEY,
  keyword_group_id  uuid NOT NULL REFERENCES ingestion.keyword_groups(id) ON DELETE CASCADE,
  keyword           text NOT NULL,
  created_at        timestamptz DEFAULT now(),
  updated_at        timestamptz
);
CREATE INDEX IF NOT EXISTS keyword_items_group_id_idx ON ingestion.keyword_items(keyword_group_id);
CREATE UNIQUE INDEX IF NOT EXISTS keyword_items_group_keyword_unique ON ingestion.keyword_items(keyword_group_id, keyword);

-- Channel snapshots table
CREATE TABLE IF NOT EXISTS ingestion.channel_snapshots (
  id                  uuid PRIMARY KEY,
  channel_id          uuid NOT NULL REFERENCES ingestion.channels(id) ON DELETE CASCADE,
  measured_at         timestamptz NOT NULL,
  subscription_count  integer NOT NULL,
  view_count          bigint,
  video_count         bigint,
  created_at          timestamptz DEFAULT now(),
  updated_at          timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS channel_snapshots_channel_id_idx ON ingestion.channel_snapshots(channel_id);
CREATE INDEX IF NOT EXISTS channel_snapshots_measured_at_idx ON ingestion.channel_snapshots(measured_at DESC);
ALTER TABLE ingestion.channel_snapshots 
  ADD CONSTRAINT channel_snapshots_channel_measured_unique UNIQUE (channel_id, measured_at);

-- Video snapshots table
CREATE TABLE IF NOT EXISTS ingestion.video_snapshots (
  id                  uuid PRIMARY KEY,
  video_id            uuid NOT NULL REFERENCES ingestion.videos(id) ON DELETE CASCADE,
  checkpoint_hour     integer NOT NULL CHECK (checkpoint_hour IN (0, 3, 6, 12, 24, 48, 72, 168)),
  measured_at         timestamptz NOT NULL,
  view_count          bigint NOT NULL DEFAULT 0,
  like_count          bigint NOT NULL DEFAULT 0,
  subscription_count  bigint NOT NULL DEFAULT 0,
  source              text NOT NULL CHECK (source IN ('websub', 'task', 'manual')),
  created_at          timestamptz DEFAULT now(),
  updated_at          timestamptz NOT NULL DEFAULT now()
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

-- Audit logs table
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

-- Batch jobs table
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