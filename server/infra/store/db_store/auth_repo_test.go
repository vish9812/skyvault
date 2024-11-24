package db_store

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	t.Run("create user", func(t *testing.T) {
		t.Parallel()
		testStore := newTestDBStore(t)
		authRepo := testStore.NewAuthRepo()
		createTestUser(t, authRepo)
	})

	t.Run("duplicate email id", func(t *testing.T) {
		t.Parallel()

		testStore := newTestDBStore(t)
		authRepo := testStore.NewAuthRepo()
		userA := createTestUser(t, authRepo)
		require.NotNil(t, userA, "failed to create userA")

		// Attempt to create another user with same email ID
		userB := newTestUser()
		userB.Email = userA.Email

		err := authRepo.CreateUser(context.Background(), userB)
		require.Error(t, err, "should have received error when creating userB")
	})
}

func TestGetUserByEmail(t *testing.T) {
	t.Parallel()
	testStore := newTestDBStore(t)
	authRepo := testStore.NewAuthRepo()

	userA := createTestUser(t, authRepo)
	userB := createTestUser(t, authRepo)

	user, err := authRepo.GetUserByEmail(context.Background(), userA.Email)
	require.Nil(t, err)
	require.EqualValues(t, userA.ID, user.ID, "userA ID mismatched")

	user, err = authRepo.GetUserByEmail(context.Background(), userB.Email)
	require.Nil(t, err)
	require.EqualValues(t, userB.ID, user.ID, "userB ID mismatched")

	_, err = authRepo.GetUserByEmail(context.Background(), "not-existed")
	require.ErrorIs(t, err, ErrNoRows, "error mismatched")
}
