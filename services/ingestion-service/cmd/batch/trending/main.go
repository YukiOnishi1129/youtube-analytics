package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/adapter/gateway/postgres"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/adapter/gateway/pubsub"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/adapter/gateway/youtube"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/driver/config"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/driver/datastore"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/usecase"
)

func main() {
	// Parse command line arguments
	var (
		genreID = flag.String("genre", "", "Genre ID to collect trending videos for (optional, all genres if not specified)")
		dryRun  = flag.Bool("dry-run", false, "Dry run mode - only log what would be done")
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
	channelRepo := postgres.NewChannelRepository(pgRepo)
	videoRepo := postgres.NewVideoRepository(pgRepo)
	genreRepo := postgres.NewGenreRepository(pgRepo)
	// keywordGroupRepo := postgres.NewKeywordGroupRepository(pgRepo) // TODO: Implement

	// Initialize gateways
	youtubeClient, err := youtube.NewClient(cfg.YouTubeAPIKey)
	if err != nil {
		log.Fatalf("Failed to create YouTube client: %v", err)
	}
	eventPublisher, err := pubsub.NewEventPublisher(cfg.PubSubProjectID)
	if err != nil {
		log.Fatalf("Failed to create event publisher: %v", err)
	}
	// idGen := uuid.NewGenerator() // Not needed for this batch

	// Initialize use cases
	videoUseCase := usecase.NewVideoUseCase(
		videoRepo,
		channelRepo,
		youtubeClient,
		eventPublisher,
	)

	genreUseCase := usecase.NewGenreUseCase(genreRepo)
	// keywordUseCase := usecase.NewKeywordUseCase(keywordGroupRepo) // TODO: Implement

	// For now, use video use case directly
	// TODO: Implement proper trending collection use case

	// Log start
	log.Printf("Starting trending video collection batch (dry-run=%v)", *dryRun)
	start := time.Now()

	// Execute collection
	if *genreID != "" {
		// Collect for specific genre
		log.Printf("Collecting trending videos for genre: %s", *genreID)
		result, err := videoUseCase.CollectTrending(ctx, genreID)
		if err != nil {
			log.Fatalf("Failed to collect trending videos: %v", err)
		}
		log.Printf("Completed: collected=%d, created=%d, updated=%d, duration=%s",
			result.VideosCollected, result.VideosCreated, result.VideosUpdated, result.Duration)
	} else {
		// Collect for all enabled genres
		log.Println("Collecting trending videos for all enabled genres")

		// Get all enabled genres
		genres, err := genreUseCase.ListGenres(ctx)
		if err != nil {
			log.Fatalf("Failed to get genres: %v", err)
		}

		totalCollected := 0
		totalCreated := 0
		totalUpdated := 0

		for _, genre := range genres {
			// Skip disabled genres
			if !genre.Enabled {
				continue
			}

			log.Printf("Processing genre: %s", genre.Code)
			genreIDStr := string(genre.ID)
			result, err := videoUseCase.CollectTrending(ctx, &genreIDStr)
			if err != nil {
				log.Printf("Error collecting for genre %s: %v", genre.Code, err)
				continue
			}
			totalCollected += result.VideosCollected
			totalCreated += result.VideosCreated
			totalUpdated += result.VideosUpdated
			log.Printf("  Genre %s: collected=%d, created=%d, updated=%d",
				genre.Code, result.VideosCollected, result.VideosCreated, result.VideosUpdated)
		}

		log.Printf("Completed: genres=%d, total_collected=%d, total_created=%d, total_updated=%d, duration=%s",
			len(genres), totalCollected, totalCreated, totalUpdated, time.Since(start))
	}

	log.Printf("Total execution time: %s", time.Since(start))
}
