package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	_ "embed"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/driver/datastore"
	_ "github.com/lib/pq"
)

// Master data (production-ready)
//go:embed seeds/master/youtube_categories.sql
var youtubeCategoriesSeedSQL string

//go:embed seeds/master/genres.sql
var genresSeedSQL string

//go:embed seeds/master/keyword_groups.sql
var keywordGroupsSeedSQL string

//go:embed seeds/master/keyword_items.sql
var keywordItemsSeedSQL string

// All seeds are master data
var seedFiles = map[string]string{
	"youtube_categories": youtubeCategoriesSeedSQL,
	"genres":             genresSeedSQL,
	"keyword_groups":     keywordGroupsSeedSQL,
	"keyword_items":      keywordItemsSeedSQL,
}

func main() {
	var (
		target = flag.String("target", "all", "Seed target: all, youtube_categories, genres, keyword_groups, keyword_items")
		dryRun = flag.Bool("dry-run", false, "Show SQL without executing")
	)
	flag.Parse()

	if *dryRun {
		showSQL(*target)
		return
	}

	// Connect to database
	db, err := datastore.OpenPostgres("")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Execute seeds
	if err := executeSeed(db, *target); err != nil {
		log.Fatalf("Failed to execute seed: %v", err)
	}

	log.Printf("Seeding completed successfully for target: %s", *target)
}

func executeSeed(db *sql.DB, target string) error {
	switch target {
	case "all":
		// Execute in dependency order
		seedOrder := []string{"youtube_categories", "genres", "keyword_groups", "keyword_items"}
		for _, name := range seedOrder {
			sql, ok := seedFiles[name]
			if !ok {
				continue
			}
			log.Printf("Executing seed: %s", name)
			if _, err := db.Exec(sql); err != nil {
				return fmt.Errorf("failed to execute %s: %w", name, err)
			}
		}
	case "youtube_categories", "genres", "keyword_groups", "keyword_items":
		sql, ok := seedFiles[target]
		if !ok {
			return fmt.Errorf("unknown target: %s", target)
		}
		log.Printf("Executing seed: %s", target)
		if _, err := db.Exec(sql); err != nil {
			return fmt.Errorf("failed to execute %s: %w", target, err)
		}
	default:
		return fmt.Errorf("unknown target: %s (available: all, youtube_categories, genres, keyword_groups, keyword_items)", target)
	}
	return nil
}

func showSQL(target string) {
	switch target {
	case "all":
		seedOrder := []string{"youtube_categories", "genres", "keyword_groups", "keyword_items"}
		for _, name := range seedOrder {
			sql, ok := seedFiles[name]
			if !ok {
				continue
			}
			fmt.Printf("-- Seed: %s\n%s\n\n", name, sql)
		}
	case "youtube_categories", "genres", "keyword_groups", "keyword_items":
		sql, ok := seedFiles[target]
		if !ok {
			log.Fatalf("Unknown target: %s", target)
		}
		fmt.Printf("-- Seed: %s\n%s\n", target, sql)
	default:
		log.Fatalf("Unknown target: %s (available: all, youtube_categories, genres, keyword_groups, keyword_items)", target)
	}
}