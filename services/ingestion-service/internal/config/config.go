package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for the ingestion service
type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	YouTube    YouTubeConfig
	GCP        GCPConfig
	CloudTasks CloudTasksConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port     int
	GRPCPort int
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	URL string
}

// YouTubeConfig holds YouTube API configuration
type YouTubeConfig struct {
	APIKey string
}

// GCPConfig holds Google Cloud Platform configuration
type GCPConfig struct {
	ProjectID string
}

// CloudTasksConfig holds Cloud Tasks configuration
type CloudTasksConfig struct {
	Location  string
	QueueName string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{}

	// Server configuration
	if port := os.Getenv("PORT"); port != "" {
		p, err := strconv.Atoi(port)
		if err != nil {
			return nil, fmt.Errorf("invalid PORT: %w", err)
		}
		cfg.Server.Port = p
	} else {
		cfg.Server.Port = 8080 // Default HTTP port
	}

	if grpcPort := os.Getenv("GRPC_PORT"); grpcPort != "" {
		p, err := strconv.Atoi(grpcPort)
		if err != nil {
			return nil, fmt.Errorf("invalid GRPC_PORT: %w", err)
		}
		cfg.Server.GRPCPort = p
	} else {
		cfg.Server.GRPCPort = 50051 // Default gRPC port
	}

	// Database configuration
	cfg.Database.URL = os.Getenv("DATABASE_URL")
	if cfg.Database.URL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	// YouTube configuration
	cfg.YouTube.APIKey = os.Getenv("YOUTUBE_API_KEY")
	if cfg.YouTube.APIKey == "" {
		return nil, fmt.Errorf("YOUTUBE_API_KEY is required")
	}

	// GCP configuration
	cfg.GCP.ProjectID = os.Getenv("GCP_PROJECT_ID")
	if cfg.GCP.ProjectID == "" {
		return nil, fmt.Errorf("GCP_PROJECT_ID is required")
	}

	// Cloud Tasks configuration
	cfg.CloudTasks.Location = os.Getenv("CLOUD_TASKS_LOCATION")
	if cfg.CloudTasks.Location == "" {
		cfg.CloudTasks.Location = "us-central1" // Default location
	}

	cfg.CloudTasks.QueueName = os.Getenv("CLOUD_TASKS_QUEUE_NAME")
	if cfg.CloudTasks.QueueName == "" {
		cfg.CloudTasks.QueueName = "ingestion-tasks" // Default queue name
	}

	return cfg, nil
}