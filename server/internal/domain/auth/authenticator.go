package auth

import (
	"context"
)

type Claims interface {
	GetProfileID() string
}

// CredKey is a key for credentials map.
// These keys depend on the Provider.
type CredKey string

const (
	CredKeyPassword     CredKey = "password"
	CredKeyPasswordHash CredKey = "password_hash"
)

// Authenticator is implemented by each Provider.
type Authenticator interface {
	GenerateToken(ctx context.Context, profileID string) (string, error)

	// App Errors:
	// - ErrAuthInvalidToken
	// - ErrAuthTokenExpired
	ValidateToken(ctx context.Context, token string) (Claims, error)

	// App Errors:
	// - ErrCommonInvalidValue
	// - ErrAuthInvalidCredentials
	ValidateCredentials(ctx context.Context, credentials map[CredKey]any) error
}

// AuthenticatorFactory returns an Authenticator instance for a given Provider.
type AuthenticatorFactory interface {
	// App Errors:
	// - ErrAuthWrongProvider
	GetAuthenticator(provider Provider) (Authenticator, error)
}
