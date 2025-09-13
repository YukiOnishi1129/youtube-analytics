package main

import (
	"fmt"
	"log"
	"os"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/adapter/gateway/cloudtasks"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/adapter/gateway/mock"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/adapter/gateway/postgres"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/adapter/gateway/youtube"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/config"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/driver/datastore"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/driver/transport"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := datastore.OpenPostgres(cfg.Database.URL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	repo := postgres.NewRepository(db)
	channelRepo := postgres.NewChannelRepository(repo)
	videoRepo := postgres.NewVideoRepository(repo)
	channelSnapshotRepo := postgres.NewChannelSnapshotRepository(repo)
	videoSnapshotRepo := postgres.NewVideoSnapshotRepository(repo)

	// Use mock keyword repository for now until SQL queries are generated
	keywordRepo := mock.NewKeywordRepository()

	// Initialize YouTube client
	youtubeClient, err := youtube.NewClient(cfg.YouTube.APIKey)
	if err != nil {
		log.Fatalf("Failed to create YouTube client: %v", err)
	}

	// Initialize task scheduler
	taskScheduler, err := cloudtasks.NewTaskScheduler(
		cfg.GCP.ProjectID,
		cfg.CloudTasks.Location,
		cfg.CloudTasks.QueueName,
		"", // Base URL will be set later
	)
	if err != nil {
		log.Fatalf("Failed to create task scheduler: %v", err)
	}

	// Initialize event publisher (using mock for now)
	eventPublisher := mock.NewEventPublisher()

	// Determine address
	addr := fmt.Sprintf(":%d", cfg.Server.GRPCPort)
	if envAddr := os.Getenv("GRPC_ADDR"); envAddr != "" {
		addr = envAddr
	}

	// Bootstrap and start gRPC server
	if err := transport.BootstrapGRPCWithKeyword(
		addr,
		channelRepo,
		channelSnapshotRepo,
		videoRepo,
		videoSnapshotRepo,
		keywordRepo,
		youtubeClient,
		taskScheduler,
		eventPublisher,
	); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
