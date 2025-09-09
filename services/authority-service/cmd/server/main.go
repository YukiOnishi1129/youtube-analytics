package main

import (
    "context"
    "log"
    "time"

    fb "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/adapter/gateway/firebase"
    insecure "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/adapter/gateway/insecure"
    cfgpkg "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/driver/config"
    "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/driver/datastore"
    "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/driver/transport"
    outgateway "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/port/output/gateway"
)

type systemClock struct{}

func (systemClock) Now() time.Time { return time.Now() }

func main() {
	// Load configuration
	cfg := cfgpkg.Load()

	// Repositories (wired to Postgres)
	var (
		accountRepo outgateway.AccountRepository
		idRepo      outgateway.IdentityRepository
		roleRepo    outgateway.RoleRepository
	)
	var idp outgateway.IdentityProvider
	var verifier outgateway.TokenVerifier
	var clock outgateway.Clock = systemClock{}

	// Datastore: open Postgres (reads DB_* envs if DATABASE_URL is empty)
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

    // Identity Platform (fallback to dummy in local when unset)
    if cfg.FirebaseAPIKey == "" {
        log.Printf("[WARN] FIREBASE_API_KEY not set; using dummy IdentityProvider for local dev")
        idp = insecure.DummyIDP{}
    } else {
        idp = fb.New(cfg.FirebaseAPIKey)
    }

    // OIDC verifier (fallback to no-op in local when unset)
    if cfg.OIDCIssuer == "" || cfg.OIDCAudience == "" {
        log.Printf("[WARN] OIDC_ISSUER/AUDIENCE not set; using NoopVerifier for local dev")
        verifier = insecure.NoopVerifier{}
    } else {
        v, err := fb.NewOIDCVerifier(context.Background(), cfg.OIDCIssuer, cfg.OIDCAudience)
        if err != nil {
            log.Fatalf("OIDC verifier init failed: %v", err)
        }
        verifier = v
    }

	if err := transport.Bootstrap(cfg.GRPCAddr, accountRepo, idRepo, roleRepo, verifier, idp, clock); err != nil {
		log.Fatal(err)
	}
}
