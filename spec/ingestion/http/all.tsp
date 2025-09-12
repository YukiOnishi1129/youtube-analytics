import "@typespec/http";
import "@typespec/rest";

using TypeSpec.Http;
using TypeSpec.Rest;

@service(#{
  title: "Ingestion Service API",
})
@server("https://ingestion.youtube-analytics.example.com", "Production server")
namespace IngestionService;

@error
model Error {
  code: string;
  message: string;
}

// ==================== Admin API ====================
@route("/admin")
namespace Admin {
  @doc("Collect trending videos")
  @route("/collect-trending")
  @post
  op collectTrending(
    @header("X-CloudScheduler-JobName") jobName: string,
    @header("X-CloudScheduler-ScheduleTime") scheduleTime: string,
  ): CollectTrendingResponse | {@statusCode statusCode: 500; @body error: Error;};

  @doc("Collect subscribed channel videos")
  @route("/collect-subscriptions")
  @post
  op collectSubscriptions(
    @header("X-CloudScheduler-JobName") jobName: string,
    @header("X-CloudScheduler-ScheduleTime") scheduleTime: string,
  ): CollectSubscriptionsResponse | {@statusCode statusCode: 500; @body error: Error;};

  @doc("Schedule video snapshots")
  @route("/schedule-snapshots")
  @post
  op scheduleSnapshots(
    @header("X-CloudScheduler-JobName") jobName: string,
    @header("X-CloudScheduler-ScheduleTime") scheduleTime: string,
  ): ScheduleSnapshotsResponse | {@statusCode statusCode: 500; @body error: Error;};

  @doc("Update channel metadata")
  @route("/update-channels")
  @post
  op updateChannels(
    @header("X-CloudScheduler-JobName") jobName: string,
    @header("X-CloudScheduler-ScheduleTime") scheduleTime: string,
  ): UpdateChannelsResponse | {@statusCode statusCode: 500; @body error: Error;};
}

model CollectTrendingResponse {
  videosProcessed: int32;
  videosAdded: int32;
  duration: string;
}

model CollectSubscriptionsResponse {
  channelsProcessed: int32;
  videosProcessed: int32;
  videosAdded: int32;
  duration: string;
}

model ScheduleSnapshotsResponse {
  videosProcessed: int32;
  tasksScheduled: int32;
  duration: string;
}

model UpdateChannelsResponse {
  channelsProcessed: int32;
  channelsUpdated: int32;
  duration: string;
}

// ==================== WebSub API ====================
@route("/websub/youtube")
namespace WebSub {
  @doc("WebSub verification endpoint for YouTube")
  @route("/verify")
  @get
  op verify(
    @query("hub.mode") mode: "subscribe" | "unsubscribe",
    @query("hub.topic") topic: string,
    @query("hub.challenge") challenge: string,
    @query("hub.lease_seconds") leaseSeconds?: int32,
  ): {
    @statusCode statusCode: 200;
    @header contentType: "text/plain";
    @body body: string;
  } | {
    @statusCode statusCode: 404;
  };

  @doc("WebSub notification endpoint for YouTube")
  @route("/notify")  
  @post
  op notify(
    @header("X-Hub-Signature") signature: string,
    @header("Content-Type") contentType: "application/atom+xml",
    @body body: string,
  ): {
    @statusCode statusCode: 204;
  } | {
    @statusCode statusCode: 400;
    @body error: Error;
  };
}

// ==================== Tasks API ====================
@route("/tasks")
namespace Tasks {
  @doc("Create video snapshot task")
  @route("/snapshot")
  @post
  op createSnapshot(
    @header("X-CloudTasks-TaskName") taskName: string,
    @header("X-CloudTasks-QueueName") queueName: string,
    @header("X-CloudTasks-TaskRetryCount") retryCount: int32,
    @header("X-CloudTasks-TaskExecutionCount") executionCount: int32,
    @header("X-CloudTasks-TaskETA") eta: string,
    @body request: CreateSnapshotRequest,
  ): void | {@statusCode statusCode: 400; @body error: Error;} | {@statusCode statusCode: 500; @body error: Error;};
}

model CreateSnapshotRequest {
  videoId: string;
  checkpointHour: 0 | 3 | 6 | 12 | 24 | 48 | 72 | 168;
}