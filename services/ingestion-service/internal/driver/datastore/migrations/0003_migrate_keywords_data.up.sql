-- Data migration for existing keywords

-- First, create a default genre for existing keywords
INSERT INTO ingestion.genres (id, code, name, language, region_code, category_ids, enabled)
VALUES (
  '550e8400-e29b-41d4-a716-446655440000'::uuid,
  'engineering_jp',
  'Engineering (JP)',
  'ja',
  'JP',
  ARRAY[27, 28],
  true
) ON CONFLICT (code) DO NOTHING;

-- Update all existing keywords to belong to the default genre
UPDATE ingestion.keywords
SET genre_id = '550e8400-e29b-41d4-a716-446655440000'::uuid
WHERE genre_id IS NULL;