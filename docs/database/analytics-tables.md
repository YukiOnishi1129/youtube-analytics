# Analytics (Processing & Serving) Tables

## Metrics by Checkpoint (Read Model for Rankings)

```sql
CREATE TABLE video_metrics_checkpoint (
  video_id                           text NOT NULL REFERENCES videos(id),
  checkpoint_hour                    smallint NOT NULL,
  CHECK (checkpoint_hour IN (3,6,12,24,48,72,168)),

  -- Redundantly stored for fast Published range search
  published_at                       timestamptz NOT NULL,

  -- Raw data at point X (display auxiliary)
  views_count                        bigint   NOT NULL,  -- views@X
  likes_count                        bigint   NOT NULL,  -- likes@X
  subscription_count                 bigint   NOT NULL,  -- subs@X

  -- 0h (baseline)
  views_baseline_count               bigint   NOT NULL,  -- views@0
  likes_baseline_count               bigint   NOT NULL,  -- likes@0
  subscription_baseline_count        bigint   NOT NULL,  -- subs@0

  -- Main metrics (0→X per-hour increment, point ratio, quality/heat)
  view_growth_rate_per_hour          double precision NOT NULL,  -- (viewsX - views0)/X[h]
  like_growth_rate_per_hour          double precision NOT NULL,  -- (likesX - likes0)/X[h]
  like_growth_rate_per_subscription_per_hour double precision,   -- optional display

  views_per_subscription_rate        double precision NOT NULL,  -- viewsX / subsX
  like_rate_at_checkpoint            double precision,           -- likesX / viewsX (reference)
  wilson_like_rate_lower_bound       double precision,           -- Wilson lower bound for like rate
  likes_per_subscription_shrunk_rate double precision,           -- SCALE*likesX/(subsX+OFFSET)

  -- Exclude from ranking for low samples, etc.
  exclude_from_ranking               boolean DEFAULT false,

  computed_at                        timestamptz DEFAULT now(),

  PRIMARY KEY (video_id, checkpoint_hour)
);

-- Indexes for typical queries
CREATE INDEX vmc_idx_pub_cp    ON video_metrics_checkpoint (published_at, checkpoint_hour);
CREATE INDEX vmc_idx_cp_view   ON video_metrics_checkpoint (checkpoint_hour, view_growth_rate_per_hour DESC);
CREATE INDEX vmc_idx_cp_rel    ON video_metrics_checkpoint (checkpoint_hour, views_per_subscription_rate DESC);
CREATE INDEX vmc_idx_cp_wilson ON video_metrics_checkpoint (checkpoint_hour, wilson_like_rate_lower_bound DESC);
CREATE INDEX vmc_idx_cp_lps    ON video_metrics_checkpoint (checkpoint_hour, likes_per_subscription_shrunk_rate DESC);
```

## Ranking History (TopN Frozen)

```sql
CREATE TABLE ranking_snapshots (
  id              uuid PRIMARY KEY,                     -- v7
  snapshot_at     timestamptz NOT NULL,                -- Capture time
  ranking_kind    text NOT NULL CHECK (ranking_kind IN
                      ('speed_views','speed_likes','relative_views','quality','heat')),
  checkpoint_hour smallint NOT NULL CHECK (checkpoint_hour IN (3,6,12,24,48,72,168)),
  published_from  timestamptz NOT NULL,
  published_to    timestamptz NOT NULL,
  top_n           int NOT NULL DEFAULT 10,
  created_at      timestamptz DEFAULT now(),
  updated_at      timestamptz,
  deleted_at      timestamptz
);
CREATE INDEX rs_idx_when ON ranking_snapshots (snapshot_at DESC);
CREATE INDEX rs_idx_meta ON ranking_snapshots (ranking_kind, checkpoint_hour, snapshot_at DESC);

CREATE TABLE ranking_snapshot_items (
  id              uuid PRIMARY KEY,                     -- v7
  snapshot_id     uuid NOT NULL REFERENCES ranking_snapshots(id) ON DELETE CASCADE,
  rank            int NOT NULL,
  video_id        text NOT NULL,                        -- YouTube videoId
  title           text NOT NULL,                        -- Frozen for display at that time
  channel_id      text NOT NULL,
  thumbnail_url   text,
  video_url       text,
  published_at    timestamptz NOT NULL,
  checkpoint_hour smallint NOT NULL,
  main_metric     double precision NOT NULL,            -- Main metric used for ordering (frozen)
  views_count     bigint,
  likes_count     bigint,
  created_at      timestamptz DEFAULT now(),
  UNIQUE (snapshot_id, rank)
);
CREATE INDEX rsi_idx_snapshot ON ranking_snapshot_items (snapshot_id, rank);
CREATE INDEX rsi_idx_video    ON ranking_snapshot_items (video_id);
```