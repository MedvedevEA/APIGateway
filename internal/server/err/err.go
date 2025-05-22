package err

import "errors"

var (
	ErrRouteNotFound      = errors.New("route not found")
	ErrInvalidTokenFormat = errors.New("invalid token format")
	ErrInvalidTokenType   = errors.New("invalid token type")
	ErrInvalidToken       = errors.New("invalid token")
)
