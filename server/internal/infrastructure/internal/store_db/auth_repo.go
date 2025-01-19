package store_db

import (
	"context"
	"database/sql"

	"skyvault/internal/domain/auth"
	"skyvault/internal/infrastructure/internal/store_db/internal/gen_jet/skyvault/public/model"
	. "skyvault/internal/infrastructure/internal/store_db/internal/gen_jet/skyvault/public/table"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/jinzhu/copier"
)

var _ auth.Repo = (*AuthRepo)(nil)

type AuthRepo struct {
	store *Store
}

func NewAuthRepo(store *Store) *AuthRepo {
	return &AuthRepo{store: store}
}

func (r *AuthRepo) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.store.DB.BeginTx(ctx, nil)
}

func (r *AuthRepo) WithTx(ctx context.Context, tx *sql.Tx) auth.Repo {
	return &AuthRepo{store: r.store.WithTx(ctx, tx)}
}

func (r *AuthRepo) Create(ctx context.Context, au *auth.Auth) (*auth.Auth, error) {
	dbAuth := new(model.Auths)
	err := copier.Copy(dbAuth, au)
	if err != nil {
		return nil, err
	}

	stmt := Auths.INSERT(
		Auths.MutableColumns.Except(Auths.CreatedAt, Auths.UpdatedAt),
	).MODEL(dbAuth).RETURNING(Auths.AllColumns)

	return runSelect[model.Auths, auth.Auth](ctx, stmt, r.store.DBTX)
}

func (r *AuthRepo) Get(ctx context.Context, id int64) (*auth.Auth, error) {
	stmt := SELECT(Auths.AllColumns).
		FROM(Auths).
		WHERE(Auths.ID.EQ(Int(id)))

	return runSelect[model.Auths, auth.Auth](ctx, stmt, r.store.DBTX)
}

func (r *AuthRepo) GetByProfileID(ctx context.Context, id int64) (*auth.Auth, error) {
	stmt := SELECT(Auths.AllColumns).
		FROM(Auths).
		WHERE(Auths.ProfileID.EQ(Int(id)))

	return runSelect[model.Auths, auth.Auth](ctx, stmt, r.store.DBTX)
}