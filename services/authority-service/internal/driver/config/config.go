package config

import (
    "os"
)

// Config holds runtime configuration loaded from environment variables.
type Config struct {
    GRPCAddr      string // e.g. ":8080"
    DatabaseURL   string // Postgres DSN
    DBHost        string // Postgres host
    DBPort        string // Postgres port
    DBUser        string // Postgres user
    DBPassword    string // Postgres password
    DBName        string // Postgres database name
    DBSSLMode     string // Postgres sslmode
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
        DBHost:        os.Getenv("DB_HOST"),
        DBPort:        getEnv("DB_PORT", ""),
        DBUser:        os.Getenv("DB_USER"),
        DBPassword:    os.Getenv("DB_PASSWORD"),
        DBName:        os.Getenv("DB_NAME"),
        DBSSLMode:     os.Getenv("DB_SSLMODE"),
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
