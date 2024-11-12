package auth

import (
	"context"
)

type Repo interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	// GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	// GetUserPassword(ctx context.Context, id ids_utils.ID) (string, error)
}
