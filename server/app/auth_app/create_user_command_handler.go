package auth_app

import (
	"context"
	"skyvault/domain/auth"
	"skyvault/infra/store"
)

type ICreateUserCommandHandler interface {
	Handle(ctx context.Context, cmd *CreateUserCommand) (*auth.User, error)
}

type CreateUserCommandHandler struct {
	Store *store.Store
}

func (h *CreateUserCommandHandler) Handle(ctx context.Context, cmd *CreateUserCommand) (*auth.User, error) {
	user := auth.NewUser()
	user.Email = cmd.Email
	user.FirstName = cmd.FirstName
	user.LastName = cmd.LastName
	hashedPassword, err := hashPassword(cmd.Password)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = hashedPassword

	err = h.Store.NewAuthRepo().CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
