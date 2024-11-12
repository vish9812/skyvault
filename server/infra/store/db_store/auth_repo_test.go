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

func TestGetUserByUsername(t *testing.T) {
	t.Parallel()
	testStore := newTestDBStore(t)
	authRepo := testStore.NewAuthRepo()

	userA := createTestUser(t, authRepo)
	userB := createTestUser(t, authRepo)

	user, err := authRepo.GetUserByUsername(context.Background(), userA.Username)
	require.Nil(t, err)
	require.EqualValues(t, userA.ID, user.ID, "userA ID mismatched")

	user, err = authRepo.GetUserByUsername(context.Background(), userB.Username)
	require.Nil(t, err)
	require.EqualValues(t, userB.ID, user.ID, "userB ID mismatched")

	_, err = authRepo.GetUserByUsername(context.Background(), "not-existed")
	require.ErrorIs(t, err, ErrNoRows, "error mismatched")
}
