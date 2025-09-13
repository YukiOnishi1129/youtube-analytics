-- name: CreateChannel :exec
INSERT INTO ingestion.channels (
    id, youtube_channel_id, title, thumbnail_url, subscribed, created_at
) VALUES ($1, $2, $3, $4, $5, $6);

-- name: UpdateChannel :exec
UPDATE ingestion.channels
SET title = $2, thumbnail_url = $3, subscribed = $4, updated_at = $5, deleted_at = $6
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetChannelByID :one
SELECT id, youtube_channel_id, title, thumbnail_url, subscribed, created_at, updated_at
FROM ingestion.channels
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetChannelByYouTubeID :one
SELECT id, youtube_channel_id, title, thumbnail_url, subscribed, created_at, updated_at
FROM ingestion.channels
WHERE youtube_channel_id = $1 AND deleted_at IS NULL;

-- name: ListChannels :many
SELECT id, youtube_channel_id, title, thumbnail_url, subscribed, created_at, updated_at
FROM ingestion.channels
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListSubscribedChannels :many
SELECT id, youtube_channel_id, title, thumbnail_url, subscribed, created_at, updated_at
FROM ingestion.channels
WHERE subscribed = true AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: ListActiveChannels :many
SELECT id, youtube_channel_id, title, thumbnail_url, subscribed, created_at, updated_at
FROM ingestion.channels
WHERE deleted_at IS NULL
ORDER BY created_at DESC;

-- name: CountChannels :one
SELECT COUNT(*) FROM ingestion.channels WHERE deleted_at IS NULL;

-- name: CreateVideo :exec
INSERT INTO ingestion.videos (
    id, youtube_video_id, channel_id, youtube_channel_id, title,
    published_at, category_id, thumbnail_url, video_url, created_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: GetVideoByID :one
SELECT id, youtube_video_id, youtube_channel_id, title, thumbnail_url, published_at, created_at
FROM ingestion.videos
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetVideoByYouTubeID :one
SELECT id, youtube_video_id, youtube_channel_id, title, thumbnail_url, published_at, created_at
FROM ingestion.videos
WHERE youtube_video_id = $1 AND deleted_at IS NULL;

-- name: CheckVideoExists :one
SELECT EXISTS(
    SELECT 1 FROM ingestion.videos 
    WHERE youtube_video_id = $1 AND deleted_at IS NULL
);

-- name: ListVideosByChannel :many
SELECT id, youtube_video_id, youtube_channel_id, title, thumbnail_url, published_at, created_at
FROM ingestion.videos
WHERE channel_id = $1 AND deleted_at IS NULL
ORDER BY published_at DESC
LIMIT $2 OFFSET $3;

-- name: ListActiveVideos :many
SELECT id, youtube_video_id, youtube_channel_id, title, thumbnail_url, published_at, created_at
FROM ingestion.videos
WHERE published_at > $1 AND deleted_at IS NULL
ORDER BY published_at DESC;

-- name: CountVideosByChannel :one
SELECT COUNT(*) FROM ingestion.videos
WHERE channel_id = $1 AND deleted_at IS NULL;

-- name: CreateVideoSnapshot :exec
INSERT INTO ingestion.video_snapshots (
    id, video_id, checkpoint_hour, measured_at, views_count,
    likes_count, subscription_count, source, created_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: GetVideoSnapshotByVideoAndCheckpoint :one
SELECT id, video_id, checkpoint_hour, measured_at, views_count, 
       likes_count, subscription_count, source, created_at
FROM ingestion.video_snapshots
WHERE video_id = $1 AND checkpoint_hour = $2;

-- name: ListVideoSnapshots :many
SELECT id, video_id, checkpoint_hour, measured_at, views_count, 
       likes_count, subscription_count, source, created_at
FROM ingestion.video_snapshots
WHERE video_id = $1
ORDER BY checkpoint_hour ASC;

-- name: CreateKeyword :exec
INSERT INTO ingestion.keywords (
    id, name, filter_type, pattern, enabled, description, created_at
) VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: GetKeywordByID :one
SELECT id, name, filter_type, pattern, enabled, description, created_at, updated_at
FROM ingestion.keywords
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListEnabledKeywords :many
SELECT id, name, filter_type, pattern, enabled, description, created_at, updated_at
FROM ingestion.keywords
WHERE enabled = true AND deleted_at IS NULL
ORDER BY name ASC;

-- name: CreateSnapshotTask :exec
INSERT INTO ingestion.snapshot_tasks (
    video_id, checkpoint_hour, scheduled_at
) VALUES ($1, $2, $3);

-- name: DeleteSnapshotTask :exec
DELETE FROM ingestion.snapshot_tasks
WHERE video_id = $1 AND checkpoint_hour = $2;

-- name: GetPendingSnapshotTasks :many
SELECT video_id, checkpoint_hour, scheduled_at
FROM ingestion.snapshot_tasks
WHERE scheduled_at <= $1
ORDER BY scheduled_at ASC
LIMIT $2;