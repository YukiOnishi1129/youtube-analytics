# Data Collection Use Cases

## Automatic Trending Video Collection

- **Actor**: System (automated batch)
- **Timing**: Several times per day at specified times (e.g., daily at 03:00 and 15:00)
- **Input**: Enabled genres from the system
- **Processing**: For each enabled genre:
  1. Get genre configuration (region, YouTube categories, keywords)
  2. Fetch trending videos from YouTube API for that region and categories
  3. Match each video title against genre-specific keyword patterns
  4. Discard if exclusion rules match
  5. Register as "monitored video" if inclusion rules match
  6. Associate video with the matching genre (videos can belong to multiple genres)
- **Output**: Monitored videos are added to the system with their genre associations (M:N relationship)
- **Examples**:
  - Genre "Engineering (JP)": Fetch JP trending from categories 27,28, filter by Japanese tech keywords
  - Genre "Engineering (EN)": Fetch US trending from categories 27,28, filter by English tech keywords
  - Genre "English Learning (JP)": Fetch JP trending from category 27, filter by English learning keywords

## New Video Monitoring (WebSub)

- **Actor**: YouTube (notification), System (reception)
- **Timing**: When monitored channels publish new videos
- **Input**: Notification (video ID)
- **Processing**:
  1. Fetch initial video data (views, likes, etc.) from YouTube API
  2. Save as "0-hour point data"
  3. Schedule rechecking tasks for "3 hours, 6 hours, 12 hours, 24 hours, 48 hours, 72 hours, 7 days later"
- **Output**: Video enters monitoring list with initial data and scheduled follow-up checks

## Periodic Snapshots

- **Actor**: System (task execution)
- **Timing**: At each checkpoint after video publication: 3h/6h/12h/24h/48h/72h/7d
- **Input**: Video ID, target checkpoint time
- **Processing**:
  1. Fetch latest views and likes from YouTube API
  2. Save as checkpoint-specific snapshot
- **Output**: Time-series data accumulates

## Growth Rate and Metrics Calculation

- **Actor**: System
- **Timing**: Immediately after snapshot is saved
- **Input**: Snapshot data (0-hour point for comparison)
- **Processing**:
  1. View count growth rate from 0→X hours
  2. Like count growth rate from 0→X hours
  3. Views per subscriber
  4. Like rate (likes/views) and its confidence lower bound (Wilson score)
  5. Likes per subscriber (shrunk version)
  6. Flag videos with insufficient data as excluded from rankings
- **Output**: Metrics usable for rankings are organized
