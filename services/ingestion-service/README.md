# Ingestion Service

This service handles the ingestion of YouTube video and channel data.

## Architecture

The service follows Clean Architecture principles with the following layers:

- **Domain**: Core business logic and entities
- **Port**: Interfaces for input (use cases) and output (gateways)
- **Adapter**: Implementation of gateways (database, YouTube API, etc.)
- **Driver**: Entry points (HTTP and gRPC servers)
- **UseCase**: Application business logic

## Prerequisites

- Go 1.21 or later
- PostgreSQL
- Protocol Buffers compiler (protoc)
- YouTube Data API key
- Google Cloud Project (for Cloud Tasks)

## Setup

### 1. Install Dependencies

```bash
go mod download
```

### 2. Install Proto Tools

```bash
make proto-install
```

### 3. Generate Proto Files

```bash
# From the project root
make proto
```

### 4. Set Environment Variables

Create a `.env` file:

```env
# Server
PORT=8080
GRPC_PORT=50051

# Database
DATABASE_URL=postgres://user:password@localhost:5432/youtube_analytics?sslmode=disable

# YouTube API
YOUTUBE_API_KEY=your-api-key

# Google Cloud
GCP_PROJECT_ID=your-project-id
CLOUD_TASKS_LOCATION=us-central1
CLOUD_TASKS_QUEUE_NAME=ingestion-tasks
```

### 5. Run Database Migrations

```bash
# Make sure your database is running
# Run migrations (implement your migration tool)
```

## Running the Service

### HTTP Server

```bash
make ingestion-http
# or
go run cmd/http/main.go
```

The HTTP server will start on port 8080 (or PORT env var).

### gRPC Server

```bash
make ingestion-grpc
# or
go run cmd/grpc/main.go
```

The gRPC server will start on port 50051 (or GRPC_PORT env var).

## API Endpoints

### HTTP API

- `POST /admin/collect-trending` - Collect trending videos
- `POST /admin/collect-subscriptions` - Collect videos from subscribed channels
- `POST /admin/schedule-snapshots` - Schedule video snapshots
- `POST /admin/update-channels` - Update channel metadata
- `POST /tasks/snapshot` - Create video snapshot (Cloud Tasks endpoint)
- `GET /websub/youtube/verify` - WebSub verification endpoint
- `POST /websub/youtube/notify` - WebSub notification endpoint

### gRPC API

See `proto/ingestion/v1/ingestion.proto` for the complete service definition.

Key services:
- Channel operations (GetChannel, ListChannels, etc.)
- Video operations (GetVideo, ListVideos, CollectTrending, etc.)
- Snapshot operations (CreateSnapshot, ListSnapshots, etc.)
- System operations (ScheduleSnapshots, UpdateChannels)

## Development

### Running Tests

```bash
make test
```

### Linting

```bash
make lint
```

### Generating OpenAPI Types

```bash
cd services/ingestion-service
oapi-codegen -generate types,gin -package generated -o internal/driver/http/generated/openapi.gen.go ../../spec/ingestion/http/dist/openapi.yaml
```

## Architecture Notes

### Use Case Pattern

Use cases now use input structs for better API design:

```go
type CreateKeywordInput struct {
    Name        string
    FilterType  string
    Pattern     string
    Description *string
}

// Usage
input := &CreateKeywordInput{
    Name: "Gaming",
    FilterType: "include",
    Pattern: "game|gaming",
}
result, err := keywordUseCase.CreateKeyword(ctx, input)
```

### Repository Pattern

Repositories implement both `GetByID` and `FindByID` methods:
- `GetByID`: Used by use cases for primary operations
- `FindByID`: Used internally by repositories

### Error Handling

Domain-specific errors are defined in the domain layer:
- `ErrChannelNotFound`
- `ErrVideoNotFound`
- `ErrInvalidFilterType`

## Troubleshooting

### Proto Generation Issues

If you encounter issues with proto generation:

1. Ensure `protoc` is installed
2. Run `make proto-install` to install Go plugins
3. Check that proto files are in the correct location
4. Verify import paths in proto files

### Database Connection Issues

1. Check DATABASE_URL format
2. Ensure PostgreSQL is running
3. Verify database exists
4. Check network connectivity

### YouTube API Issues

1. Verify API key is valid
2. Check API quotas in Google Cloud Console
3. Enable YouTube Data API v3 in your project