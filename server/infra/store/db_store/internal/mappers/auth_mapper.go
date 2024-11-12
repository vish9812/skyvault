package mappers

import (
	"skyvault/common/utils"
	"skyvault/domain/auth"
	"skyvault/infra/store/db_store/internal/gen_jet/skyvault/public/model"
)

func DomainUserToDBUser(user *auth.User) *model.Users {
	return &model.Users{
		ID:           user.ID.ToUUID(),
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
	}
}

func DBUserToDomainUser(dbUser *model.Users) *auth.User {
	return &auth.User{
		ID:        utils.ToID(dbUser.ID),
		FirstName: dbUser.FirstName,
		LastName:  dbUser.LastName,
		Email:     dbUser.Email,
		Username:  dbUser.Username,
	}
}
