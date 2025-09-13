-- Default keywords for video filtering (Japanese programming content)
INSERT INTO keywords (id, name, filter_type, pattern, description, created_at) VALUES
-- Programming languages with Japanese terms
('550e8400-e29b-41d4-a716-446655440001', 'JavaScript', 'TITLE', 'javascript|js|node.js|nodejs|react|vue|angular|ジャバスクリプト|リアクト|ビュー', 'JavaScript and related frameworks', NOW()),
('550e8400-e29b-41d4-a716-446655440002', 'Python', 'TITLE', 'python|django|flask|pandas|numpy|machine learning|ML|AI|パイソン|機械学習|人工知能', 'Python programming and data science', NOW()),
('550e8400-e29b-41d4-a716-446655440003', 'Go/Golang', 'TITLE', 'golang|go programming|go lang|ゴー言語|Go言語', 'Go programming language', NOW()),
('550e8400-e29b-41d4-a716-446655440004', 'Rust', 'TITLE', 'rust programming|rust lang|ラスト|Rust言語', 'Rust programming language', NOW()),
('550e8400-e29b-41d4-a716-446655440005', 'TypeScript', 'TITLE', 'typescript|ts|タイプスクリプト', 'TypeScript programming language', NOW()),

-- Japanese programming terms
('550e8400-e29b-41d4-a716-446655440006', 'プログラミング', 'TITLE', 'プログラミング|プログラム|programming|開発|コーディング|coding', 'General programming terms in Japanese', NOW()),
('550e8400-e29b-41d4-a716-446655440007', 'エンジニア', 'TITLE', 'エンジニア|engineer|開発者|developer|プログラマー|programmer', 'Engineer/developer terms in Japanese', NOW()),
('550e8400-e29b-41d4-a716-446655440008', 'Web開発', 'TITLE', 'web開発|ウェブ開発|フロントエンド|frontend|バックエンド|backend|フルスタック', 'Web development terms in Japanese', NOW()),
('550e8400-e29b-41d4-a716-446655440009', 'アプリ開発', 'TITLE', 'アプリ開発|モバイル開発|iOS|Android|スマホアプリ|mobile', 'Mobile app development terms', NOW()),

-- Cloud platforms with Japanese terms
('550e8400-e29b-41d4-a716-446655440010', 'AWS', 'TITLE', 'aws|amazon web services|ec2|s3|lambda|アマゾンウェブサービス', 'Amazon Web Services', NOW()),
('550e8400-e29b-41d4-a716-446655440011', 'Google Cloud', 'TITLE', 'gcp|google cloud|firebase|bigquery|グーグルクラウド|ファイアベース', 'Google Cloud Platform', NOW()),
('550e8400-e29b-41d4-a716-446655440012', 'Azure', 'TITLE', 'azure|microsoft azure|アジュール|マイクロソフトアジュール', 'Microsoft Azure', NOW()),

-- DevOps tools with Japanese terms
('550e8400-e29b-41d4-a716-446655440020', 'Docker', 'TITLE', 'docker|container|dockerfile|ドッカー|コンテナ|仮想化', 'Docker containerization', NOW()),
('550e8400-e29b-41d4-a716-446655440021', 'Kubernetes', 'TITLE', 'kubernetes|k8s|kubectl|クバネティス|クーバネティス', 'Kubernetes orchestration', NOW()),
('550e8400-e29b-41d4-a716-446655440022', 'CI/CD', 'TITLE', 'ci/cd|continuous integration|continuous deployment|github actions|jenkins|継続的インテグレーション|自動化', 'CI/CD tools and practices', NOW()),

-- Japanese tech tutorial terms
('550e8400-e29b-41d4-a716-446655440030', '入門', 'TITLE', '入門|初心者|beginner|基礎|tutorial|チュートリアル', 'Beginner/tutorial content', NOW()),
('550e8400-e29b-41d4-a716-446655440031', '解説', 'TITLE', '解説|説明|紹介|レビュー|review', 'Explanation/review content', NOW()),
('550e8400-e29b-41d4-a716-446655440032', 'ハンズオン', 'TITLE', 'ハンズオン|hands-on|実践|実装|作ってみた', 'Hands-on/practical content', NOW())

ON CONFLICT (id) DO NOTHING;