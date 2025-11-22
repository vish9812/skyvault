//lint:file-ignore ST1001 Using dot import to make SQL queries more readable
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
	return &ProfileRepository{repository: r.repository.withTx(ctx, tx)}
}

func (r *ProfileRepository) Create(ctx context.Context, pro *profile.Profile) (*profile.Profile, error) {
	dbModel := new(model.Profile)
	err := copier.Copy(dbModel, pro)
	if err != nil {
		return nil, apperror.NewAppError(err, "repository.Create:copier.Copy")
	}

	stmt := Profile.INSERT(
		Profile.AllColumns,
	).MODEL(dbModel).RETURNING(Profile.AllColumns)

	return runInsert[model.Profile, profile.Profile](ctx, stmt, r.repository.dbTx)
}

func (r *ProfileRepository) Get(ctx context.Context, id string) (*profile.Profile, error) {
	stmt := SELECT(Profile.AllColumns).
		FROM(Profile).
		WHERE(Profile.ID.EQ(UUID(UUIDStr(id))))

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
	).MODEL(dbModel).WHERE(Profile.ID.EQ(UUID(UUIDStr(pro.ID))))

	return runUpdateOrDelete(ctx, stmt, r.repository.dbTx)
}

func (r *ProfileRepository) Delete(ctx context.Context, id string) error {
	stmt := Profile.DELETE().
		WHERE(Profile.ID.EQ(UUID(UUIDStr(id))))

	return runUpdateOrDelete(ctx, stmt, r.repository.dbTx)
}

func (r *ProfileRepository) AtomicAllocateStorage(ctx context.Context, profileID string, bytes int64) error {
	stmt := Profile.UPDATE(Profile.StorageUsed).
		SET(Profile.StorageUsed.SET(Profile.StorageUsed.ADD(Int64(bytes)))).
		WHERE(
			Profile.ID.EQ(UUID(UUIDStr(profileID))).
				AND(Profile.StorageUsed.ADD(Int64(bytes)).LT_EQ(Profile.StorageQuota)),
		)

	result, err := stmt.ExecContext(ctx, r.repository.dbTx)
	if err != nil {
		return apperror.NewAppError(err, "repository.AtomicAllocateStorage:ExecContext")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return apperror.NewAppError(err, "repository.AtomicAllocateStorage:RowsAffected")
	}

	if rowsAffected == 0 {
		// Either profile doesn't exist or quota would be exceeded
		// Try to fetch profile to provide better error context
		pro, err := r.Get(ctx, profileID)
		if err != nil {
			// Profile doesn't exist
			return apperror.NewAppError(err, "repository.AtomicAllocateStorage:Get")
		}

		// Quota would be exceeded
		return apperror.NewAppError(apperror.ErrStorageQuotaExceeded, "repository.AtomicAllocateStorage:QuotaCheck").
			WithMetadata("profile_id", profileID).
			WithMetadata("requested_bytes", bytes).
			WithMetadata("storage_quota_bytes", pro.StorageQuota).
			WithMetadata("storage_used_bytes", pro.StorageUsed)
	}

	return nil
}

func (r *ProfileRepository) IncrementStorageUsage(ctx context.Context, profileID string, bytes int64) error {
	stmt := Profile.UPDATE(Profile.StorageUsed).
		SET(Profile.StorageUsed.SET(Profile.StorageUsed.ADD(Int64(bytes)))).
		WHERE(Profile.ID.EQ(UUID(UUIDStr(profileID))))

	return runUpdateOrDelete(ctx, stmt, r.repository.dbTx)
}

func (r *ProfileRepository) DecrementStorageUsage(ctx context.Context, profileID string, bytes int64) error {
	stmt := Profile.UPDATE(Profile.StorageUsed).
		SET(Profile.StorageUsed.SET(Profile.StorageUsed.SUB(Int64(bytes)))).
		WHERE(Profile.ID.EQ(UUID(UUIDStr(profileID))))

	return runUpdateOrDelete(ctx, stmt, r.repository.dbTx)
}
