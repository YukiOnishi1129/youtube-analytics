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
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/driver/config"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/driver/datastore"
)

func main() {
	// Parse command line arguments
	var (
		genreID    = flag.String("genre", "", "Genre ID to generate rankings for (optional, all genres if not specified)")
		checkpoint = flag.Int("checkpoint", 24, "Checkpoint hour for ranking (default: 24h)")
		topN       = flag.Int("top", 10, "Number of top videos to include in ranking")
		dryRun     = flag.Bool("dry-run", false, "Dry run mode - only log what would be done")
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
	genreRepo := postgres.NewGenreRepository(pgRepo)
	// metricsRepo := postgres.NewVideoMetricsRepository(pgRepo) // TODO: Implement
	// rankingRepo := postgres.NewRankingHistoryRepository(pgRepo) // TODO: Implement

	// Log start
	log.Printf("Starting ranking generation batch (genre=%s, checkpoint=%dh, top=%d, dry-run=%v)",
		*genreID, *checkpoint, *topN, *dryRun)
	start := time.Now()

	// TODO: Implement ranking generation
	log.Printf("Ranking generation batch not yet implemented")
	log.Printf("Would generate top %d ranking for checkpoint %dh", *topN, *checkpoint)
	
	if *genreID != "" {
		// Generate for specific genre
		id := valueobject.UUID(*genreID)
		genre, err := genreRepo.FindByID(ctx, id)
		if err != nil {
			log.Printf("Error getting genre: %v", err)
		} else {
			log.Printf("Would generate ranking for genre: %s", genre.Code)
		}
		
		// Show sample of what videos would be ranked
		videos, err := videoRepo.ListActive(ctx, time.Now().Add(-7*24*time.Hour))
		if err != nil {
			log.Printf("Error getting videos: %v", err)
		} else {
			log.Printf("Sample active videos: %d", len(videos))
		}
	} else {
		// Generate for all genres
		genres, err := genreRepo.FindAll(ctx)
		if err != nil {
			log.Printf("Error getting genres: %v", err)
		} else {
			log.Printf("Would generate rankings for %d genres", len(genres))
			for _, genre := range genres {
				if genre.Enabled {
					log.Printf("  - %s", genre.Code)
				}
			}
		}
	}

	// Log results
	log.Printf("Completed in %s", time.Since(start))
}