-- YouTube Categories Seed Data
-- Based on official YouTube API categories

INSERT INTO youtube_categories (id, name, assignable) VALUES
    (1, 'Film & Animation', true),
    (2, 'Autos & Vehicles', true),
    (10, 'Music', true),
    (15, 'Pets & Animals', true),
    (17, 'Sports', true),
    (19, 'Travel & Events', true),
    (20, 'Gaming', true),
    (22, 'People & Blogs', true),
    (23, 'Comedy', true),
    (24, 'Entertainment', true),
    (25, 'News & Politics', true),
    (26, 'Howto & Style', true),
    (27, 'Education', true),
    (28, 'Science & Technology', true),
    (29, 'Nonprofits & Activism', true)
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    assignable = EXCLUDED.assignable;