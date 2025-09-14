-- Remove NOT NULL constraint from genre_id
ALTER TABLE ingestion.keywords
  ALTER COLUMN genre_id DROP NOT NULL;