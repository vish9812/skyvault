package auth

import (
	"fmt"
	"skyvault/pkg/apperror"
	"skyvault/pkg/utils"
	"strings"
	"time"
)

const (
	passwordMinLen = 4
	passwordMaxLen = 50
)

type Auth struct {
	ID             int64
	ProfileID      int64
	Provider       Provider // E.g., "email", "oidc", "ldap"
	ProviderUserID string   // External user ID for provider
	PasswordHash   *string  // Optional if using external providers
	CreatedAt      time.Time
	UpdatedAt      time.Time
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

// App Errors:
// - apperror.ErrInvalidValue
func NewAuth(profileID int64, provider Provider, providerUserID string, password *string) (*Auth, error) {
	if provider == "" {
		return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "auth.NewAuth:Provider")
	}

	if providerUserID == "" {
		return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "auth.NewAuth:ProviderUserID:Empty")
	}

	var passwordHash *string
	if provider == ProviderEmail {
		if err := utils.IsValidEmail(providerUserID); err != nil {
			return nil, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "auth.NewAuth:IsValidEmail")
		}

		if password == nil {
			return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "auth.NewAuth:Password:Empty")
		}

		pwd := strings.TrimSpace(*password)
		if len(pwd) < passwordMinLen {
			return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "auth.NewAuth:Password:MinLength")
		}

		if len(pwd) > passwordMaxLen {
			return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "auth.NewAuth:Password:MaxLength")
		}

		hash, err := utils.HashPassword(pwd)
		if err != nil {
			return nil, apperror.NewAppError(err, "auth.NewAuth:HashPassword")
		}

		passwordHash = &hash
	}

	return &Auth{
		ProfileID:      profileID,
		Provider:       provider,
		ProviderUserID: providerUserID,
		PasswordHash:   passwordHash,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}
