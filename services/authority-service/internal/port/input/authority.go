package input

import "context"

// AuthorityInputPort defines use case entrypoints for the authority-service.
type AuthorityInputPort interface {
    // GetAccount presents the current user's account using auth context.
    GetAccount(ctx context.Context) error

    // SignUp registers a new account via identity provider and persists account profile.
    SignUp(ctx context.Context, email, password string) error

    // SignIn authenticates via identity provider and returns tokens.
    SignIn(ctx context.Context, email, password string) error

    // SignOut revokes refresh token via identity provider.
    SignOut(ctx context.Context, refreshToken string) error

    // ResetPassword triggers a password reset flow via identity provider.
    ResetPassword(ctx context.Context, email string) error
}
