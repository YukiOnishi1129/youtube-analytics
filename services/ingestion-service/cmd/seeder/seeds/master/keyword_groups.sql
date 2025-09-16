-- Keyword Groups for Japanese Engineering genre
-- Groups organize related keywords for filtering

-- Japanese Engineering keyword groups (genre_id: 550e8400-e29b-41d4-a716-446655440001)
INSERT INTO keyword_groups (id, genre_id, name, filter_type, target_field, description, enabled, created_at, updated_at) VALUES

-- Programming Languages
('550e8400-e29b-41d4-a716-446655440101', '550e8400-e29b-41d4-a716-446655440001', 'JavaScript/TypeScript', 'include', 'title', 'JavaScript, TypeScript and related technologies', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440103', '550e8400-e29b-41d4-a716-446655440001', 'Go/Golang', 'include', 'title', 'Go programming language', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440105', '550e8400-e29b-41d4-a716-446655440001', 'Java', 'include', 'title', 'Java programming language', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440106', '550e8400-e29b-41d4-a716-446655440001', 'Ruby/Rails', 'include', 'title', 'Ruby and Ruby on Rails', true, NOW(), NOW()),

-- Career and Job Related
('550e8400-e29b-41d4-a716-446655440202', '550e8400-e29b-41d4-a716-446655440001', 'エンジニア職種', 'include', 'title', 'エンジニア・SIer・SES関連', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440203', '550e8400-e29b-41d4-a716-446655440001', 'キャリア/報酬', 'include', 'title', 'キャリア・年収・単価関連', true, NOW(), NOW()),

-- AI/ML
('550e8400-e29b-41d4-a716-446655440301', '550e8400-e29b-41d4-a716-446655440001', 'AI', 'include', 'title', 'AI・人工知能関連', true, NOW(), NOW()),

-- Cloud
('550e8400-e29b-41d4-a716-446655440403', '550e8400-e29b-41d4-a716-446655440001', 'クラウドプラットフォーム', 'include', 'title', 'AWS・GCP等のクラウドサービス', true, NOW(), NOW())

ON CONFLICT (id) DO UPDATE SET
    genre_id = EXCLUDED.genre_id,
    name = EXCLUDED.name,
    filter_type = EXCLUDED.filter_type,
    target_field = EXCLUDED.target_field,
    description = EXCLUDED.description,
    enabled = EXCLUDED.enabled,
    updated_at = NOW();