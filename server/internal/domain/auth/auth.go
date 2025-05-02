package auth

import (
	"skyvault/pkg/apperror"
	"skyvault/pkg/utils"
	"skyvault/pkg/validate"
	"time"
)

type Auth struct {
	ID             int64
	ProfileID      int64
	Provider       Provider // E.g., "email", "oidc", "ldap"
	ProviderUserID string   // userID provided by the provider
	PasswordHash   *string  // Nil if not using email provider
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Provider string

const (
	ProviderEmail Provider = "email"
	ProviderOIDC  Provider = "oidc"
	ProviderLDAP  Provider = "ldap"
)

// App Errors:
// - ErrCommonInvalidValue
func NewAuth(profileID int64, provider Provider, providerUserID string, password *string) (*Auth, error) {
	var passwordHash *string
	if provider == ProviderEmail {
		if password == nil {
			return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "auth.NewAuth:Password")
		}

		if email, err := validate.Email(providerUserID); err != nil {
			return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "auth.NewAuth:ValidateEmail")
		} else {
			providerUserID = email
		}

		hash, err := utils.HashPassword(*password)
		if err != nil {
			return nil, apperror.NewAppError(err, "auth.NewAuth:HashPassword")
		}

		passwordHash = &hash
	}

	now := time.Now().UTC()
	return &Auth{
		ProfileID:      profileID,
		Provider:       provider,
		ProviderUserID: providerUserID,
		PasswordHash:   passwordHash,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// App Errors:
// - ErrCommonNoAccess
func (a *Auth) ValidateAccess(accessedByID int64) error {
	if a.ProfileID != accessedByID {
		return apperror.ErrCommonNoAccess
	}
	return nil
}
