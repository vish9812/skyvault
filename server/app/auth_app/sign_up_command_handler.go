package auth_app

import (
	"context"
	"net/http"
	"skyvault/common/utils/utils_api"
	"skyvault/domain"
	"skyvault/domain/auth"
)

type ISignUpCommandHandler interface {
	Handle(ctx context.Context, cmd *SignUpCommand) (*auth.User, error)
}

type SignUpCommandHandler struct {
	Store domain.IStore
}

func (h *SignUpCommandHandler) Handle(ctx context.Context, cmd *SignUpCommand) (*auth.User, error) {
	user := auth.NewUser()
	user.Email = cmd.Email
	user.FirstName = cmd.FirstName
	user.LastName = cmd.LastName
	hashedPassword, err := hashPassword(cmd.Password)
	if err != nil {
		return nil, utils_api.NewAPIError(http.StatusInternalServerError, "failed to hash password", "", err)
	}
	user.PasswordHash = hashedPassword
	
	err = h.Store.NewAuthRepo().CreateUser(ctx, user)
	if err != nil {
		return nil, utils_api.NewAPIError(http.StatusInternalServerError, "failed to create user", "", err)
	}

	return user, nil
}
