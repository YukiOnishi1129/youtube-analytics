package identityauth

import (
    "context"
)

type Claims struct {
    Subject       string
    Email         string
    EmailVerified bool
    DisplayName   string
    PhotoURL      string
}

type ctxKey string

const claimsKey ctxKey = "identityauth.claims"

func WithClaims(ctx context.Context, c Claims) context.Context {
    return context.WithValue(ctx, claimsKey, c)
}

func FromContext(ctx context.Context) (Claims, bool) {
    v := ctx.Value(claimsKey)
    if v == nil {
        return Claims{}, false
    }
    c, ok := v.(Claims)
    return c, ok
}

