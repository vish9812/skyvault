package auth_app

import (
	"context"
	"skyvault/common/utils/utils_api"
	"skyvault/domain/auth"
	"skyvault/infra/store"
)

type ISignUpCommandHandler interface {
	Handle(ctx context.Context, cmd *SignUpCommand) (*auth.User, error)
}

type SignUpCommandHandler struct {
	Store *store.Store
	AuthRepo auth.Repo
}

func (h *SignUpCommandHandler) Handle(ctx context.Context, cmd *SignUpCommand) (*auth.User, error) {
	user := auth.NewUser()
	user.Email = cmd.Email
	user.FirstName = cmd.FirstName
	user.LastName = cmd.LastName
	hashedPassword, err := hashPassword(cmd.Password)
	if err != nil {
		return nil, utils_api.APIError{
			Message: "failed to hash password",
			Err:     err,
		}
	}
	user.PasswordHash = hashedPassword

	err = h.Store.NewAuthRepo().CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
