package security

import (
    "context"
    "strings"
    "google.golang.org/grpc"
    "google.golang.org/grpc/metadata"
    outgateway "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/port/output/gateway"
    adaptersec "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/adapter/security"
)

// UnaryAuthInterceptor verifies Authorization: Bearer <token> with the provided TokenVerifier,
// and injects TokenClaims into context for use cases via ClaimsProvider.
func UnaryAuthInterceptor(verifier outgateway.TokenVerifier) grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        md, ok := metadata.FromIncomingContext(ctx)
        if !ok {
            return handler(ctx, req)
        }
        authz := ""
        if vals := md.Get("authorization"); len(vals) > 0 {
            authz = vals[0]
        }
        if strings.HasPrefix(strings.ToLower(authz), "bearer ") {
            token := strings.TrimSpace(authz[len("bearer "):])
            if token != "" {
                if claims, err := verifier.Verify(ctx, token); err == nil {
                    ctx = adaptersec.WithClaims(ctx, outgateway.TokenClaims{
                        Subject:       claims.Subject,
                        Email:         claims.Email,
                        EmailVerified: claims.EmailVerified,
                        DisplayName:   claims.DisplayName,
                        PhotoURL:      claims.PhotoURL,
                        Provider:      claims.Provider,
                        ProviderUID:   claims.ProviderUID,
                    })
                }
            }
        }
        return handler(ctx, req)
    }
}

