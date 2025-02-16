package serverErrors

import "errors"

var (
	ErrItemNotFound    = errors.New("item not found")
	ErrUserNotFound    = errors.New("user not found")
	ErrBadCreditonals  = errors.New("bad creditionals")
	ErrBadToken        = errors.New("bad token")
	ErrInternal        = errors.New("internal error")
	ErrNotEnoughtMoney = errors.New("not enough money")
	ErrNotAllowed      = errors.New("not allowed to update such user")
)
