package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAuth(t *testing.T) {
	validPassword := "password123"
	validEmail := "test@example.com"
	validProfileID := "1"

	tests := []struct {
		name           string
		profileID      string
		provider       Provider
		providerUserID string
		password       *string
		wantErr        bool
	}{
		{
			name:           "valid email auth",
			profileID:      validProfileID,
			provider:       ProviderEmail,
			providerUserID: validEmail,
			password:       &validPassword,
		},
		{
			name:           "valid oidc auth",
			profileID:      validProfileID,
			provider:       ProviderOIDC,
			providerUserID: "oidc-user-123",
			password:       nil,
		},
		{
			name:           "valid ldap auth",
			profileID:      validProfileID,
			provider:       ProviderLDAP,
			providerUserID: "ldap-user-123",
			password:       nil,
		},
		{
			name:           "email auth missing password",
			profileID:      validProfileID,
			provider:       ProviderEmail,
			providerUserID: validEmail,
			password:       nil,
			wantErr:        true,
		},
		{
			name:           "email auth invalid email",
			profileID:      validProfileID,
			provider:       ProviderEmail,
			providerUserID: "invalid-email",
			password:       &validPassword,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth, err := NewAuth(tt.profileID, tt.provider, tt.providerUserID, tt.password)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, auth)
			assert.Equal(t, tt.profileID, auth.ProfileID)
			assert.Equal(t, tt.provider, auth.Provider)
			assert.Equal(t, tt.providerUserID, auth.ProviderUserID)

			// For email provider, verify password was hashed
			if tt.provider == ProviderEmail {
				require.NotNil(t, auth.PasswordHash)
				assert.NotEqual(t, *tt.password, *auth.PasswordHash)
			} else {
				require.Nil(t, auth.PasswordHash)
			}
		})
	}
}

func TestAuth_ValidateAccess(t *testing.T) {
	validProfileID := "1"

	auth := &Auth{
		ProfileID: validProfileID,
	}

	tests := []struct {
		name         string
		accessedByID string
		wantErr      bool
	}{
		{
			name:         "valid access",
			accessedByID: validProfileID,
		},
		{
			name:         "invalid access",
			accessedByID: "2",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := auth.ValidateAccess(tt.accessedByID)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
