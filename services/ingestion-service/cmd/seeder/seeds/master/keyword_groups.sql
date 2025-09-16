-- Keyword Groups for Japanese Engineering genre
-- Groups organize related keywords for filtering

-- Japanese Engineering keyword groups (genre_id: 550e8400-e29b-41d4-a716-446655440001)
INSERT INTO keyword_groups (id, genre_id, name, filter_type, target_field, description, enabled, created_at, updated_at) VALUES

-- Programming Languages
('550e8400-e29b-41d4-a716-446655440101', '550e8400-e29b-41d4-a716-446655440001', 'JavaScript/JS', 'include', 'title', 'JavaScript and related frameworks', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440102', '550e8400-e29b-41d4-a716-446655440001', 'Python', 'include', 'title', 'Python programming and libraries', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440103', '550e8400-e29b-41d4-a716-446655440001', 'Go/Golang', 'include', 'title', 'Go programming language', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440104', '550e8400-e29b-41d4-a716-446655440001', 'Rust', 'include', 'title', 'Rust programming language', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440105', '550e8400-e29b-41d4-a716-446655440001', 'Java', 'include', 'title', 'Java and JVM languages', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440106', '550e8400-e29b-41d4-a716-446655440001', 'C/C++', 'include', 'title', 'C and C++ programming', true, NOW(), NOW()),

-- General Programming Terms
('550e8400-e29b-41d4-a716-446655440201', '550e8400-e29b-41d4-a716-446655440001', 'プログラミング全般', 'include', 'title', '一般的なプログラミング用語', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440202', '550e8400-e29b-41d4-a716-446655440001', 'エンジニア/開発者', 'include', 'title', 'エンジニア・開発者関連用語', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440203', '550e8400-e29b-41d4-a716-446655440001', 'Web開発', 'include', 'title', 'Web開発関連用語', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440204', '550e8400-e29b-41d4-a716-446655440001', 'アプリ開発', 'include', 'title', 'モバイルアプリ開発', true, NOW(), NOW()),

-- AI/ML
('550e8400-e29b-41d4-a716-446655440301', '550e8400-e29b-41d4-a716-446655440001', 'AI/機械学習', 'include', 'title', 'AI・機械学習関連', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440302', '550e8400-e29b-41d4-a716-446655440001', 'データサイエンス', 'include', 'title', 'データサイエンス関連', true, NOW(), NOW()),

-- DevOps/Infrastructure
('550e8400-e29b-41d4-a716-446655440401', '550e8400-e29b-41d4-a716-446655440001', 'DevOps/インフラ', 'include', 'title', 'DevOps・インフラ関連', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440402', '550e8400-e29b-41d4-a716-446655440001', 'Docker/K8s', 'include', 'title', 'コンテナ技術', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440403', '550e8400-e29b-41d4-a716-446655440001', 'クラウド', 'include', 'title', 'クラウドプラットフォーム', true, NOW(), NOW()),

-- Educational
('550e8400-e29b-41d4-a716-446655440501', '550e8400-e29b-41d4-a716-446655440001', '入門/初心者向け', 'include', 'title', '入門・初心者向けコンテンツ', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440502', '550e8400-e29b-41d4-a716-446655440001', '解説/説明', 'include', 'title', '解説・教育コンテンツ', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440503', '550e8400-e29b-41d4-a716-446655440001', '実践/ハンズオン', 'include', 'title', '実践的コンテンツ', true, NOW(), NOW()),

-- Technical Terms
('550e8400-e29b-41d4-a716-446655440601', '550e8400-e29b-41d4-a716-446655440001', 'データベース', 'include', 'title', 'データベース技術', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440602', '550e8400-e29b-41d4-a716-446655440001', 'セキュリティ', 'include', 'title', 'セキュリティ関連', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440603', '550e8400-e29b-41d4-a716-446655440001', 'アーキテクチャ', 'include', 'title', 'ソフトウェアアーキテクチャ', true, NOW(), NOW())

ON CONFLICT (id) DO UPDATE SET
    genre_id = EXCLUDED.genre_id,
    name = EXCLUDED.name,
    filter_type = EXCLUDED.filter_type,
    target_field = EXCLUDED.target_field,
    description = EXCLUDED.description,
    enabled = EXCLUDED.enabled,
    updated_at = NOW();