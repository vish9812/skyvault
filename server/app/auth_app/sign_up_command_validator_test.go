package auth_app

import (
	"context"
	"skyvault/domain/auth"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockSignUpCommandHandler struct {
	mock.Mock
}

func (m *MockSignUpCommandHandler) Handle(ctx context.Context, cmd *SignUpCommand) (*auth.User, error) {
	args := m.Called(ctx, cmd)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.User), args.Error(1)
}

func TestHandle_Success(t *testing.T) {
	t.Parallel()
	// Create a mock implementation of the ISignUpCommandHandler interface
	handlerMock := new(MockSignUpCommandHandler)

	// Initialize the validator with the handler mock
	validator := &SignUpCommandValidator{Handler: handlerMock}

	// Create a sample SignUpCommand
	cmd := &SignUpCommand{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "password123",
	}

	// Define what the handler mock should do when its Handle method is called
	handlerMock.On("Handle", mock.Anything, cmd).Return(auth.NewUser(), nil)

	// Call the Handle method of the validator under test
	user, err := validator.Handle(context.Background(), cmd)
	require.NoError(t, err)
	assert.NotNil(t, user)

	// Assert that the mock's Handle method was called as expected
	handlerMock.AssertExpectations(t)
}

func TestHandle_MissingFirstName(t *testing.T) {
	t.Parallel()
	validator := &SignUpCommandValidator{}
	cmd := &SignUpCommand{
		LastName: "Doe",
		Email:    "john.doe@example.com",
		Password: "password123",
	}

	_, err := validator.Handle(context.Background(), cmd)
	assert.ErrorIs(t, err, errFirstNameRequired)
}

func TestHandle_MissingLastName(t *testing.T) {
	t.Parallel()
	validator := &SignUpCommandValidator{}
	cmd := &SignUpCommand{
		FirstName: "John",
		Email:     "john.doe@example.com",
		Password:  "password123",
	}

	_, err := validator.Handle(context.Background(), cmd)
	assert.ErrorIs(t, err, errLastNameRequired)
}

func TestHandle_MissingEmail(t *testing.T) {
	t.Parallel()
	validator := &SignUpCommandValidator{}
	cmd := &SignUpCommand{
		FirstName: "John",
		LastName:  "Doe",
		Password:  "password123",
	}

	_, err := validator.Handle(context.Background(), cmd)
	assert.ErrorIs(t, err, errEmailRequired)
}

func TestHandle_InvalidEmail(t *testing.T) {
	t.Parallel()
	validator := &SignUpCommandValidator{}
	cmd := &SignUpCommand{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "invalid-email",
		Password:  "password123",
	}

	_, err := validator.Handle(context.Background(), cmd)
	assert.ErrorIs(t, err, errEmailInvalid)
}

func TestHandle_MissingPassword(t *testing.T) {
	t.Parallel()
	validator := &SignUpCommandValidator{}
	cmd := &SignUpCommand{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}

	_, err := validator.Handle(context.Background(), cmd)
	assert.ErrorIs(t, err, errPasswordRequired)
}

func TestHandle_PasswordLen(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{"short password", strings.Repeat("1", passwordMinLen-1), errPasswordMinLen},
		{"long password", strings.Repeat("1", passwordMaxLen+1), errPasswordMaxLen},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			validator := &SignUpCommandValidator{}
			cmd := &SignUpCommand{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
				Password:  tt.input,
			}

			_, err := validator.Handle(context.Background(), cmd)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
