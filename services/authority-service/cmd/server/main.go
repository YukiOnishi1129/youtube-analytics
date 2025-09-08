package main

import (
    "context"
    "log"
    "time"
    "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/adapter/gateway"
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
    // In-memory repos and stubs for local bootstrap
    accountRepo := gateway.NewInMemoryAccountRepo()
    idRepo := &gateway.InMemoryIdentityRepo{}
    roleRepo := &gateway.InMemoryRoleRepo{}
    idp := &gateway.IDPClient{}
    var verifier outgateway.TokenVerifier = noopVerifier{}
    var clock outgateway.Clock = systemClock{}

    if err := transport.Bootstrap(":8080", accountRepo, idRepo, roleRepo, verifier, idp, clock); err != nil {
        log.Fatal(err)
    }
}
