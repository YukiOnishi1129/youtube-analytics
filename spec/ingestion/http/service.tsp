import "@typespec/http";
import "@typespec/rest";
import "./websub.tsp";
import "./tasks.tsp";
import "./admin.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

@service(#{
  title: "Ingestion Service API",
})
@server("https://ingestion.youtube-analytics.example.com", "Production server")
namespace IngestionService;