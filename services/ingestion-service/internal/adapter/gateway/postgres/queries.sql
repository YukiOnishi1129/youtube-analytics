-- name: CreateChannel :exec
INSERT INTO ingestion.channels (
    id, youtube_channel_id, title, thumbnail_url, description, country, 
    view_count, subscription_count, video_count, subscribed, created_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);

-- name: UpdateChannel :exec
UPDATE ingestion.channels
SET title = $2, thumbnail_url = $3, description = $4, country = $5,
    view_count = $6, subscription_count = $7, video_count = $8,
    subscribed = $9, updated_at = $10, deleted_at = $11
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetChannelByID :one
SELECT id, youtube_channel_id, title, thumbnail_url, description, country,
       view_count, subscription_count, video_count, subscribed, created_at, updated_at
FROM ingestion.channels
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetChannelByYouTubeID :one
SELECT id, youtube_channel_id, title, thumbnail_url, description, country,
       view_count, subscription_count, video_count, subscribed, created_at, updated_at
FROM ingestion.channels
WHERE youtube_channel_id = $1 AND deleted_at IS NULL;

-- name: ListChannels :many
SELECT id, youtube_channel_id, title, thumbnail_url, description, country,
       view_count, subscription_count, video_count, subscribed, created_at, updated_at
FROM ingestion.channels
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListSubscribedChannels :many
SELECT id, youtube_channel_id, title, thumbnail_url, description, country,
       view_count, subscription_count, video_count, subscribed, created_at, updated_at
FROM ingestion.channels
WHERE subscribed = true AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: ListActiveChannels :many
SELECT id, youtube_channel_id, title, thumbnail_url, description, country,
       view_count, subscription_count, video_count, subscribed, created_at, updated_at
FROM ingestion.channels
WHERE deleted_at IS NULL
ORDER BY created_at DESC;

-- name: CountChannels :one
SELECT COUNT(*) FROM ingestion.channels WHERE deleted_at IS NULL;

