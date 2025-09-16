package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/adapter/gateway/cloudtasks"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/adapter/gateway/postgres"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/adapter/gateway/youtube"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/service"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/driver/config"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/driver/datastore"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/usecase"
)

func main() {
	// Parse command line arguments
	var (
		hours  = flag.Int("hours", 24, "Schedule snapshots for videos published within the last N hours")
		dryRun = flag.Bool("dry-run", false, "Dry run mode - only log what would be done")
	)
	flag.Parse()

	// Load configuration
	cfg := config.Load()

	// Setup signal handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("Shutting down...")
		cancel()
	}()

	// Initialize database connection
	db, err := datastore.OpenPostgres(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	pgRepo := postgres.NewRepository(db)
	videoRepo := postgres.NewVideoRepository(pgRepo)
	snapshotRepo := postgres.NewVideoSnapshotRepository(pgRepo)

	// Initialize task scheduler
	taskScheduler, err := cloudtasks.NewTaskScheduler(
		cfg.CloudTasksProjectID,
		cfg.CloudTasksLocation,
		cfg.CloudTasksQueueName,
		cfg.CloudTasksServiceURL,
	)
	if err != nil {
		log.Fatalf("Failed to create task scheduler: %v", err)
	}

	// Initialize YouTube client
	youtubeClient, err := youtube.NewClient(cfg.YouTubeAPIKey)
	if err != nil {
		log.Fatalf("Failed to create YouTube client: %v", err)
	}
	
	// Initialize snapshot scheduler
	snapshotScheduler := service.NewSnapshotScheduler()
	
	// Initialize use case
	systemUseCase := usecase.NewSystemUseCase(
		videoRepo,
		snapshotRepo,
		taskScheduler,
		snapshotScheduler,
		youtubeClient,
	)

	// Log start
	log.Printf("Starting snapshot scheduling batch (hours=%d, dry-run=%v)", *hours, *dryRun)
	start := time.Now()

	// Execute scheduling
	result, err := systemUseCase.ScheduleSnapshots(ctx)
	if err != nil {
		log.Fatalf("Failed to schedule snapshots: %v", err)
	}

	// Log results
	log.Printf("Completed: videos=%d, tasks=%d, duration=%s",
		result.VideosProcessed, result.TasksScheduled, time.Since(start))
}