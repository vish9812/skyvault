//lint:file-ignore ST1001 Using dot import to make SQL queries more readable
package repository

import (
	"context"
	"database/sql"

	"skyvault/internal/domain/auth"
	"skyvault/internal/infrastructure/internal/repository/internal/gen_jet/skyvault/public/model"
	. "skyvault/internal/infrastructure/internal/repository/internal/gen_jet/skyvault/public/table"
	"skyvault/pkg/apperror"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/jinzhu/copier"
)

var _ auth.Repository = (*AuthRepository)(nil)

type AuthRepository struct {
	repository *Repository
}

func NewAuthRepository(repo *Repository) *AuthRepository {
	return &AuthRepository{repository: repo}
}

func (r *AuthRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.repository.db.BeginTx(ctx, nil)
}

func (r *AuthRepository) WithTx(ctx context.Context, tx *sql.Tx) auth.Repository {
	return &AuthRepository{repository: r.repository.withTx(ctx, tx)}
}

func (r *AuthRepository) Create(ctx context.Context, au *auth.Auth) (*auth.Auth, error) {
	dbModel := new(model.Auth)
	err := copier.Copy(dbModel, au)
	if err != nil {
		return nil, apperror.NewAppError(err, "repository.Create:copier.Copy")
	}

	stmt := Auth.INSERT(
		Auth.MutableColumns,
	).MODEL(dbModel).RETURNING(Auth.AllColumns)

	return runInsert[model.Auth, auth.Auth](ctx, stmt, r.repository.dbTx)
}

func (r *AuthRepository) Get(ctx context.Context, id int64) (*auth.Auth, error) {
	stmt := SELECT(Auth.AllColumns).
		FROM(Auth).
		WHERE(Auth.ID.EQ(Int(id)))

	return runSelect[model.Auth, auth.Auth](ctx, stmt, r.repository.dbTx)
}

func (r *AuthRepository) GetByProfileID(ctx context.Context, id int64) ([]*auth.Auth, error) {
	stmt := SELECT(Auth.AllColumns).
		FROM(Auth).
		WHERE(Auth.ProfileID.EQ(Int(id)))

	return runSelectSliceAll[model.Auth, auth.Auth](ctx, stmt, r.repository.dbTx)
}

func (r *AuthRepository) GetByProvider(ctx context.Context, provider auth.Provider, providerUserID string) (*auth.Auth, error) {
	stmt := SELECT(Auth.AllColumns).
		FROM(Auth).
		WHERE(
			Auth.Provider.EQ(String(string(provider))).AND(
				Auth.ProviderUserID.EQ(String(providerUserID)),
			),
		)

	return runSelect[model.Auth, auth.Auth](ctx, stmt, r.repository.dbTx)
}

func (r *AuthRepository) Update(ctx context.Context, au *auth.Auth) error {
	dbModel := new(model.Auth)
	err := copier.Copy(dbModel, au)
	if err != nil {
		return apperror.NewAppError(err, "repository.Update:copier.Copy")
	}

	stmt := Auth.UPDATE(Auth.MutableColumns).
		MODEL(dbModel).
		WHERE(Auth.ID.EQ(Int(au.ID)))

	return runUpdateOrDelete(ctx, stmt, r.repository.dbTx)
}

func (r *AuthRepository) Delete(ctx context.Context, id int64) error {
	stmt := Auth.DELETE().
		WHERE(Auth.ID.EQ(Int(id)))

	return runUpdateOrDelete(ctx, stmt, r.repository.dbTx)
}
