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

//go:embed seeds/keywords.sql
var keywordsSeedSQL string

var seedFiles = map[string]string{
	"keywords": keywordsSeedSQL,
}

func main() {
	var (
		target = flag.String("target", "all", "Seed target: all, keywords")
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
		for name, sql := range seedFiles {
			log.Printf("Executing seed: %s", name)
			if _, err := db.Exec(sql); err != nil {
				return fmt.Errorf("failed to execute %s: %w", name, err)
			}
		}
	case "keywords":
		sql, ok := seedFiles[target]
		if !ok {
			return fmt.Errorf("unknown target: %s", target)
		}
		log.Printf("Executing seed: %s", target)
		if _, err := db.Exec(sql); err != nil {
			return fmt.Errorf("failed to execute %s: %w", target, err)
		}
	default:
		return fmt.Errorf("unknown target: %s (available: all, keywords)", target)
	}
	return nil
}

func showSQL(target string) {
	switch target {
	case "all":
		for name, sql := range seedFiles {
			fmt.Printf("-- Seed: %s\n%s\n\n", name, sql)
		}
	case "keywords":
		sql, ok := seedFiles[target]
		if !ok {
			log.Fatalf("Unknown target: %s", target)
		}
		fmt.Printf("-- Seed: %s\n%s\n", target, sql)
	default:
		log.Fatalf("Unknown target: %s (available: all, keywords)", target)
	}
}