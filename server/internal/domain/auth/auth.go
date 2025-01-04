package auth

import "errors"

var ErrAuthAlreadyExists = errors.New("auth method already exists")

type Auth struct {
	ID             int64
	ProfileID         int64
	Provider       Provider // E.g., "email", "oidc", "ldap"
	ProviderUserID string   // External user ID for provider
	PasswordHash   *string  // Optional if using external providers
}

type Provider string

const (
	ProviderEmail Provider = "email"
	ProviderOIDC  Provider = "oidc"
	ProviderLDAP  Provider = "ldap"
)

func Providers() []Provider {
	return []Provider{
		ProviderEmail,
		ProviderOIDC,
		ProviderLDAP,
	}
}

func NewAuth(profileID int64) *Auth {
	return &Auth{
		ProfileID: profileID,
		// Email is the default provider
		Provider: ProviderEmail,
	}
}
