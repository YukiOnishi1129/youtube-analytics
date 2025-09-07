# Derived Metrics

## Basic Metric Calculations

- `rel_views = views / greatest(subscriber_count, 1)`
- `like_rate = likes / greatest(views, 1)`
- `momentum = 0.6*Δviews(6-24h)/18 + 0.4*Δviews(0-6h)/6`
- Z-score normalization over last 7 days → `param_score = 0.4*z_momentum + 0.35*z_rel_views + 0.25*z_quality`

## Ranking Metrics

### 1. Initial Speed (Views) - speed_views
`view_growth_rate_per_hour = (viewsX - views0) / X[h]`

### 2. Initial Speed (Likes) - speed_likes
`like_growth_rate_per_hour = (likesX - likes0) / X[h]`

### 3. Penetration Rate - relative_views
`views_per_subscription_rate = viewsX / subsX`

### 4. Quality (Wilson Lower Bound) - quality
Calculate the Wilson confidence interval lower bound for like rate

### 5. Heat (LPS) - heat
`likes_per_subscription_shrunk_rate = SCALE * likesX / (subsX + OFFSET)`