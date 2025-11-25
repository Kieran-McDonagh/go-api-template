package types

import "errors"

type APIError error

var (
	ErrNotFound       = errors.New("not found")
	ErrNotUniqueEmail = errors.New("not unique email")
	ErrIncorrectLogin = errors.New("incorrect email or password")
)
