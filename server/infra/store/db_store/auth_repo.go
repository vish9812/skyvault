package db_store

import (
	"context"
	"skyvault/domain/auth"
	"time"

	"skyvault/infra/store/db_store/internal/gen_jet/skyvault/public/model"
	. "skyvault/infra/store/db_store/internal/gen_jet/skyvault/public/table"
	// . "github.com/go-jet/jet/v2/postgres"
)

type AuthRepo struct {
	DB *DBStore
}

func domainUserToDBUser(user *auth.User) *model.Users {
	t := time.Now().UTC()

	return &model.Users{
		ID:           user.ID.ToUUID(),
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		CreatedAt:    t,
		UpdatedAt:    t,
	}
}

func (r *AuthRepo) CreateUser(ctx context.Context, user *auth.User) error {
	dbUser := domainUserToDBUser(user)

	stmt := Users.INSERT(Users.AllColumns).MODEL(dbUser)

	query, args := stmt.Sql()
	_, err := r.DB.Exec(ctx, query, args...)

	return err
}

// func (r *AuthRepo) Get(ctx context.Context, id int) *auth.User {
// 	entity := &UserEntity{
// 		ID:       id,
// 		Username: "me",
// 		Password: "pass",
// 	}
// 	return entity.ToModel()
// }
