package presenter

import (
    account "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/domain/account"
)

// AuthorityPresenter is the output boundary that shapes responses
// for authority use cases.
type AuthorityPresenter interface {
    // PresentGetMe outputs the current user's account profile.
    PresentGetMe(a *account.Account) error
}

