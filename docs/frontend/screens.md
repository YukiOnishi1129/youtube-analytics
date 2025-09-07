# Screen Specifications (OOUI Approach)

## Common (All Screens)

- **Header**: Service Name / Navigation (Rankings, History, Channels, Keywords, My Page) / Profile Icon
- **Notifications**: Toasts (Success/Warning/Error)
- **Loading**: List = Skeleton, Details = Skeleton Card
- **Empty State**: Guidance to relax conditions (Reset CP to 24h / Show low-sample)
- **Error**: Display current search conditions + Retry button
- **Responsive**: ≥1024px = Table/2 columns, <1024px = Vertical card layout

## P1. Ranking List /ranking

**Purpose**: View "current" rankings by Published range × CP × metric

- **Search Header**
  - Published Range (From/To, upper limit is exclusive)
  - Checkpoint: 3/6/12/24/48/72/7d (default=24h)
  - Tabs (RankingKind): Initial Speed (Views) / Initial Speed (Likes) / Penetration / Quality (Wilson) / Heat (LPS)
  - Category (All/27/28) (Optional)
  - Hide low-sample (ON by default)
- **List (Row)**
  - Thumbnail + Title (to video details) / Channel Name
  - Published Date / CP Chip (e.g., 24h)
  - **Main Metric (Bold)** + Auxiliary (Views@X / Likes@X)
- **Actions**
  - "Save as History" (Save current TopN to History)
  - Pagination (limit/offset)

## P2. Video Details /videos/:videoId

**Purpose**: Understand video progression and growth rate by CP

- **Header**: Thumbnail, Title, Channel Name, Published Date
- **Graph**: Views / Likes line chart (0/3/6/12/24/48/72/7d)
- **Table**: Growth rates for each CP 0→X (Views/hr, Likes/hr) / Views/Sub, Like rate, Wilson, LPS
- 24→72 interval rates shown only here (not in list)
- **Actions**: CSV export for this video (snapshot time series)

## P3. History

### P3-1. History List /history

**Purpose**: Find saved rankings (TopN at that time)

- **Filters**: Date/Week / RankingKind / CP
- **Cards**: snapshotAt / kind / CP / Published Range / TopN

### P3-2. History Details /history/:snapshotId

**Purpose**: View TopN from that time / CSV export

- **Table**: rank / Thumbnail, Title / Published Date / Main Metric (frozen) / Views/Likes (frozen)
- **CSV Download**

## P4. Channel List /channels

**Purpose**: Manage and switch monitored channels

- **Filters**: Subscribed (ON/ALL/OFF), Sort (Latest Video/Subscribers/Title), Keyword Search
- **Columns**: Thumbnail | Title (to details) | Subscribed Toggle | Subscribers (latest) | Latest Video (published date)
- **Row Click**: To channel details

## P5. Channel Details /channels/:channelId

**Purpose**: Ranking/sort for videos from this channel only

- **Top Section**: Channel Name, Icon, Subscribed Toggle, Subscriber Sparkline (last 30 days)
- **Header**: Published Range / Checkpoint / RankingKind / Hide low-sample
- **Table**: Same structure as P1 (limited to this channel's videos)

## P6. Keyword Management /keywords

**Purpose**: Manage collection filters (include/exclude)

- **List**: name | filterType | enabled | description | (advanced) pattern
- **Add/Edit Dialog**:
  - Input: name (display name), filterType (include/exclude), description
  - On Save: Server auto-generates pattern (handles variations)
  - (Optional) Test field: Enter title to visualize hits

## P7. My Page /mypage

**Purpose**: Profile verification

- **Display**: displayName / photoUrl / email / emailVerified / linked providers
- **Actions**: Sign out (external)