package repository

import (
	"context"
	"database/sql"

	"skyvault/internal/domain/media"
	"skyvault/internal/infrastructure/internal/repository/internal/gen_jet/skyvault/public/model"
	. "skyvault/internal/infrastructure/internal/repository/internal/gen_jet/skyvault/public/table"
	"skyvault/pkg/apperror"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/jinzhu/copier"
)

var _ media.Repository = (*MediaRepository)(nil)

type MediaRepository struct {
	repository *Repository
}

func NewMediaRepository(repo *Repository) *MediaRepository {
	return &MediaRepository{repository: repo}
}

func (r *MediaRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.repository.db.BeginTx(ctx, nil)
}

func (r *MediaRepository) WithTx(ctx context.Context, tx *sql.Tx) media.Repository {
	return &MediaRepository{repository: r.repository.WithTx(ctx, tx)}
}

//--------------------------------
// File
//--------------------------------

func (r *MediaRepository) CreateFileInfo(ctx context.Context, info *media.FileInfo) (*media.FileInfo, error) {
	dbModel := new(model.FileInfo)
	err := copier.Copy(dbModel, info)
	if err != nil {
		return nil, apperror.NewAppError(err, "repository.CreateFileInfo:copier.Copy")
	}

	stmt := FileInfo.INSERT(
		FileInfo.MutableColumns,
	).MODEL(dbModel).RETURNING(FileInfo.AllColumns)

	return runInsert[model.FileInfo, media.FileInfo](ctx, stmt, r.repository.dbTx)
}

func (r *MediaRepository) GetFileInfo(ctx context.Context, fileID int64) (*media.FileInfo, error) {
	stmt := SELECT(FileInfo.AllColumns).
		FROM(FileInfo).
		WHERE(FileInfo.ID.EQ(Int64(fileID)).
			AND(FileInfo.TrashedAt.IS_NULL()),
		)

	return runSelect[model.FileInfo, media.FileInfo](ctx, stmt, r.repository.dbTx)
}

func (r *MediaRepository) GetFilesInfo(ctx context.Context, ownerID int64, folderID *int64) ([]*media.FileInfo, error) {
	var folderCond BoolExpression
	if folderID == nil {
		// Get all files in the root folder of the owner.
		folderCond = FileInfo.FolderID.IS_NULL().AND(FileInfo.OwnerID.EQ(Int64(ownerID)))
	} else {
		// Get all files in the specified folder.
		folderCond = FileInfo.FolderID.EQ(Int64(*folderID))
	}

	stmt := SELECT(FileInfo.AllColumns).
		FROM(FileInfo).
		WHERE(
			folderCond.
				AND(FileInfo.TrashedAt.IS_NULL()),
		)

	return runSelectSlice[model.FileInfo, media.FileInfo](ctx, stmt, r.repository.dbTx)
}

func (r *MediaRepository) UpdateFileInfo(ctx context.Context, info *media.FileInfo) error {
	stmt := FileInfo.UPDATE(FileInfo.MutableColumns).
		MODEL(info).
		WHERE(FileInfo.ID.EQ(Int64(info.ID)))

	return runUpdateOrDelete(ctx, stmt, r.repository.dbTx)
}

func (r *MediaRepository) DeleteFileInfo(ctx context.Context, fileID int64) error {
	stmt := FileInfo.DELETE().
		WHERE(FileInfo.ID.EQ(Int64(fileID)))

	return runUpdateOrDelete(ctx, stmt, r.repository.dbTx)
}

//--------------------------------
// Folder
//--------------------------------

func (r *MediaRepository) CreateFolderInfo(ctx context.Context, info *media.FolderInfo) (*media.FolderInfo, error) {
	dbModel := new(model.FolderInfo)
	err := copier.Copy(dbModel, info)
	if err != nil {
		return nil, apperror.NewAppError(err, "repository.CreateFolderInfo:copier.Copy")
	}

	stmt := FolderInfo.INSERT(
		FolderInfo.MutableColumns,
	).MODEL(dbModel).RETURNING(FolderInfo.AllColumns)

	return runInsert[model.FolderInfo, media.FolderInfo](ctx, stmt, r.repository.dbTx)
}

func (r *MediaRepository) GetFolderInfo(ctx context.Context, folderID int64) (*media.FolderInfo, error) {
	stmt := SELECT(FolderInfo.AllColumns).
		FROM(FolderInfo).
		WHERE(FolderInfo.ID.EQ(Int64(folderID)).AND(FolderInfo.TrashedAt.IS_NULL()))

	return runSelect[model.FolderInfo, media.FolderInfo](ctx, stmt, r.repository.dbTx)
}

func (r *MediaRepository) GetFoldersInfo(ctx context.Context, ownerID int64, parentFolderID *int64) ([]*media.FolderInfo, error) {
	var parentFolderCond BoolExpression
	if parentFolderID == nil {
		// Get all folders in the root folder of the owner.
		parentFolderCond = FolderInfo.ParentFolderID.IS_NULL().AND(FolderInfo.OwnerID.EQ(Int64(ownerID)))
	} else {
		// Get all folders in the specified parent folder.
		parentFolderCond = FolderInfo.ParentFolderID.EQ(Int64(*parentFolderID))
	}

	stmt := SELECT(FolderInfo.AllColumns).
		FROM(FolderInfo).
		WHERE(
			parentFolderCond.
				AND(FolderInfo.TrashedAt.IS_NULL()),
		)

	return runSelectSlice[model.FolderInfo, media.FolderInfo](ctx, stmt, r.repository.dbTx)
}

func (r *MediaRepository) UpdateFolderInfo(ctx context.Context, info *media.FolderInfo) error {
	stmt := FolderInfo.UPDATE(FolderInfo.MutableColumns).
		MODEL(info).
		WHERE(FolderInfo.ID.EQ(Int64(info.ID)))

	return runUpdateOrDelete(ctx, stmt, r.repository.dbTx)
}

func (r *MediaRepository) DeleteFolderInfo(ctx context.Context, folderID int64) error {
	stmt := FolderInfo.DELETE().
		WHERE(FolderInfo.ID.EQ(Int64(folderID)))

	return runUpdateOrDelete(ctx, stmt, r.repository.dbTx)
}
