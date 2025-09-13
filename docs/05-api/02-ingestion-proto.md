# Ingestion API (ingestion/v1/ingestion.proto)

```protobuf
syntax = "proto3";

package ingestion.v1;

option go_package = "github.com/yourorg/youtube-analytics/proto/ingestion/v1;ingestionv1";

// ===== Keyword =====
enum FilterType {
  FILTER_TYPE_UNSPECIFIED = 0;
  INCLUDE = 1;
  EXCLUDE = 2;
}

message Keyword {
  string id = 1;             // uuid v7
  string name = 2;           // display name
  FilterType filter_type = 3;
  bool   enabled = 4;
  string description = 5;
  string pattern = 6;        // shown only in advanced mode
}

message ListKeywordsRequest {}
message ListKeywordsResponse { 
  repeated Keyword items = 1; 
}

message CreateKeywordRequest { 
  string name = 1; 
  FilterType filter_type = 2; 
  string description = 3; 
}

message UpdateKeywordRequest { 
  string id = 1; 
  string name = 2; 
  FilterType filter_type = 3; 
  bool enabled = 4; 
  string description = 5; 
}

message DeleteKeywordRequest { 
  string id = 1; 
}

// ===== Channel =====
message ListChannelsRequest {
  string subscribed = 1;     // "on"|"off"|"all"
  string sort = 2;           // "latest_video"|"subs"|"title"
  string q = 3;              // prefix match
  int32  limit = 4;
  int32  offset = 5;
}

message Channel {
  string channel_id = 1;
  string title = 2;
  string thumbnail_url = 3;
  bool   subscribed = 4;
  int64  subscriber_count = 5;           // latest (optional)
  string latest_video_published_at = 6;  // RFC3339
}

message ListChannelsResponse { 
  repeated Channel items = 1; 
}

message SetChannelSubscriptionRequest {
  string channel_id = 1;
  bool   subscribed = 2;
}

// ===== Snapshot (internal) =====
message SnapshotRequest {
  string video_id = 1;
  int32  checkpoint_hour = 2;
}

service IngestionService {
  // Keywords
  rpc ListKeywords (ListKeywordsRequest) returns (ListKeywordsResponse);
  rpc CreateKeyword (CreateKeywordRequest) returns (Keyword);
  rpc UpdateKeyword (UpdateKeywordRequest) returns (Keyword);
  rpc DeleteKeyword (DeleteKeywordRequest) returns (.google.protobuf.Empty);

  // Channels
  rpc ListChannels (ListChannelsRequest) returns (ListChannelsResponse);
  rpc SetChannelSubscription (SetChannelSubscriptionRequest) returns (.google.protobuf.Empty);

  // Snapshots (internal)
  rpc InsertSnapshot (SnapshotRequest) returns (.google.protobuf.Empty);
}
```
