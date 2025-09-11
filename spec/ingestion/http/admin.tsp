import "@typespec/http";
import "@typespec/rest";

using TypeSpec.Http;
using TypeSpec.Rest;

@service({
  title: "Ingestion Service Admin API",
})
@server("https://ingestion.youtube-analytics.example.com", "Production server")
namespace IngestionService.Admin;

@route("/admin")
namespace Admin {
  @doc("Collect trending videos")
  @route("/collect-trending")
  @post
  op collectTrending(
    @header("X-CloudScheduler-JobName") jobName: string,
    @header("X-CloudScheduler-ScheduleTime") scheduleTime: string,
  ): {
    @statusCode statusCode: 200;
    @body response: CollectTrendingResponse;
  } | {
    @statusCode statusCode: 500;
    @body error: Error;
  };

  @doc("Collect subscribed channel videos")
  @route("/collect-subscriptions")
  @post
  op collectSubscriptions(
    @header("X-CloudScheduler-JobName") jobName: string,
    @header("X-CloudScheduler-ScheduleTime") scheduleTime: string,
  ): {
    @statusCode statusCode: 200;
    @body response: CollectSubscriptionsResponse;
  } | {
    @statusCode statusCode: 500;
    @body error: Error;
  };

  @doc("Schedule video snapshots")
  @route("/schedule-snapshots")
  @post
  op scheduleSnapshots(
    @header("X-CloudScheduler-JobName") jobName: string,
    @header("X-CloudScheduler-ScheduleTime") scheduleTime: string,
  ): {
    @statusCode statusCode: 200;
    @body response: ScheduleSnapshotsResponse;
  } | {
    @statusCode statusCode: 500;
    @body error: Error;
  };

  @doc("Update channel metadata")
  @route("/update-channels")
  @post
  op updateChannels(
    @header("X-CloudScheduler-JobName") jobName: string,
    @header("X-CloudScheduler-ScheduleTime") scheduleTime: string,
  ): {
    @statusCode statusCode: 200;
    @body response: UpdateChannelsResponse;
  } | {
    @statusCode statusCode: 500;
    @body error: Error;
  };
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

@error
model Error {
  code: string;
  message: string;
}