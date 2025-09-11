import "@typespec/http";
import "@typespec/rest";

using TypeSpec.Http;
using TypeSpec.Rest;

@service({
  title: "Ingestion Service Cloud Tasks API",
})
@server("https://ingestion.youtube-analytics.example.com", "Production server")
namespace IngestionService.Tasks;

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
  ): {
    @statusCode statusCode: 200;
  } | {
    @statusCode statusCode: 400;
    @body error: Error;
  } | {
    @statusCode statusCode: 500;
    @body error: Error;
  };
}

model CreateSnapshotRequest {
  videoId: string;
  checkpointHour: 0 | 3 | 6 | 12 | 24 | 48 | 72 | 168;
}

@error
model Error {
  code: string;
  message: string;
}