package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "mypassword123",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  false, // bcrypt allows empty passwords
		},
		{
			name:     "long password",
			password: strings.Repeat("a", 72), // bcrypt max length
			wantErr:  false,
		},
		{
			name:     "too long password",
			password: strings.Repeat("a", 73), // bcrypt max length + 1
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			hash, err := HashPassword(tt.password)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.True(t, strings.HasPrefix(hash, "$2a$"), "hash should start with $2a$")
		})
	}
}

func TestIsValidPassword(t *testing.T) {
	t.Parallel()
	// Create a known password and hash for testing
	password := "mypassword123"
	hash, err := HashPassword(password)
	require.NoError(t, err, "Failed to create test hash")
	wrongPassword := "wrongpassword"
	wrongPasswordHash, err := HashPassword(wrongPassword)
	require.NoError(t, err, "Failed to create test wrong hash")

	tests := []struct {
		name        string
		hash        string
		password    string
		want        bool
		wantErr     bool
	}{
		{
			name:     "valid password",
			hash:     hash,
			password: password,
			want:     true,
			wantErr:  false,
		},
		{
			name:     "invalid password",
			hash:     hash,
			password: wrongPassword,
			want:     false,
			wantErr:  false,
		},
		{
			name:     "wrong hash",
			hash:     wrongPasswordHash,
			password: password,
			want:     false,
			wantErr:  false,
		},
		{
			name:        "invalid hash format",
			hash:        strings.Repeat("a", 100),
			password:    password,
			want:        false,
			wantErr:     true,
		},
		{
			name:     "empty password",
			hash:     hash,
			password: "",
			want:     false,
			wantErr:  false,
		},
		{
			name:     "empty hash",
			hash:     "",
			password: password,
			want:     false,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := IsValidPassword(tt.hash, tt.password)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestHashPassword_Cost(t *testing.T) {
	t.Parallel()
	password := "mypassword123"
	hash, err := HashPassword(password)
	require.NoError(t, err, "HashPassword() failed")

	cost, err := bcrypt.Cost([]byte(hash))
	require.NoError(t, err, "bcrypt.Cost() failed")
	require.Equal(t, bcrypt.DefaultCost, cost)
}