-- name: CreateVideo :exec
INSERT INTO ingestion.videos (
    id, youtube_video_id, channel_id, youtube_channel_id, title,
    published_at, category_id, created_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetVideoByID :one
SELECT id, youtube_video_id, channel_id, youtube_channel_id, title, published_at, category_id, created_at
FROM ingestion.videos
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetVideoByYouTubeID :one
SELECT id, youtube_video_id, channel_id, youtube_channel_id, title, published_at, category_id, created_at
FROM ingestion.videos
WHERE youtube_video_id = $1 AND deleted_at IS NULL;

-- name: CheckVideoExists :one
SELECT EXISTS(
    SELECT 1 FROM ingestion.videos 
    WHERE youtube_video_id = $1 AND deleted_at IS NULL
);

-- name: ListVideosByChannel :many
SELECT id, youtube_video_id, channel_id, youtube_channel_id, title, published_at, category_id, created_at
FROM ingestion.videos
WHERE channel_id = $1 AND deleted_at IS NULL
ORDER BY published_at DESC
LIMIT $2 OFFSET $3;

-- name: ListActiveVideos :many
SELECT id, youtube_video_id, channel_id, youtube_channel_id, title, published_at, category_id, created_at
FROM ingestion.videos
WHERE published_at > $1 AND deleted_at IS NULL
ORDER BY published_at DESC;

-- name: CountVideosByChannel :one
SELECT COUNT(*) FROM ingestion.videos
WHERE channel_id = $1 AND deleted_at IS NULL;

-- name: CreateVideoSnapshot :exec
INSERT INTO ingestion.video_snapshots (
    id, video_id, checkpoint_hour, measured_at, view_count,
    like_count, subscription_count, source, created_at, updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: GetVideoSnapshotByVideoAndCheckpoint :one
SELECT id, video_id, checkpoint_hour, measured_at, view_count, 
       like_count, subscription_count, source, created_at, updated_at
FROM ingestion.video_snapshots
WHERE video_id = $1 AND checkpoint_hour = $2;

-- name: ListVideoSnapshots :many
SELECT id, video_id, checkpoint_hour, measured_at, view_count, 
       like_count, subscription_count, source, created_at, updated_at
FROM ingestion.video_snapshots
WHERE video_id = $1
ORDER BY checkpoint_hour ASC;

-- name: CreateChannelSnapshot :exec
INSERT INTO ingestion.channel_snapshots (
    id, channel_id, measured_at, subscription_count, view_count, video_count, created_at, updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetLatestChannelSnapshot :one
SELECT id, channel_id, measured_at, subscription_count, view_count, video_count, created_at, updated_at
FROM ingestion.channel_snapshots
WHERE channel_id = $1
ORDER BY measured_at DESC
LIMIT 1;

-- name: ListChannelSnapshots :many
SELECT id, channel_id, measured_at, subscription_count, view_count, video_count, created_at, updated_at
FROM ingestion.channel_snapshots
WHERE channel_id = $1
ORDER BY measured_at DESC
LIMIT $2;

-- name: CreateKeyword :exec
INSERT INTO ingestion.keywords (
    id, genre_id, name, filter_type, pattern, target_field, enabled, description, created_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: GetKeywordByID :one
SELECT id, genre_id, name, filter_type, pattern, target_field, enabled, description, created_at, updated_at
FROM ingestion.keywords
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListEnabledKeywords :many
SELECT id, genre_id, name, filter_type, pattern, target_field, enabled, description, created_at, updated_at
FROM ingestion.keywords
WHERE enabled = true AND deleted_at IS NULL
ORDER BY name ASC;

-- name: UpdateKeyword :exec
UPDATE ingestion.keywords
SET name = $2, filter_type = $3, pattern = $4, target_field = $5, 
    enabled = $6, description = $7, updated_at = $8
WHERE id = $1 AND deleted_at IS NULL;

-- name: SoftDeleteKeyword :exec
UPDATE ingestion.keywords
SET deleted_at = $2, updated_at = $2
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListKeywordsByGenre :many
SELECT id, genre_id, name, filter_type, pattern, target_field, enabled, description, created_at, updated_at
FROM ingestion.keywords
WHERE genre_id = $1 AND ($2::boolean IS NULL OR enabled = $2) AND deleted_at IS NULL
ORDER BY name ASC;

-- name: ListKeywordsByGenreAndType :many
SELECT id, genre_id, name, filter_type, pattern, target_field, enabled, description, created_at, updated_at
FROM ingestion.keywords
WHERE genre_id = $1 AND filter_type = $2 AND ($3::boolean IS NULL OR enabled = $3) AND deleted_at IS NULL
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

-- Genre queries
-- name: CreateGenre :exec
INSERT INTO ingestion.genres (
    id, code, name, language, region_code, category_ids, enabled, created_at, updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: UpdateGenre :exec
UPDATE ingestion.genres
SET name = $2, category_ids = $3, enabled = $4, updated_at = $5
WHERE id = $1;

-- name: GetGenreByID :one
SELECT id, code, name, language, region_code, category_ids, enabled, created_at, updated_at
FROM ingestion.genres
WHERE id = $1;

-- name: GetGenreByCode :one
SELECT id, code, name, language, region_code, category_ids, enabled, created_at, updated_at
FROM ingestion.genres
WHERE code = $1;

-- name: ListGenres :many
SELECT id, code, name, language, region_code, category_ids, enabled, created_at, updated_at
FROM ingestion.genres
ORDER BY code ASC;

-- name: ListEnabledGenres :many
SELECT id, code, name, language, region_code, category_ids, enabled, created_at, updated_at
FROM ingestion.genres
WHERE enabled = true
ORDER BY code ASC;

-- YouTube Category queries
-- name: CreateYouTubeCategory :exec
INSERT INTO ingestion.youtube_categories (
    id, name, assignable, created_at, updated_at
) VALUES ($1, $2, $3, $4, $5);

-- name: UpdateYouTubeCategory :exec
UPDATE ingestion.youtube_categories
SET name = $2, assignable = $3, updated_at = $4
WHERE id = $1;

-- name: GetYouTubeCategoryByID :one
SELECT id, name, assignable, created_at, updated_at
FROM ingestion.youtube_categories
WHERE id = $1;

-- name: ListYouTubeCategories :many
SELECT id, name, assignable, created_at, updated_at
FROM ingestion.youtube_categories
ORDER BY id ASC;

-- name: ListAssignableYouTubeCategories :many
SELECT id, name, assignable, created_at, updated_at
FROM ingestion.youtube_categories
WHERE assignable = true
ORDER BY id ASC;

-- Video Genre queries
-- name: CreateVideoGenre :exec
INSERT INTO ingestion.video_genres (
    id, video_id, genre_id, created_at
) VALUES ($1, $2, $3, $4);

-- name: ListVideoGenresByVideo :many
SELECT id, video_id, genre_id, created_at
FROM ingestion.video_genres
WHERE video_id = $1;

-- name: ListVideoGenresByGenre :many
SELECT id, video_id, genre_id, created_at
FROM ingestion.video_genres
WHERE genre_id = $1;

-- name: CheckVideoGenreExists :one
SELECT EXISTS(
    SELECT 1 FROM ingestion.video_genres
    WHERE video_id = $1 AND genre_id = $2
);

-- name: DeleteVideoGenresByVideo :exec
DELETE FROM ingestion.video_genres
WHERE video_id = $1;

-- name: DeleteVideoGenresByGenre :exec
DELETE FROM ingestion.video_genres
WHERE genre_id = $1;

-- Audit Log queries
-- name: CreateAuditLog :exec
INSERT INTO ingestion.audit_logs (
    id, actor_id, actor_email, action, resource_type, resource_id,
    old_values, new_values, ip_address, user_agent, created_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);

-- name: ListAuditLogsByActor :many
SELECT id, actor_id, actor_email, action, resource_type, resource_id,
       old_values, new_values, ip_address, user_agent, created_at
FROM ingestion.audit_logs
WHERE actor_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListAuditLogsByResource :many
SELECT id, actor_id, actor_email, action, resource_type, resource_id,
       old_values, new_values, ip_address, user_agent, created_at
FROM ingestion.audit_logs
WHERE resource_type = $1 AND resource_id = $2
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: ListRecentAuditLogs :many
SELECT id, actor_id, actor_email, action, resource_type, resource_id,
       old_values, new_values, ip_address, user_agent, created_at
FROM ingestion.audit_logs
ORDER BY created_at DESC
LIMIT $1;

-- Batch Job queries
-- name: CreateBatchJob :exec
INSERT INTO ingestion.batch_jobs (
    id, job_type, status, parameters, created_at
) VALUES ($1, $2, $3, $4, $5);

-- name: UpdateBatchJob :exec
UPDATE ingestion.batch_jobs
SET status = $2, started_at = $3, completed_at = $4, error_message = $5, statistics = $6
WHERE id = $1;

-- name: GetBatchJobByID :one
SELECT id, job_type, status, parameters, started_at, completed_at,
       error_message, statistics, created_at
FROM ingestion.batch_jobs
WHERE id = $1;

-- name: ListBatchJobsByTypeAndStatus :many
SELECT id, job_type, status, parameters, started_at, completed_at,
       error_message, statistics, created_at
FROM ingestion.batch_jobs
WHERE job_type = $1 AND status = $2
ORDER BY created_at DESC;

-- name: ListRecentBatchJobs :many
SELECT id, job_type, status, parameters, started_at, completed_at,
       error_message, statistics, created_at
FROM ingestion.batch_jobs
WHERE ($1::text IS NULL OR job_type = $1)
ORDER BY created_at DESC
LIMIT $2;

-- name: ListRunningBatchJobs :many
SELECT id, job_type, status, parameters, started_at, completed_at,
       error_message, statistics, created_at
FROM ingestion.batch_jobs
WHERE status = 'running'
ORDER BY started_at ASC;