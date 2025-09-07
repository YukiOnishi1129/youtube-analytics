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

### Ingestion (Collection & Storage)

#### A-1. Channels (Subscription & Profile)

```sql
CREATE TABLE channels (
  id               text PRIMARY KEY,                 -- YouTube channelId
  title            text,
  description      text,
  thumbnail_url    text,
  subscribed       boolean DEFAULT false,            -- WebSub subscription target
  created_at       timestamptz DEFAULT now(),
  updated_at       timestamptz,
  deleted_at       timestamptz
);
```

#### A-2. Channel Subscriber Snapshots (Optional: Daily)

```sql
CREATE TABLE channel_snapshots (
  channel_id         text NOT NULL REFERENCES channels(id),
  measured_at        timestamptz NOT NULL,
  subscription_count bigint NOT NULL,
  PRIMARY KEY (channel_id, measured_at)
);
CREATE INDEX channel_snapshots_idx ON channel_snapshots (channel_id, measured_at DESC);
```

#### A-3. Video Metadata

```sql
CREATE TABLE videos (
  id                   text PRIMARY KEY,             -- YouTube videoId
  channel_id           text NOT NULL REFERENCES channels(id),
  title                text,
  description          text,
  published_at         timestamptz NOT NULL,
  youtube_category_id  int REFERENCES video_categories(youtube_category_id),
  thumbnail_url        text,
  video_url            text,                         -- https://www.youtube.com/watch?v={id}
  created_at           timestamptz DEFAULT now(),
  updated_at           timestamptz,
  deleted_at           timestamptz
);
CREATE INDEX videos_channel_idx   ON videos (channel_id);
CREATE INDEX videos_published_idx ON videos (published_at);
```

#### A-4. Video Snapshots (Each Checkpoint: insert-only)

```sql
CREATE TABLE video_snapshots (
  video_id            text NOT NULL REFERENCES videos(id),
  checkpoint_hour     smallint NOT NULL,             -- 0,3,6,12,24,48,72,168
  CHECK (checkpoint_hour IN (0,3,6,12,24,48,72,168)),

  measured_at         timestamptz NOT NULL,          -- Actual measurement time (aligned to CP)
  created_at          timestamptz NOT NULL DEFAULT now(),

  views_count         bigint NOT NULL,
  likes_count         bigint NOT NULL,
  comment_count       bigint,
  subscription_count  bigint NOT NULL,               -- Channel subscriber count at the time (embedded)
  source              text NOT NULL DEFAULT 'task',  -- 'websub'|'task'|'manual'

  PRIMARY KEY (video_id, checkpoint_hour)
);
CREATE INDEX video_snapshots_measured_idx ON video_snapshots (video_id, measured_at);
CREATE INDEX video_snapshots_cp_idx       ON video_snapshots (checkpoint_hour);
```

#### A-5. Filter Keywords (Single Table Operation & Soft Delete)

```sql
CREATE TABLE video_filter_keywords (
  id               uuid PRIMARY KEY,                            -- v7 generated in Go
  name             text NOT NULL,                               -- Display/aggregation name (e.g.: 'next.js')
  filter_type      text NOT NULL CHECK (filter_type IN ('include','exclude')),
  pattern          text NOT NULL,                               -- regex from BuildPattern(name)
  enabled          boolean DEFAULT true,
  description      text,
  created_at       timestamptz DEFAULT now(),
  updated_at       timestamptz,
  deleted_at       timestamptz
);
CREATE INDEX vfk_enabled_idx     ON video_filter_keywords (enabled) WHERE deleted_at IS NULL;
CREATE INDEX vfk_filter_type_idx ON video_filter_keywords (filter_type) WHERE deleted_at IS NULL;
```