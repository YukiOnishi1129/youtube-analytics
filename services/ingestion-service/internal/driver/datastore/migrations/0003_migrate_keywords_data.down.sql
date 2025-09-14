-- Revert data migration

-- Remove genre_id from keywords (handled in 0002_genre_based_design.down.sql)

-- Remove the default genre
DELETE FROM ingestion.genres WHERE code = 'engineering_jp';