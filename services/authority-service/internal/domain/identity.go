package domain

// Identity stores authentication provider linkage for an account.
type Identity struct {
    Provider    Provider
    ProviderUID string
}

// Provider represents an authentication provider kind.
type Provider string

const (
    ProviderGoogle   Provider = "google"
    ProviderPassword Provider = "password"
    ProviderGithub   Provider = "github"
)

