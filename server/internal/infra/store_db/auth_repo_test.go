package store_db

import (
	"context"
	"skyvault/internal/domain/auth"
	"skyvault/pkg/common"
	"skyvault/pkg/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		t.Parallel()
		store := newTestStore(t)
		authRepo := NewAuthRepo(store)
		profileRepo := NewProfileRepo(store)
		profileCreate(t, profileRepo, authRepo)
	})

	t.Run("duplicate auth provider", func(t *testing.T) {
		t.Parallel()
		store := newTestStore(t)
		authRepo := NewAuthRepo(store)
		pro, authA := profileCreate(t, NewProfileRepo(store), authRepo)

		// Create another Auth with the same provider and providerUserID
		authB := authRandom(pro.ID, pro.Email)
		authB.Provider = authA.Provider
		authB.ProviderUserID = authA.ProviderUserID
		_, err := authRepo.Create(context.Background(), authB)
		require.ErrorIs(t, err, common.ErrDBUniqueConstraint, "expected error for duplicate provider details of userB")

		// Create another Auth with the same provider and different providerUserID
		userC := authRandom(pro.ID, pro.Email)
		userC.Provider = authA.Provider
		_, err = authRepo.Create(context.Background(), userC)
		require.NoError(t, err, "expected no error for different providerUserID of userC")

		// Create another Auth with the same providerUserID and different provider
		userD := authRandom(pro.ID, pro.Email)
		userD.Provider = utils.RandomItemExcept(auth.Providers(), authA.Provider)
		userD.ProviderUserID = authA.ProviderUserID
		_, err = authRepo.Create(context.Background(), userD)
		require.NoError(t, err, "expected no error for different provider of userD")
	})
}
