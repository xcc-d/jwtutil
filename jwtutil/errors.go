package jwtutil

import "errors"

// 1
var (
	ErrInvalidToken         = errors.New("invalid token")
	ErrTokenExpired         = errors.New("token expired")
	ErrInvalidSigningMethod = errors.New("invalid signing method")
	ErrMissingKey           = errors.New("signing key is missing")
	ErrInvalidClaims        = errors.New("invalid claims")
	ErrTokenMalformed       = errors.New("token is malformed")
)
