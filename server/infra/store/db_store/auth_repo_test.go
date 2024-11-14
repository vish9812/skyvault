package db_store

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()
	testStore := newTestDBStore(t)
	authRepo := testStore.NewAuthRepo()
	createTestUser(t, authRepo)
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
