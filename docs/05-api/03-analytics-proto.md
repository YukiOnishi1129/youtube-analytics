# Analytics API (analytics/v1/analytics.proto)

```protobuf
syntax = "proto3";

package analytics.v1;

option go_package = "github.com/yourorg/youtube-analytics/proto/analytics/v1;analyticsv1";

// ========= Common =========
enum Checkpoint {
  CHECKPOINT_UNSPECIFIED = 0;
  CHECKPOINT_3H   = 3;
  CHECKPOINT_6H   = 6;
  CHECKPOINT_12H  = 12;
  CHECKPOINT_24H  = 24;
  CHECKPOINT_48H  = 48;
  CHECKPOINT_72H  = 72;
  CHECKPOINT_168H = 168; // 7d
}

enum RankingKind {
  RANKING_KIND_UNSPECIFIED = 0;
  SPEED_VIEWS     = 1;  // view_growth_rate_per_hour
  SPEED_LIKES     = 2;  // like_growth_rate_per_hour
  RELATIVE_VIEWS  = 3;  // views_per_subscription_rate
  QUALITY         = 4;  // wilson_like_rate_lower_bound
  HEAT            = 5;  // likes_per_subscription_shrunk_rate
}

// ========= Messages =========
message Video {
  string video_id    = 1;
  string title       = 2;
  string channel_id  = 3;
  string thumbnail_url = 4;
  string video_url     = 5;
  string published_at  = 6; // RFC3339
}

message SnapshotPoint {
  Checkpoint checkpoint_hour = 1;
  int64 views_count = 2;
  int64 likes_count = 3;
}

message MetricPoint {
  Checkpoint checkpoint_hour = 1;
  double view_growth_rate_per_hour = 2;
  double like_growth_rate_per_hour = 3;
}

// ========= Ranking =========
message ListRankingRequest {
  string published_from = 1;      // RFC3339 (UTC)
  string published_to   = 2;      // RFC3339 (UTC, exclusive)
  Checkpoint checkpoint_hour = 3; // default: 24h (specified by client)
  RankingKind ranking_kind = 4;   // required: which metric to sort by
  bool hide_low_sample = 5;       // default true
  int32 category = 6;             // optional: YouTube category id
  int32 limit = 7;
  int32 offset = 8;
}

message RankingItem {
  Video video = 1;
  Checkpoint checkpoint_hour = 2;
  double main_metric = 3;
  int64 views_count = 4;
  int64 likes_count = 5;
}

message ListRankingResponse {
  repeated RankingItem items = 1;
}

// ========= Channel Ranking =========
message ListChannelRankingRequest {
  string channel_id = 1;
  string published_from = 2;
  string published_to   = 3;
  Checkpoint checkpoint_hour = 4;
  RankingKind ranking_kind = 5;
  bool hide_low_sample = 6;
  int32 limit = 7;
  int32 offset = 8;
}

message ListChannelRankingResponse {
  repeated RankingItem items = 1;
}

// ========= Video Detail =========
message GetVideoDetailRequest {
  string video_id = 1;
}

message GetVideoDetailResponse {
  Video video = 1;
  repeated SnapshotPoint snapshots = 2;
  repeated MetricPoint   metrics   = 3;
}

// ========= History =========
message ListHistoryRequest {
  string from = 1;              // RFC3339 date or datetime
  string to   = 2;              // RFC3339 (exclusive)
  RankingKind ranking_kind = 3; // optional
  Checkpoint  checkpoint_hour = 4; // optional
  int32 limit = 5;
  int32 offset = 6;
}

message History {
  string snapshot_id = 1;
  string snapshot_at = 2;       // RFC3339
  RankingKind ranking_kind = 3;
  Checkpoint checkpoint_hour = 4;
  string published_from = 5;
  string published_to   = 6;
  int32  top_n          = 7;
}

message ListHistoryResponse {
  repeated History items = 1;
}

message GetHistoryItemsRequest {
  string snapshot_id = 1;
}

message HistoryItem {
  int32 rank = 1;
  RankingItem ranking = 2;      // video + mainMetric + counts + CP
}

message GetHistoryItemsResponse {
  repeated HistoryItem items = 1;
}

// ========= Service =========
service AnalyticsService {
  rpc ListRanking         (ListRankingRequest)        returns (ListRankingResponse);
  rpc ListChannelRanking  (ListChannelRankingRequest) returns (ListChannelRankingResponse);
  rpc GetVideoDetail      (GetVideoDetailRequest)     returns (GetVideoDetailResponse);
  rpc ListHistory         (ListHistoryRequest)        returns (ListHistoryResponse);
  rpc GetHistoryItems     (GetHistoryItemsRequest)    returns (GetHistoryItemsResponse);
}
```
