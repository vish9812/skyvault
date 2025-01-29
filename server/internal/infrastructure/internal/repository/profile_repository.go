package repository

import (
	"context"
	"database/sql"

	"skyvault/internal/domain/profile"
	"skyvault/internal/infrastructure/internal/repository/internal/gen_jet/skyvault/public/model"
	. "skyvault/internal/infrastructure/internal/repository/internal/gen_jet/skyvault/public/table"
	"skyvault/pkg/apperror"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/jinzhu/copier"
)

var _ profile.Repository = (*ProfileRepository)(nil)

type ProfileRepository struct {
	repository *Repository
}

func NewProfileRepository(repo *Repository) *ProfileRepository {
	return &ProfileRepository{repository: repo}
}

func (r *ProfileRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.repository.db.BeginTx(ctx, nil)
}

func (r *ProfileRepository) WithTx(ctx context.Context, tx *sql.Tx) profile.Repository {
	return &ProfileRepository{repository: r.repository.WithTx(ctx, tx)}
}

func (r *ProfileRepository) Create(ctx context.Context, pro *profile.Profile) (*profile.Profile, error) {
	dbModel := new(model.Profile)
	err := copier.Copy(dbModel, pro)
	if err != nil {
		return nil, apperror.NewAppError(err, "repository.Create:copier.Copy")
	}

	stmt := Profile.INSERT(
		Profile.MutableColumns,
	).MODEL(dbModel).RETURNING(Profile.AllColumns)

	return runInsert[model.Profile, profile.Profile](ctx, stmt, r.repository.dbTx)
}

func (r *ProfileRepository) Get(ctx context.Context, id int64) (*profile.Profile, error) {
	stmt := SELECT(Profile.AllColumns).
		FROM(Profile).
		WHERE(Profile.ID.EQ(Int64(id)))

	return runSelect[model.Profile, profile.Profile](ctx, stmt, r.repository.dbTx)
}

func (r *ProfileRepository) GetByEmail(ctx context.Context, email string) (*profile.Profile, error) {
	stmt := SELECT(Profile.AllColumns).
		FROM(Profile).
		WHERE(Profile.Email.EQ(String(email)))

	return runSelect[model.Profile, profile.Profile](ctx, stmt, r.repository.dbTx)
}

func (r *ProfileRepository) Update(ctx context.Context, pro *profile.Profile) error {
	dbModel := new(model.Profile)
	err := copier.Copy(dbModel, pro)
	if err != nil {
		return apperror.NewAppError(err, "repository.Update:copier.Copy")
	}

	stmt := Profile.UPDATE(
		Profile.MutableColumns,
	).MODEL(dbModel).WHERE(Profile.ID.EQ(Int64(pro.ID)))

	return runUpdateOrDelete(ctx, stmt, r.repository.dbTx)
}

func (r *ProfileRepository) Delete(ctx context.Context, id int64) error {
	stmt := Profile.DELETE().
		WHERE(Profile.ID.EQ(Int64(id)))

	return runUpdateOrDelete(ctx, stmt, r.repository.dbTx)
}
