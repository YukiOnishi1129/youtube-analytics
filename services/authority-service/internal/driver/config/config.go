package config

import (
    "os"
)

// Config holds runtime configuration loaded from environment variables.
type Config struct {
    GRPCAddr      string // e.g. ":8080"
    DatabaseURL   string // Postgres DSN
    MigrationsDir string // optional, path to migrations

    FirebaseAPIKey string // Identity Platform API key
    OIDCIssuer     string // OIDC issuer (e.g., https://securetoken.google.com/<PROJECT>)
    OIDCAudience   string // OIDC audience/clientID (usually PROJECT_ID)
}

// Load loads configuration from environment variables with sensible defaults.
func Load() Config {
    cfg := Config{
        GRPCAddr:      getEnv("GRPC_ADDR", ":8080"),
        DatabaseURL:   os.Getenv("DATABASE_URL"),
        MigrationsDir: getEnv("MIGRATIONS_DIR", "services/authority-service/internal/driver/datastore/migrations"),
        FirebaseAPIKey: os.Getenv("FIREBASE_API_KEY"),
        OIDCIssuer:     os.Getenv("OIDC_ISSUER"),
        OIDCAudience:   os.Getenv("OIDC_AUDIENCE"),
    }
    return cfg
}

func getEnv(k, def string) string {
    if v := os.Getenv(k); v != "" { return v }
    return def
}

