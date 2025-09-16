package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/adapter/gateway/postgres"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/driver/config"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/driver/datastore"
)

func main() {
	// Parse command line arguments
	var (
		days   = flag.Int("days", 7, "Renew subscriptions expiring within N days")
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
	channelRepo := postgres.NewChannelRepository(pgRepo)

	// Log start
	log.Printf("Starting WebSub renewal batch (days=%d, dry-run=%v)", *days, *dryRun)
	
	// TODO: Implement WebSub client and subscription renewal use case
	log.Printf("WebSub renewal batch not yet implemented")
	log.Printf("Would check channels from repository and renew subscriptions expiring in %d days", *days)
	
	// For now, just show what channels we would process
	channels, err := channelRepo.ListSubscribed(ctx)
	if err != nil {
		log.Printf("Error listing channels: %v", err)
	} else {
		log.Printf("Found %d subscribed channels that might need renewal", len(channels))
	}
}