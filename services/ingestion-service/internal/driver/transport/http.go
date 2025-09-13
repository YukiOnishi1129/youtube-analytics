package transport

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	cloudtasksgw "github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/adapter/gateway/cloudtasks"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/adapter/gateway/postgres"
	pubsubgw "github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/adapter/gateway/pubsub"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/adapter/gateway/youtube"
	httpPresenter "github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/adapter/presenter/http"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/service"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/usecase"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/driver/datastore"
	httpdriver "github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/driver/http"
)

// BootstrapHTTP bootstraps HTTP server with all dependencies
func BootstrapHTTP(
	addr string,
	db *sql.DB,
	projectID string,
	youtubeAPIKey string,
	region string,
	taskQueue string,
	eventTopic string,
) error {
	// Initialize presenters
	channelPresenter := httpPresenter.NewChannelPresenter()
	videoPresenter := httpPresenter.NewVideoPresenter()
	systemPresenter := httpPresenter.NewSystemPresenter()
	keywordPresenter := httpPresenter.NewKeywordPresenter()

	// Initialize gateways
	repo := postgres.NewRepository(db)

	// Initialize repositories
	channelRepo := postgres.NewChannelRepository(repo)
	videoRepo := postgres.NewVideoRepository(repo)
	videoSnapshotRepo := postgres.NewVideoSnapshotRepository(repo)
	keywordRepo := postgres.NewKeywordRepository(repo)

	// Initialize external service clients
	youtubeClient, err := youtube.NewClient(youtubeAPIKey)
	if err != nil {
		return fmt.Errorf("failed to create YouTube client: %w", err)
	}

	// websub client will be used in future implementations
	// webSubClient := websub.NewHubClient()

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "https://ingestion-service.example.com" // Default for local dev
	}

	taskScheduler, err := cloudtasksgw.NewTaskScheduler(projectID, region, taskQueue, baseURL)
	if err != nil {
		return fmt.Errorf("failed to create task scheduler: %w", err)
	}

	eventPublisher, err := pubsubgw.NewEventPublisher(projectID)
	if err != nil {
		return fmt.Errorf("failed to create event publisher: %w", err)
	}

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

	// Create snapshot scheduler service
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

	// Create router with presenters
	router := httpdriver.SetupRouterWithPresenters(
		channelUseCase,
		videoUseCase,
		systemUseCase,
		keywordUseCase,
		channelPresenter,
		videoPresenter,
		systemPresenter,
		keywordPresenter,
	)

	// Create HTTP server
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("HTTP server starting on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down HTTP server...")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server forced to shutdown: %v", err)
	}

	log.Println("HTTP server exited")
	return nil
}

// BootstrapHTTPWithDefaults bootstraps HTTP server with default configuration
func BootstrapHTTPWithDefaults(addr string) error {
	// Get configuration from environment
	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		return fmt.Errorf("GCP_PROJECT_ID environment variable is required")
	}

	youtubeAPIKey := os.Getenv("YOUTUBE_API_KEY")
	if youtubeAPIKey == "" {
		return fmt.Errorf("YOUTUBE_API_KEY environment variable is required")
	}

	region := getEnvOrDefault("GCP_REGION", "us-central1")
	taskQueue := getEnvOrDefault("TASK_QUEUE", "video-snapshots")
	eventTopic := getEnvOrDefault("EVENT_TOPIC", "ingestion-events")

	// Initialize database
	db, err := datastore.OpenPostgres("")
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	return BootstrapHTTP(addr, db, projectID, youtubeAPIKey, region, taskQueue, eventTopic)
}

// getEnvOrDefault returns environment variable value or default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
