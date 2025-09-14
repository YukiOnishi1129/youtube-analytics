-- Add NOT NULL constraint to genre_id after data migration
ALTER TABLE ingestion.keywords
  ALTER COLUMN genre_id SET NOT NULL;