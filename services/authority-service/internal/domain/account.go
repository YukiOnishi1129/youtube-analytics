package domain

import (
    "errors"
    "strings"
    "time"
)

// Account represents a user account within the Authority context aggregate root.
// Aggregates: Account (root) -> Identities (children)
// Invariants:
//  - email is unique among active accounts (enforced by repository/storage)
//  - (account_id, provider) is unique (enforced by repository/storage)
//  - at least one identity must remain
//  - multiple providers can be linked
type Account struct {
    ID            string
    Email         string
    EmailVerified bool
    DisplayName   string
    PhotoURL      string
    IsActive      bool
    LastLoginAt   *time.Time

    // Child entities
    Identities []Identity
    Roles      []Role
}

// Domain errors
var (
    ErrInvalidEmail          = errors.New("invalid email")
    ErrIdentityAlreadyLinked = errors.New("identity already linked")
    ErrIdentityNotFound      = errors.New("identity not found")
    ErrCannotUnlinkLastID    = errors.New("cannot unlink the last identity")
    ErrAccountInactive       = errors.New("account is inactive")
    ErrRoleAlreadyAssigned   = errors.New("role already assigned")
    ErrRoleNotAssigned       = errors.New("role not assigned")
    ErrNotFound              = errors.New("account not found")
    ErrEmailAlreadyUsed      = errors.New("email already in use")
)

// NewAccount creates a new Account aggregate with initial values.
// ID generation policy is handled by the caller (infra), expected to be UUID v7.
func NewAccount(id string, email string) (*Account, error) {
    if !isValidEmail(email) {
        return nil, ErrInvalidEmail
    }
    return &Account{
        ID:            id,
        Email:         strings.ToLower(strings.TrimSpace(email)),
        EmailVerified: false,
        DisplayName:   "",
        PhotoURL:      "",
        IsActive:      true,
        LastLoginAt:   nil,
        Identities:    []Identity{},
        Roles:         []Role{},
    }, nil
}

// UpdateProfile updates profile fields.
func (a *Account) UpdateProfile(displayName, photoURL string) {
    if displayName != "" {
        a.DisplayName = displayName
    }
    if photoURL != "" {
        a.PhotoURL = photoURL
    }
}

// VerifyEmail marks the account's email as verified.
func (a *Account) VerifyEmail() { a.EmailVerified = true }

// TouchLogin updates last login timestamp and activation state.
func (a *Account) TouchLogin(now time.Time) {
    a.LastLoginAt = &now
    if !a.IsActive {
        a.IsActive = true
    }
}

// LinkIdentity links an authentication provider to this account if not already linked.
func (a *Account) LinkIdentity(p Provider, providerUID string) error {
    for _, id := range a.Identities {
        if id.Provider == p {
            return ErrIdentityAlreadyLinked
        }
    }
    a.Identities = append(a.Identities, Identity{Provider: p, ProviderUID: providerUID})
    return nil
}

// HasProvider returns true if the account has a link to the specified provider.
func (a *Account) HasProvider(p Provider) bool {
    for _, id := range a.Identities {
        if id.Provider == p {
            return true
        }
    }
    return false
}

// ReplaceIdentity replaces or inserts an identity for the given provider.
// Useful when syncing provider UID changes from upstream identity platform.
func (a *Account) ReplaceIdentity(p Provider, providerUID string) {
    for i, id := range a.Identities {
        if id.Provider == p {
            a.Identities[i].ProviderUID = providerUID
            return
        }
    }
    a.Identities = append(a.Identities, Identity{Provider: p, ProviderUID: providerUID})
}

// UnlinkIdentity removes an identity; enforces at least one identity remains.
func (a *Account) UnlinkIdentity(p Provider) error {
    idx := -1
    for i, id := range a.Identities {
        if id.Provider == p {
            idx = i
            break
        }
    }
    if idx == -1 {
        return ErrIdentityNotFound
    }
    if len(a.Identities) <= 1 {
        return ErrCannotUnlinkLastID
    }
    a.Identities = append(a.Identities[:idx], a.Identities[idx+1:]...)
    return nil
}

// Deactivate marks the account as inactive (cannot sign in).
func (a *Account) Deactivate() { a.IsActive = false }

// Reactivate marks the account as active again.
func (a *Account) Reactivate() { a.IsActive = true }

// AssignRole assigns a role if not already present.
func (a *Account) AssignRole(r Role) error {
    if a.hasRole(r.Name) {
        return ErrRoleAlreadyAssigned
    }
    a.Roles = append(a.Roles, r)
    return nil
}

// RevokeRole removes a role if present.
func (a *Account) RevokeRole(r Role) error {
    for i, rr := range a.Roles {
        if equalRole(rr, r) {
            a.Roles = append(a.Roles[:i], a.Roles[i+1:]...)
            return nil
        }
    }
    return ErrRoleNotAssigned
}

func (a *Account) hasRole(name string) bool {
    for _, r := range a.Roles {
        if strings.EqualFold(r.Name, name) {
            return true
        }
    }
    return false
}

func equalRole(aRole, bRole Role) bool {
    return strings.EqualFold(aRole.Name, bRole.Name)
}

// isValidEmail performs a minimal sanity check. Full validation belongs to application/infrastructure.
func isValidEmail(email string) bool {
    e := strings.TrimSpace(strings.ToLower(email))
    return e != "" && strings.Contains(e, "@")
}
