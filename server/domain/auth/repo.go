package auth

import (
	"context"
)

type Repo interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	// GetUserPassword(ctx context.Context, id ids_utils.ID) (string, error)
}
