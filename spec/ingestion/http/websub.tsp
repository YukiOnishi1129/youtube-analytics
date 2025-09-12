import "@typespec/http";
import "@typespec/rest";

using TypeSpec.Http;
using TypeSpec.Rest;

@service(#{
  title: "Ingestion Service WebSub API",
})
@server("https://ingestion.youtube-analytics.example.com", "Production server")
namespace IngestionService.WebSub;

@route("/websub/youtube")
namespace YouTube {
  @doc("WebSub verification endpoint for YouTube")
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

@error
model Error {
  code: string;
  message: string;
}