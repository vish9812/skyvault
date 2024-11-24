package auth_app

import (
	"context"
	"skyvault/common/utils"
	"skyvault/domain"
	"skyvault/domain/auth"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var _ domain.IStore = &MockStore{}
var _ auth.Repo = &MockAuthRepo{}

// MockStore is a mock implementation of the domain.IStore interface for testing.
type MockStore struct {
	mock.Mock
}

func (m *MockStore) NewAuthRepo() auth.Repo {
	args := m.Called()
	return args.Get(0).(auth.Repo)
}

// MockAuthRepo is a mock implementation of the auth.Repo interface for testing.
type MockAuthRepo struct {
	mock.Mock
}

// CreateUser implements auth.Repo.
func (m *MockAuthRepo) CreateUser(ctx context.Context, user *auth.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// GetUserByEmail implements auth.Repo.
func (m *MockAuthRepo) GetUserByEmail(ctx context.Context, email string) (*auth.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.User), args.Error(1)
}

func TestHandle(t *testing.T) {
	randName := utils.RandomName()
	cmd := &SignUpCommand{
		FirstName: randName,
		LastName:  randName,
		Email:     utils.RandomEmail(),
		Password:  "pass",
	}

	// Create mock instances of the interfaces
	mockStore := new(MockStore)
	mockRepo := new(MockAuthRepo)

	// Define expected calls and their return values
	mockStore.On("NewAuthRepo").Return(mockRepo)
	mockRepo.On("CreateUser", mock.Anything, mock.Anything).Return(nil)

	// Initialize the handler
	authApp := &AuthApp{Store: mockStore}
	handler := authApp.NewSignUpCommandHandler()

	// Call the function being tested
	actualUser, err := handler.Handle(context.Background(), cmd)
	require.NoError(t, err)
	assert.Equal(t, cmd.Email, actualUser.Email)
	assert.Equal(t, cmd.FirstName, actualUser.FirstName)
	assert.Equal(t, cmd.LastName, actualUser.LastName)
	assert.NotEmpty(t, actualUser.ID)
	assert.True(t, isValidPassword(actualUser.PasswordHash, cmd.Password))

	// Assert that the expected calls were made
	mockRepo.AssertExpectations(t)
}
