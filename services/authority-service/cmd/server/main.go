package main

import (
    "context"
    "log"
    "time"

    fb "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/adapter/gateway/firebase"
    cfgpkg "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/driver/config"
    "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/driver/datastore"
    "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/driver/transport"
    outgateway "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/port/output/gateway"
)

type systemClock struct{}

func (systemClock) Now() time.Time { return time.Now() }

// TokenVerifier stub: trusts any token and builds minimal claims (DO NOT USE IN PROD)
type noopVerifier struct{}

func (noopVerifier) Verify(_ context.Context, _ string) (outgateway.TokenClaims, error) {
	return outgateway.TokenClaims{Email: "stub@example.com", EmailVerified: true}, nil
}

func main() {
	// Load configuration
	cfg := cfgpkg.Load()

	// Repositories (wired to Postgres)
	var (
		accountRepo outgateway.AccountRepository
		idRepo      outgateway.IdentityRepository
		roleRepo    outgateway.RoleRepository
	)
	var idp outgateway.IdentityProvider = noopIDP{}
	var verifier outgateway.TokenVerifier = noopVerifier{}
	var clock outgateway.Clock = systemClock{}

	// Datastore: open Postgres and run migrations (no-op unless built with tags)
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	db, err := datastore.OpenPostgres(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Postgres connect failed: %v", err)
	}
    // Wire repositories via driver (requires -tags sqlc)
    ar, ir, rr, err := datastore.NewRepositories(db)
    if err != nil {
        log.Fatalf("repository wiring failed (build with -tags sqlc): %v", err)
    }
    accountRepo, idRepo, roleRepo = ar, ir, rr

	// Identity Platform & OIDC
	if cfg.FirebaseAPIKey != "" {
		idp = fb.New(cfg.FirebaseAPIKey)
	}
	if cfg.OIDCIssuer != "" && cfg.OIDCAudience != "" {
		if v, err := fb.NewOIDCVerifier(context.Background(), cfg.OIDCIssuer, cfg.OIDCAudience); err == nil {
			verifier = v
		} else {
			log.Printf("OIDC verifier init failed, using noop: %v", err)
		}
	}

	if err := transport.Bootstrap(cfg.GRPCAddr, accountRepo, idRepo, roleRepo, verifier, idp, clock); err != nil {
		log.Fatal(err)
	}
}

// noopIDP implements IdentityProvider for local runs without Firebase.
type noopIDP struct{}

func (noopIDP) SignUp(ctx context.Context, email, password string) (outgateway.AuthTokens, error) {
	return outgateway.AuthTokens{IDToken: "stub", RefreshToken: "stub", ExpiresIn: 3600}, nil
}
func (noopIDP) SignIn(ctx context.Context, email, password string) (outgateway.AuthTokens, error) {
	return outgateway.AuthTokens{IDToken: "stub", RefreshToken: "stub", ExpiresIn: 3600}, nil
}
func (noopIDP) SignOut(ctx context.Context, refreshToken string) error { return nil }
func (noopIDP) ResetPassword(ctx context.Context, email string) error  { return nil }
