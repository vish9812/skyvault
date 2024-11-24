package auth_app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"empty password", "", false},
		{"simple password", "password123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := hashPassword(tt.input)
			require.NoError(t, err)
			require.NotEmpty(t, got, "The hash should not be an empty string")
		})
	}
}
