package presenter

import "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/domain"

// AuthorityPresenter is the output boundary that shapes responses for authority use cases.
type AuthorityPresenter interface {
    // PresentGetAccount outputs the current user's account profile.
    PresentGetAccount(a *domain.Account) error

    // PresentSignUp outputs result for SignUp (tokens + account if needed).
    PresentSignUp(account *domain.Account, idToken, refreshToken string) error

    // PresentSignIn outputs tokens after successful authentication.
    PresentSignIn(idToken, refreshToken string, expiresIn int64) error

    // PresentSignOut confirms successful sign-out.
    PresentSignOut(success bool) error

    // PresentResetPassword confirms password reset initiation.
    PresentResetPassword(emailSent bool) error
}
