import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi";

using TypeSpec.Http;
using TypeSpec.Rest;
using TypeSpec.OpenAPI;

namespace IngestionService;

// Channel model
model Channel {
  id: string;
  youtubeChannelId: string;
  title: string;
  thumbnailUrl: string;
  subscribed: boolean;
  createdAt: utcDateTime;
  updatedAt?: utcDateTime;
}

// Video model
model Video {
  id: string;
  youtubeVideoId: string;
  channelId: string;
  title: string;
  description?: string;
  thumbnailUrl: string;
  publishedAt: utcDateTime;
  viewCount: int32;
  likeCount: int32;
  dislikeCount: int32;
  commentCount: int32;
  categoryId: string;
  tags: string[];
  durationSeconds: int32;
  privacyStatus: string;
  isDeleted: boolean;
  createdAt: utcDateTime;
  updatedAt?: utcDateTime;
}

// VideoSnapshot model
model VideoSnapshot {
  id: string;
  videoId: string;
  viewCount: int32;
  likeCount: int32;
  dislikeCount: int32;
  commentCount: int32;
  snapshotAt: utcDateTime;
  createdAt: utcDateTime;
}

// Keyword model
model Keyword {
  id: string;
  name: string;
  filterType: string;
  pattern: string;
  description?: string;
  createdAt: utcDateTime;
  updatedAt?: utcDateTime;
}

// List responses
model ListChannelsResponse {
  channels: Channel[];
}

model ListVideosResponse {
  videos: Video[];
}

model ListSnapshotsResponse {
  snapshots: VideoSnapshot[];
}

model ListKeywordsResponse {
  keywords: Keyword[];
}

// Update result responses
model UpdateChannelsResponse {
  channelsProcessed: int32;
  channelsUpdated: int32;
  duration: string;
}

model CollectTrendingResponse {
  videosCollected: int32;
  videosUpdated: int32;
  videosCreated: int32;
  duration: string;
}

model CollectSubscriptionsResponse {
  channelsProcessed: int32;
  videosCollected: int32;
  videosCreated: int32;
  duration: string;
}

model ScheduleSnapshotsResponse {
  videosProcessed: int32;
  tasksScheduled: int32;
  duration: string;
}