package store_db

import (
	"context"
	"database/sql"

	"skyvault/internal/domain/profile"
	"skyvault/internal/infra/store_db/internal/gen_jet/skyvault/public/model"
	. "skyvault/internal/infra/store_db/internal/gen_jet/skyvault/public/table"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/jinzhu/copier"
)

var _ profile.Repo = (*ProfileRepo)(nil)

type ProfileRepo struct {
	store *store
}

func NewProfileRepo(db *store) *ProfileRepo {
	return &ProfileRepo{store: db}
}

func (r *ProfileRepo) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.store.db.BeginTx(ctx, nil)
}

func (r *ProfileRepo) WithTx(ctx context.Context, tx *sql.Tx) profile.Repo {
	return &ProfileRepo{store: r.store.WithTx(ctx, tx)}
}

func (r *ProfileRepo) Create(ctx context.Context, pro *profile.Profile) (*profile.Profile, error) {
	dbProfile := new(model.Profiles)
	err := copier.Copy(dbProfile, pro)
	if err != nil {
		return nil, err
	}

	stmt := Profiles.INSERT(
		Profiles.MutableColumns.Except(Profiles.CreatedAt, Profiles.UpdatedAt),
	).MODEL(dbProfile).RETURNING(Profiles.AllColumns)

	return runSelect[model.Profiles, profile.Profile](ctx, stmt, r.store.dbTx)
}

func (r *ProfileRepo) Get(ctx context.Context, id int64) (*profile.Profile, error) {
	stmt := SELECT(Profiles.AllColumns).
		FROM(Profiles).
		WHERE(Profiles.ID.EQ(Int64(id)))

	return runSelect[model.Profiles, profile.Profile](ctx, stmt, r.store.dbTx)
}

func (r *ProfileRepo) GetByEmail(ctx context.Context, email string) (*profile.Profile, error) {
	stmt := SELECT(Profiles.AllColumns).
		FROM(Profiles).
		WHERE(Profiles.Email.EQ(String(email)))

	return runSelect[model.Profiles, profile.Profile](ctx, stmt, r.store.dbTx)
}
