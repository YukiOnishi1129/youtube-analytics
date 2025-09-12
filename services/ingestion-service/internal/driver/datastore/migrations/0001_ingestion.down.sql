-- Down migration: drop ingestion tables and schema
DROP TABLE IF EXISTS ingestion.snapshot_tasks;
DROP TABLE IF EXISTS ingestion.video_snapshots;
DROP TABLE IF EXISTS ingestion.channel_snapshots;
DROP TABLE IF EXISTS ingestion.keywords;
DROP TABLE IF EXISTS ingestion.videos;
DROP TABLE IF EXISTS ingestion.channels;

DROP SCHEMA IF EXISTS ingestion;