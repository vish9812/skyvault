package store_db

import (
	"context"
	"database/sql"

	"skyvault/internal/domain/media"
	"skyvault/internal/infra/store_db/internal/gen_jet/skyvault/public/model"
	. "skyvault/internal/infra/store_db/internal/gen_jet/skyvault/public/table"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/jinzhu/copier"
)

var _ media.Repo = (*MediaRepo)(nil)

type MediaRepo struct {
	store *store
}

func NewMediaRepo(db *store) *MediaRepo {
	return &MediaRepo{store: db}
}

func (r *MediaRepo) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.store.db.BeginTx(ctx, nil)
}

func (r *MediaRepo) WithTx(ctx context.Context, tx *sql.Tx) media.Repo {
	return &MediaRepo{store: r.store.WithTx(ctx, tx)}
}

func (r *MediaRepo) CreateFile(ctx context.Context, f *media.FileInfo) (*media.FileInfo, error) {
	dbFile := new(model.Files)
	err := copier.Copy(dbFile, f)
	if err != nil {
		return nil, err
	}

	stmt := Files.INSERT(
		Files.MutableColumns.Except(Files.CreatedAt, Files.UpdatedAt),
	).MODEL(dbFile).RETURNING(Files.AllColumns)

	return query[model.Files, media.FileInfo](ctx, stmt, r.store.dbTx)
}

func (r *MediaRepo) CreateFolder(ctx context.Context, folder *media.FolderInfo) (*media.FolderInfo, error) {
	dbFolder := new(model.Folders)
	err := copier.Copy(dbFolder, folder)
	if err != nil {
		return nil, err
	}

	stmt := Folders.INSERT(
		Folders.MutableColumns.Except(Folders.CreatedAt, Folders.UpdatedAt),
	).MODEL(dbFolder).RETURNING(Folders.AllColumns)

	return query[model.Folders, media.FolderInfo](ctx, stmt, r.store.dbTx)
}

func (r *MediaRepo) GetFile(ctx context.Context, fileID, ownerID int64) (*media.FileInfo, error) {
	stmt := SELECT(Files.AllColumns).
		FROM(Files).
		WHERE(Files.ID.EQ(Int64(fileID)).
			AND(Files.OwnerID.EQ(Int64(ownerID))).
			AND(Files.TrashedAt.IS_NULL()),
		)

	return query[model.Files, media.FileInfo](ctx, stmt, r.store.dbTx)
}

func (r *MediaRepo) GetFiles(ctx context.Context, ownerID int64, folderID *int64) ([]*media.FileInfo, error) {
	var folderCond BoolExpression
	if folderID == nil {
		folderCond = Files.FolderID.IS_NULL()
	} else {
		folderCond = Files.FolderID.EQ(Int64(*folderID))
	}

	stmt := SELECT(Files.AllColumns).
		FROM(Files).
		WHERE(
			Files.OwnerID.EQ(Int64(ownerID)).
				AND(folderCond).
				AND(Files.TrashedAt.IS_NULL()),
		)

	return querySlice[model.Files, media.FileInfo](ctx, stmt, r.store.dbTx)
}

func (r *MediaRepo) GetFolders(ctx context.Context, ownerID int64, folderID *int64) ([]*media.FolderInfo, error) {
	var folderCond BoolExpression
	if folderID == nil {
		folderCond = Folders.ParentFolderID.IS_NULL()
	} else {
		folderCond = Folders.ParentFolderID.EQ(Int64(*folderID))
	}

	stmt := SELECT(Folders.AllColumns).
		FROM(Folders).
		WHERE(
			Folders.OwnerID.EQ(Int64(ownerID)).
				AND(folderCond).
				AND(Folders.TrashedAt.IS_NULL()),
		)

	return querySlice[model.Folders, media.FolderInfo](ctx, stmt, r.store.dbTx)
}

func (r *MediaRepo) DeleteFile(ctx context.Context, fileID, ownerID int64) error {
	stmt := Files.DELETE().
		WHERE(Files.ID.EQ(Int64(fileID)).AND(Files.OwnerID.EQ(Int64(ownerID))))

	return exec(ctx, stmt, r.store.dbTx)
}
