-- Genres Seed Data
-- Initial launch with Japanese Engineering genre only

-- Japanese Engineering genre (combining Education and Science & Technology categories)
INSERT INTO genres (id, code, name, language, region_code, category_ids, enabled, created_at, updated_at) VALUES
    ('550e8400-e29b-41d4-a716-446655440001', 'engineering_ja_jp', 'エンジニア', 'ja', 'JP', '{27,28}', true, NOW(), NOW())
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    language = EXCLUDED.language,
    region_code = EXCLUDED.region_code,
    category_ids = EXCLUDED.category_ids,
    enabled = EXCLUDED.enabled,
    updated_at = NOW();

-- Future genres (commented out for initial launch)
-- Uncomment these as needed when expanding to other genres/regions

-- -- Technology genres (other regions)
-- INSERT INTO genres (id, code, name, language, region_code, category_ids, enabled, created_at, updated_at) VALUES
--     ('550e8400-e29b-41d4-a716-446655440002', 'tech_en_us', 'Technology', 'en', 'US', '{27,28}', true, NOW(), NOW()),
--     ('550e8400-e29b-41d4-a716-446655440003', 'tech_ko_kr', '기술', 'ko', 'KR', '{27,28}', true, NOW(), NOW()),
--     ('550e8400-e29b-41d4-a716-446655440004', 'tech_zh_cn', '科技', 'zh', 'CN', '{27,28}', true, NOW(), NOW()),
--     ('550e8400-e29b-41d4-a716-446655440005', 'tech_de_de', 'Technologie', 'de', 'DE', '{27,28}', true, NOW(), NOW())
-- ON CONFLICT (code) DO UPDATE SET
--     name = EXCLUDED.name,
--     language = EXCLUDED.language,
--     region_code = EXCLUDED.region_code,
--     category_ids = EXCLUDED.category_ids,
--     enabled = EXCLUDED.enabled,
--     updated_at = NOW();

-- -- Gaming genres
-- INSERT INTO genres (id, code, name, language, region_code, category_ids, enabled, created_at, updated_at) VALUES
--     ('550e8400-e29b-41d4-a716-446655440011', 'gaming_en_us', 'Gaming', 'en', 'US', '{20}', true, NOW(), NOW()),
--     ('550e8400-e29b-41d4-a716-446655440012', 'gaming_ja_jp', 'ゲーム', 'ja', 'JP', '{20}', true, NOW(), NOW())
-- ON CONFLICT (code) DO UPDATE SET
--     name = EXCLUDED.name,
--     language = EXCLUDED.language,
--     region_code = EXCLUDED.region_code,
--     category_ids = EXCLUDED.category_ids,
--     enabled = EXCLUDED.enabled,
--     updated_at = NOW();