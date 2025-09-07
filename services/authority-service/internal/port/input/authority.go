package input

import "context"

// AuthorityInputPort defines use case entrypoints for the authority-service.
// Presenter and gateways are injected into the interactor in the use case layer.
type AuthorityInputPort interface {
    // GetMe verifies the given idToken and emits the result via presenter.
    GetMe(ctx context.Context, idToken string) error
}

