package store

import (
	"skyvault/common/utils"
	"skyvault/domain/auth"
	"testing"

	"context"

	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	testStore := newTestStore(t)
	authRepo := testStore.NewAuthRepo()

	user := auth.NewUser()
	user.Email = utils.RandomEmail()
	randStr := utils.RandomName()
	user.Username = randStr
	user.FirstName = randStr
	user.LastName = randStr
	user.PasswordHash = randStr

	err := authRepo.CreateUser(context.Background(), user)
	require.Nil(t, err)
}
