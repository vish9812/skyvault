package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"skyvault/internal/api/helper/dtos"
	"skyvault/internal/domain/auth"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAuthFlow tests the complete authentication flow including
// sign up and sign in with different providers
func TestAuthFlow(t *testing.T) {
	t.Parallel()
	env := setupTestEnv(t)

	testCases := []struct {
		name     string
		provider auth.Provider
		email    string
		password string
	}{
		{
			name:     "Email Provider",
			provider: auth.ProviderEmail,
			email:    "test@email.com",
			password: "password123",
		},
		// {
		// 	name:     "Google Provider",
		// 	provider: auth.ProviderOIDC,
		// 	email:    "test@oidc.com",
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Test Sign Up
			signUpReq := map[string]interface{}{
				"email":    tc.email,
				"fullName": "Test User",
				"provider": string(tc.provider),
			}

			if tc.provider == auth.ProviderEmail {
				signUpReq["password"] = tc.password
			}

			jsonBody, err := json.Marshal(signUpReq)
			require.NoError(t, err, "should marshal sign up request")

			req, err := http.NewRequest(http.MethodPost, "/api/v1/pub/auth/sign-up", bytes.NewBuffer(jsonBody))
			require.NoError(t, err, "should create new request for sign up")
			req.Header.Set("Content-Type", "application/json")

			resp := executeRequest(t, env, req)
			require.Equal(t, http.StatusCreated, resp.Code, "should return status created for sign up")

			var signUpRes dtos.SignUp
			err = json.NewDecoder(resp.Body).Decode(&signUpRes)
			require.NoError(t, err, "should decode sign up response")

			// Verify sign up response
			assert.NotEmpty(t, signUpRes.Token, "should return token")
			assert.NotNil(t, signUpRes.Profile, "should return profile")
			assert.Equal(t, tc.email, signUpRes.Profile.Email, "should return correct email")
			assert.Equal(t, "Test User", signUpRes.Profile.FullName, "should return correct full name")

			// Test Sign In
			signInReq := map[string]interface{}{
				"provider":       string(tc.provider),
				"providerUserId": tc.email,
			}

			if tc.provider == auth.ProviderEmail {
				signInReq["password"] = tc.password
			}

			jsonBody, err = json.Marshal(signInReq)
			require.NoError(t, err, "should marshal sign in request")

			req, err = http.NewRequest(http.MethodPost, "/api/v1/pub/auth/sign-in", bytes.NewBuffer(jsonBody))
			require.NoError(t, err, "should create new request for sign in")
			req.Header.Set("Content-Type", "application/json")

			resp = executeRequest(t, env, req)
			require.Equal(t, http.StatusOK, resp.Code, "should return status ok for sign in")

			var signInRes dtos.SignUp
			err = json.NewDecoder(resp.Body).Decode(&signInRes)
			require.NoError(t, err, "should decode sign in response")

			// Verify sign in response
			assert.NotEmpty(t, signInRes.Token, "should return token")
			assert.NotNil(t, signInRes.Profile, "should return profile")
			assert.Equal(t, tc.email, signInRes.Profile.Email, "should return correct email")
			assert.Equal(t, "Test User", signInRes.Profile.FullName, "should return correct full name")
		})
	}
}

// TestAuthErrors tests various error scenarios in authentication
func TestAuthErrors(t *testing.T) {
	t.Parallel()
	env := setupTestEnv(t)

	t.Run("Duplicate Sign Up", func(t *testing.T) {
		t.Parallel()

		// First sign up
		signUpReq := map[string]interface{}{
			"email":    "duplicate@email.com",
			"fullName": "Test User",
			"provider": string(auth.ProviderEmail),
			"password": "password123",
		}

		jsonBody, err := json.Marshal(signUpReq)
		require.NoError(t, err, "should marshal first sign up request")

		req, err := http.NewRequest(http.MethodPost, "/api/v1/pub/auth/sign-up", bytes.NewBuffer(jsonBody))
		require.NoError(t, err, "should create new request for first sign up")
		req.Header.Set("Content-Type", "application/json")

		resp := executeRequest(t, env, req)
		require.Equal(t, http.StatusCreated, resp.Code, "should return status created for first sign up")

		// Duplicate sign up
		req, err = http.NewRequest(http.MethodPost, "/api/v1/pub/auth/sign-up", bytes.NewBuffer(jsonBody))
		require.NoError(t, err, "should create new request for duplicate sign up")
		req.Header.Set("Content-Type", "application/json")

		resp = executeRequest(t, env, req)
		assert.Equal(t, http.StatusConflict, resp.Code, "should return conflict for duplicate sign up")
	})

	t.Run("Invalid Credentials Sign In", func(t *testing.T) {
		t.Parallel()

		signInReq := map[string]interface{}{
			"provider":       string(auth.ProviderEmail),
			"providerUserId": "nonexistent@email.com",
			"password":       "wrongpassword",
		}

		jsonBody, err := json.Marshal(signInReq)
		require.NoError(t, err, "should marshal sign in request")

		req, err := http.NewRequest(http.MethodPost, "/api/v1/pub/auth/sign-in", bytes.NewBuffer(jsonBody))
		require.NoError(t, err, "should create new request for sign in")
		req.Header.Set("Content-Type", "application/json")

		resp := executeRequest(t, env, req)
		assert.Equal(t, http.StatusNotFound, resp.Code, "should return NotFound for invalid credentials")
	})
}
