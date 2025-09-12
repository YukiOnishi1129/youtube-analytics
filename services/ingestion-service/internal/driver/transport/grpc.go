package transport

import (
	"fmt"
	"log"
	"net"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/service"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/driver/grpc"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/gateway"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/usecase"
	// pb "github.com/YukiOnishi1129/youtube-analytics/pkg/proto/ingestion/v1"
	googlegrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// BootstrapGRPC wires everything and starts gRPC server
func BootstrapGRPC(
	addr string,
	channelRepo gateway.ChannelRepository,
	channelSnapshotRepo gateway.ChannelSnapshotRepository,
	videoRepo gateway.VideoRepository,
	videoSnapshotRepo gateway.VideoSnapshotRepository,
	keywordRepo gateway.KeywordRepository,
	youtubeClient gateway.YouTubeClient,
	taskScheduler gateway.TaskScheduler,
	eventPublisher gateway.EventPublisher,
) error {
	// Initialize use cases
	channelUseCase := usecase.NewChannelUseCase(
		channelRepo,
		youtubeClient,
	)

	videoUseCase := usecase.NewVideoUseCase(
		videoRepo,
		channelRepo,
		youtubeClient,
		eventPublisher,
	)

	// Create snapshot scheduler
	snapshotScheduler := service.NewSnapshotScheduler()

	systemUseCase := usecase.NewSystemUseCase(
		videoRepo,
		videoSnapshotRepo,
		taskScheduler,
		snapshotScheduler,
		youtubeClient,
	)

	// Create gRPC server handler
	_ = grpc.NewServer(
		channelUseCase,
		videoUseCase,
		systemUseCase,
	)

	// Create gRPC server with options
	grpcServer := googlegrpc.NewServer(
		// Add interceptors here if needed
		// googlegrpc.UnaryInterceptor(someInterceptor),
	)

	// Register service
	// pb.RegisterIngestionServiceServer(grpcServer, handler)

	// Register reflection service on gRPC server for debugging
	reflection.Register(grpcServer)

	// Start listening
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	log.Printf("ingestion-service gRPC listening on %s", addr)
	return grpcServer.Serve(lis)
}

// BootstrapGRPCWithKeyword includes keyword use case
func BootstrapGRPCWithKeyword(
	addr string,
	channelRepo gateway.ChannelRepository,
	channelSnapshotRepo gateway.ChannelSnapshotRepository,
	videoRepo gateway.VideoRepository,
	videoSnapshotRepo gateway.VideoSnapshotRepository,
	keywordRepo gateway.KeywordRepository,
	youtubeClient gateway.YouTubeClient,
	taskScheduler gateway.TaskScheduler,
	eventPublisher gateway.EventPublisher,
) error {
	// Initialize use cases
	channelUseCase := usecase.NewChannelUseCase(
		channelRepo,
		youtubeClient,
	)

	videoUseCase := usecase.NewVideoUseCase(
		videoRepo,
		channelRepo,
		youtubeClient,
		eventPublisher,
	)

	// Create snapshot scheduler
	snapshotScheduler := service.NewSnapshotScheduler()

	systemUseCase := usecase.NewSystemUseCase(
		videoRepo,
		videoSnapshotRepo,
		taskScheduler,
		snapshotScheduler,
		youtubeClient,
	)

	keywordUseCase := usecase.NewKeywordUseCase(
		keywordRepo,
	)

	// Create extended gRPC server handler with keyword support
	_ = grpc.NewServerWithKeyword(
		channelUseCase,
		videoUseCase,
		systemUseCase,
		keywordUseCase,
	)

	// Create gRPC server
	grpcServer := googlegrpc.NewServer()

	// Register service
	// pb.RegisterIngestionServiceServer(grpcServer, handler)

	// Register reflection service
	reflection.Register(grpcServer)

	// Start listening
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	log.Printf("ingestion-service gRPC listening on %s (with keyword support)", addr)
	return grpcServer.Serve(lis)
}

