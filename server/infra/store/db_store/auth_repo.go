package db_store

import (
	"context"
	"errors"
	"skyvault/domain/auth"
	"time"

	"skyvault/infra/store/db_store/internal/gen_jet/skyvault/public/model"
	. "skyvault/infra/store/db_store/internal/gen_jet/skyvault/public/table"
	"skyvault/infra/store/db_store/internal/mappers"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
)

type AuthRepo struct {
	DB *DBStore
}

func (r *AuthRepo) CreateUser(ctx context.Context, user *auth.User) error {
	dbUser := mappers.DomainUserToDBUser(user)
	t := time.Now().UTC()
	dbUser.CreatedAt, dbUser.UpdatedAt = t, t

	stmt := Users.INSERT(Users.AllColumns).MODEL(dbUser)

	stdDB := r.DB.openStdDB()
	defer r.DB.closeStdDB(stdDB)

	_, err := stmt.ExecContext(ctx, stdDB)

	return err
}

func (r *AuthRepo) GetUserByUsername(ctx context.Context, username string) (*auth.User, error) {
	stmt := SELECT(Users.AllColumns.Except(Users.PasswordHash)).
		FROM(Users).
		WHERE(Users.Username.EQ(String(username)))

	stdDB := r.DB.openStdDB()
	defer r.DB.closeStdDB(stdDB)

	var dbUser model.Users
	err := stmt.QueryContext(ctx, stdDB, &dbUser)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, ErrNoRows
		}
		return nil, err
	}

	return mappers.DBUserToDomainUser(&dbUser), nil
}
