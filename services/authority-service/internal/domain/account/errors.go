package account

import "errors"

var (
    ErrNotFound        = errors.New("account not found")
    ErrEmailAlreadyUsed = errors.New("email already in use")
)
