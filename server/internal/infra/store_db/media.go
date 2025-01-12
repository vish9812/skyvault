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

func (r *MediaRepo) CreateFile(ctx context.Context, f *media.File) (*media.File, error) {
	dbFile := new(model.Files)
	err := copier.Copy(dbFile, f)
	if err != nil {
		return nil, err
	}

	stmt := Files.INSERT(
		Files.MutableColumns.Except(Files.CreatedAt, Files.UpdatedAt),
	).MODEL(dbFile).RETURNING(Files.AllColumns)

	return get[model.Files, media.File](ctx, stmt, r.store.exec)
}

func (r *MediaRepo) CreateFolder(ctx context.Context, folder *media.Folder) (*media.Folder, error) {
	dbFolder := new(model.Folders)
	err := copier.Copy(dbFolder, folder)
	if err != nil {
		return nil, err
	}

	stmt := Folders.INSERT(
		Folders.MutableColumns.Except(Folders.CreatedAt, Folders.UpdatedAt),
	).MODEL(dbFolder).RETURNING(Folders.AllColumns)

	return get[model.Folders, media.Folder](ctx, stmt, r.store.exec)
}

func (r *MediaRepo) GetFiles(ctx context.Context, ownerID int64, folderID *int64) ([]*media.File, error) {
	var folderCond BoolExpression
	if folderID != nil {
		folderCond = Files.FolderID.EQ(Int64(*folderID))
	} else {
		folderCond = Files.FolderID.IS_NULL()
	}

	stmt := SELECT(Files.AllColumns).
		FROM(Files).
		WHERE(
			Files.OwnerID.EQ(Int64(ownerID)).
				AND(folderCond).
				AND(Files.TrashedAt.IS_NULL()),
		)

	return getSlice[model.Files, media.File](ctx, stmt, r.store.exec)
}

func (r *MediaRepo) GetFolders(ctx context.Context, ownerID int64, folderID *int64) ([]*media.Folder, error) {
	var folderCond BoolExpression
	if folderID != nil {
		folderCond = Folders.ParentFolderID.EQ(Int64(*folderID))
	} else {
		folderCond = Folders.ParentFolderID.IS_NULL()
	}

	stmt := SELECT(Folders.AllColumns).
		FROM(Folders).
		WHERE(
			Folders.OwnerID.EQ(Int64(ownerID)).
				AND(folderCond).
				AND(Folders.TrashedAt.IS_NULL()),
		)

	return getSlice[model.Folders, media.Folder](ctx, stmt, r.store.exec)
}
