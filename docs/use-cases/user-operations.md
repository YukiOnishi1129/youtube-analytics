# User Operations (Functions Called from UI)

## View Rankings

- **Actor**: User
- **Input**: Search criteria (publication date range, checkpoint time, ranking type)
- **Processing**: Filter ranking data by search criteria and sort by specified metrics
- **Output**: Ranking list (thumbnail, title, views, likes, main metric)

### Detailed Flow

1. User sets search criteria
   - Publication date range: From/To (UTC)
   - Checkpoint: Select from 3h/6h/12h/24h/48h/72h/7d
   - Ranking type:
     - `speed_views`: Initial velocity of view count
     - `speed_likes`: Initial velocity of like count
     - `relative_views`: Engagement level (views per subscriber)
     - `quality`: Quality (Wilson lower bound)
     - `heat`: Heat level (LPS)
   - Low-sample exclusion: ON/OFF
   - Category (optional): Education(27), Science & Technology(28), etc.

2. System extracts relevant data
   - Search records from video_metrics_checkpoint table matching conditions
   - Sort by specified metric
   - Pagination processing (limit/offset)

3. Display ranking results
   - Each video: thumbnail, title, channel name
   - Metrics: main metric value, view count, like count
   - Rank number

## View Video Details

- **Actor**: User
- **Input**: Video ID
- **Processing**:
  1. Retrieve snapshot list for that video (0h,3h,6h,12h...)
  2. Display progression in graph format
  3. Display growth rates for each checkpoint in table format
- **Output**: Detailed page showing video growth progression at a glance

### Detailed Flow

1. Retrieve video basic information
   - Title, channel, publication date/time
   - Thumbnail, video URL
   - YouTube category

2. Retrieve snapshot data
   - Data for all checkpoints from video_snapshots
   - View count, like count, subscriber count at each time point

3. Retrieve metrics data
   - Pre-calculated metrics from video_metrics_checkpoint
   - Growth rate, engagement level, quality, heat level

4. Graph rendering
   - Time series progression graph (view count, like count)
   - Growth rate graph (hourly increase)

5. Detailed display in table format
   - Measured values for each checkpoint
   - Growth rate from 0-hour point
   - Various metric values

## Manage Keywords

- **Actor**: User (Administrator)
- **Input**: Keyword name, type (include or exclude), description
- **Processing**:
  - System generates and saves regular expression patterns
  - Can toggle enable/disable or delete
- **Output**: Updated keyword list

### Detailed Flow

1. Display keyword list
   - name (display name)
   - filterType (include/exclude)
   - enabled (enabled/disabled)
   - description (description)
   - pattern (for advanced users: regular expression)

2. Add keyword
   - Input name (example: "React", "Next.js", "Docker")
   - Select filterType (include: adopt, exclude: exclude)
   - Input description (optional)
   - PatternBuilder automatically generates pattern on save
     - Absorbs notation variations (case, full-width/half-width)
     - Considers word boundaries

3. Edit keyword
   - Change attributes of existing keywords
   - Toggle enable/disable
   - Pattern is automatically regenerated

4. Delete keyword
   - Soft delete (set deleted_at)
   - No impact on past filtering results

5. Test feature (optional)
   - Input title string
   - Visualize which keywords match
   - Confirm filtering results (include/exclude)

## Authentication and Profile

- **Actor**: User
- **Input**: Google account or email + password
- **Processing**: Create/update profile on successful login (name, photo, email verification status)
- **Output**: Login session, profile information

### Detailed Flow

1. Login
   - Authentication via Identity Platform
   - Provider selection: Google, email + password
   - ID token acquisition

2. Profile creation/update
   - First login: create account
   - Existing user: update last login time
   - Profile information sync:
     - displayName (display name)
     - photoUrl (profile image)
     - email (email address)
     - emailVerified (email verification status)

3. Session management
   - Session maintenance via JWT token
   - Token refresh mechanism
   - Logout processing

4. Profile display
   - Information confirmation on My Page
   - List of linked providers
   - Role display (admin/editor/viewer)

## View History

- **Actor**: User
- **Input**: Date or weekly unit selection
- **Processing**: Read saved "ranking top N at that time"
- **Output**: Historical ranking list (CSV download also available)

### Detailed Flow

1. Display history list
   - Period filter (From/To)
   - Ranking type filter
   - Checkpoint filter
   - Each history entry:
     - snapshot_id
     - snapshot_at (save date/time)
     - ranking_kind
     - checkpoint_hour
     - published_from/to
     - top_n (number of saved items)

2. Display history details
   - Retrieve details by specifying snapshot_id
   - Ranking order at that time
   - Information and metrics for each video
   - Complete reproduction of state at that time

3. CSV output
   - Download ranking history in CSV format
   - Included information:
     - Rank
     - Video title
     - Channel name
     - Publication date/time
     - Main metric value
     - View count
     - Like count
     - Various metrics

4. History utilization
   - Past trend analysis
   - Weekly/monthly report creation
   - Discovery of success patterns