# Database Schema Design

## Overview

Database design using PostgreSQL (Neon). Based on microservices architecture, tables are separated by each Bounded Context.

## Table Design (Neon / PostgreSQL)

### Auxiliary: YouTube Category Dictionary (Optional)

```sql
CREATE TABLE video_categories (
  youtube_category_id int PRIMARY KEY,         -- e.g.: 27=Education, 28=Science&Tech
  name                text NOT NULL,
  created_at          timestamptz DEFAULT now(),
  updated_at          timestamptz,
  deleted_at          timestamptz
);
```

### Ingestion Service (Collection & Storage)

```sql
-- Schema creation
CREATE SCHEMA IF NOT EXISTS ingestion;
SET search_path TO ingestion;

-- =========================================================
-- 1) Keywords (Include/Exclude Rules)
-- =========================================================
CREATE TABLE keywords (
  id               uuid PRIMARY KEY,                             -- Internal UUID (v7)
  name             text NOT NULL,
  filter_type      text NOT NULL CHECK (filter_type IN ('include','exclude')),
  pattern          text NOT NULL,                                 -- Normalized regex
  enabled          boolean NOT NULL DEFAULT true,
  description      text,
  created_at       timestamptz NOT NULL DEFAULT now(),
  updated_at       timestamptz,
  deleted_at       timestamptz
);

-- Logical uniqueness (excluding deleted): name + filter_type
CREATE UNIQUE INDEX keywords_name_type_uq
  ON keywords (lower(name), filter_type)
  WHERE deleted_at IS NULL;

-- =========================================================
-- 2) Channels (Subscribed Channels) + Snapshots (Subscriber Trends)
-- =========================================================
CREATE TABLE channels (
  id                 uuid PRIMARY KEY,                           -- Internal UUID
  youtube_channel_id text NOT NULL UNIQUE,                       -- External ID (YouTube channelId)
  title              text,
  thumbnail_url      text,
  subscribed         boolean NOT NULL DEFAULT false,             -- WebSub subscription target
  created_at         timestamptz NOT NULL DEFAULT now(),
  updated_at         timestamptz,
  deleted_at         timestamptz
);

CREATE TABLE channel_snapshots (
  id                  uuid PRIMARY KEY,                          -- Snapshot UUID
  channel_id          uuid NOT NULL REFERENCES channels(id),
  measured_at         timestamptz NOT NULL,
  subscription_count  integer NOT NULL,
  created_at          timestamptz NOT NULL DEFAULT now(),
  UNIQUE (channel_id, measured_at)                               -- insert-only
);
CREATE INDEX channel_snapshots_recent_idx
  ON channel_snapshots (channel_id, measured_at DESC);

-- =========================================================
-- 3) Videos (Monitored Videos) + Snapshots (Checkpoint Measurements)
-- =========================================================
CREATE TABLE videos (
  id                   uuid PRIMARY KEY,                         -- Internal UUID
  youtube_video_id     text NOT NULL UNIQUE,                     -- External ID (YouTube videoId)
  channel_id           uuid NOT NULL REFERENCES channels(id),    -- Internal FK
  youtube_channel_id   text NOT NULL,                            -- Redundant storage (JOIN optimization)
  title                text,
  published_at         timestamptz NOT NULL,
  category_id          integer,                                  -- YouTube categoryId (e.g.: 27,28)
  thumbnail_url        text,
  video_url            text,
  created_at           timestamptz NOT NULL DEFAULT now(),
  updated_at           timestamptz,
  deleted_at           timestamptz
);
CREATE INDEX videos_channel_idx   ON videos (channel_id);
CREATE INDEX videos_published_idx ON videos (published_at);

CREATE TABLE video_snapshots (
  id                   uuid PRIMARY KEY,                         -- Snapshot UUID
  video_id             uuid NOT NULL REFERENCES videos(id),
  checkpoint_hour      smallint NOT NULL
       CHECK (checkpoint_hour IN (0,3,6,12,24,48,72,168)),       -- 0h/3h/.../7d
  measured_at          timestamptz NOT NULL,
  views_count          bigint NOT NULL,
  likes_count          bigint NOT NULL,
  subscription_count   bigint NOT NULL,                          -- Channel subs at the time (copy)
  source               text NOT NULL DEFAULT 'task',             -- 'websub'|'task'|'manual'
  created_at           timestamptz NOT NULL DEFAULT now(),
  UNIQUE (video_id, checkpoint_hour)                             -- Idempotency (ON CONFLICT DO NOTHING)
);
CREATE INDEX video_snapshots_vh_idx   ON video_snapshots (video_id, checkpoint_hour);
CREATE INDEX video_snapshots_meas_idx ON video_snapshots (video_id, measured_at);

-- =========================================================
-- 4) WebSub Subscriptions (Optional: Lease Management)
-- =========================================================
-- Only needed if you want to track WebSub lease expiry in DB
-- CREATE TABLE websub_subscriptions (
--   channel_id       uuid PRIMARY KEY REFERENCES channels(id),
--   lease_expires_at timestamptz NOT NULL,
--   hub_callback_url text NOT NULL,
--   updated_at       timestamptz NOT NULL DEFAULT now()
-- );
```

#### Key Design Decisions

- **Internal IDs**: All internal PKs use UUID v7 (generated in application)
- **External IDs**: YouTube IDs stored in `youtube_*` columns with UNIQUE constraints
- **Idempotency**:
  - `video_snapshots`: `(video_id, checkpoint_hour) UNIQUE` → `INSERT ... ON CONFLICT DO NOTHING`
  - `videos.youtube_video_id UNIQUE` → Absorbs trending duplicates
  - `channels.youtube_channel_id UNIQUE` → Absorbs subscription duplicates
- **Schema Separation**: Uses `ingestion` schema for isolation
- **Indexes**: Optimized for common queries (recent snapshots, published date ranges)
