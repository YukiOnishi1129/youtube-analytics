package usecase

import (
    "context"
    "strings"
    "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/domain"
    inport "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/port/input"
    outgateway "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/port/output/gateway"
    outpresenter "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/port/output/presenter"
)

// Ensure Interactor implements the input port.
var _ inport.AuthorityInputPort = (*AuthorityInteractor)(nil)

type AuthorityInteractor struct {
    accounts  outgateway.AccountRepository
    ids       outgateway.IdentityRepository
    roles     outgateway.RoleRepository
    verifier  outgateway.TokenVerifier
    claims    outgateway.ClaimsProvider
    idp       outgateway.IdentityProvider
    clock     outgateway.Clock
    presenter outpresenter.AuthorityPresenter
}

func NewAuthorityInteractor(
    accounts outgateway.AccountRepository,
    ids outgateway.IdentityRepository,
    roles outgateway.RoleRepository,
    verifier outgateway.TokenVerifier,
    claims outgateway.ClaimsProvider,
    idp outgateway.IdentityProvider,
    clock outgateway.Clock,
    presenter outpresenter.AuthorityPresenter,
) *AuthorityInteractor {
    return &AuthorityInteractor{
        accounts:  accounts,
        ids:       ids,
        roles:     roles,
        verifier:  verifier,
        claims:    claims,
        idp:       idp,
        clock:     clock,
        presenter: presenter,
    }
}

// GetAccount verifies the idToken and presents the account profile.
// Minimal happy-path: find by email first, then by provider mapping.
func (uc *AuthorityInteractor) GetAccount(ctx context.Context) error {
    claims, err := uc.claims.Current(ctx)
    if err != nil {
        return err
    }

    var acc *domain.Account

    if claims.Email != "" {
        if a, err := uc.accounts.FindByEmail(ctx, claims.Email); err == nil {
            acc = a
        }
    }

    if acc == nil && claims.Provider != "" && claims.ProviderUID != "" {
        if a, err := uc.ids.FindByProvider(ctx, claims.Provider, claims.ProviderUID); err == nil {
            acc = a
        }
    }

    if acc == nil {
        return domain.ErrNotFound
    }

    if !acc.IsActive {
        return domain.ErrAccountInactive
    }

    return uc.presenter.PresentGetAccount(acc)
}

// SignUp creates an account via identity provider and persists the local profile.
func (uc *AuthorityInteractor) SignUp(ctx context.Context, email, password string) error {
    email = strings.TrimSpace(strings.ToLower(email))
    if email == "" || password == "" {
        return domain.ErrInvalidEmail
    }

    // 1) Create user in IDP and get tokens
    tokens, err := uc.idp.SignUp(ctx, email, password)
    if err != nil {
        return err
    }

    // 2) Create local account if not exists
    a, err := uc.accounts.FindByEmail(ctx, email)
    if err != nil {
        // assume not found → create
        a, err = domain.NewAccount("", email)
        if err != nil {
            return err
        }
        // default role user
        _ = a.AssignRole(domain.Role{Name: "user"})
        // link password identity
        _ = a.LinkIdentity(domain.ProviderPassword, email)
        if err := uc.accounts.Save(ctx, a); err != nil {
            return err
        }
    }

    return uc.presenter.PresentSignUp(a, tokens.IDToken, tokens.RefreshToken)
}

// SignIn authenticates via IDP and returns tokens; touches last login if account exists.
func (uc *AuthorityInteractor) SignIn(ctx context.Context, email, password string) error {
    email = strings.TrimSpace(strings.ToLower(email))
    if email == "" || password == "" {
        return domain.ErrInvalidEmail
    }

    tokens, err := uc.idp.SignIn(ctx, email, password)
    if err != nil {
        return err
    }

    if a, err := uc.accounts.FindByEmail(ctx, email); err == nil {
        if !a.IsActive {
            return domain.ErrAccountInactive
        }
        a.TouchLogin(uc.clock.Now())
        _ = uc.accounts.Save(ctx, a)
    }

    return uc.presenter.PresentSignIn(tokens.IDToken, tokens.RefreshToken, tokens.ExpiresIn)
}

// SignOut revokes refresh token via IDP.
func (uc *AuthorityInteractor) SignOut(ctx context.Context, refreshToken string) error {
    if err := uc.idp.SignOut(ctx, refreshToken); err != nil {
        return err
    }
    return uc.presenter.PresentSignOut(true)
}

// ResetPassword triggers password reset email via IDP.
func (uc *AuthorityInteractor) ResetPassword(ctx context.Context, email string) error {
    email = strings.TrimSpace(strings.ToLower(email))
    if email == "" {
        return domain.ErrInvalidEmail
    }
    if err := uc.idp.ResetPassword(ctx, email); err != nil {
        return err
    }
    return uc.presenter.PresentResetPassword(true)
}
