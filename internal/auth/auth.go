package auth

import "errors"

type Auth interface {
	Auth(token string) (err error)
	GetToken() string
}

var (
	ErrAuthTokenInternal = errors.New("auth: token internal")
	ErrAuthTokenInvalid  = errors.New("auth: token invalid")
	ErrAuthTokenNotFound = errors.New("auth: token not found")
)
