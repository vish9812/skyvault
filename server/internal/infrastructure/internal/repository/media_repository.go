package repository

import (
	"context"
	"database/sql"
	"time"

	"skyvault/internal/domain/media"
	"skyvault/internal/infrastructure/internal/repository/internal/gen_jet/skyvault/public/model"
	. "skyvault/internal/infrastructure/internal/repository/internal/gen_jet/skyvault/public/table"
	"skyvault/pkg/apperror"
	"skyvault/pkg/paging"

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
	return &MediaRepository{repository: r.repository.withTx(ctx, tx)}
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

func (r *MediaRepository) getFileInfo(ctx context.Context, fileID int64, onlyTrashed bool) (*media.FileInfo, error) {
	whereCond := FileInfo.ID.EQ(Int64(fileID))

	if onlyTrashed {
		whereCond = whereCond.AND(FileInfo.TrashedAt.IS_NOT_NULL())
	} else {
		whereCond = whereCond.AND(FileInfo.TrashedAt.IS_NULL())
	}

	stmt := SELECT(FileInfo.AllColumns).
		FROM(FileInfo).
		WHERE(whereCond)

	return runSelect[model.FileInfo, media.FileInfo](ctx, stmt, r.repository.dbTx)
}

func (r *MediaRepository) GetFileInfo(ctx context.Context, fileID int64) (*media.FileInfo, error) {
	return r.getFileInfo(ctx, fileID, false)
}

func (r *MediaRepository) GetFileInfoTrashed(ctx context.Context, fileID int64) (*media.FileInfo, error) {
	return r.getFileInfo(ctx, fileID, true)
}

func (r *MediaRepository) getFilesInfo(ctx context.Context, whereCond BoolExpression, pagingOpt *paging.Options, ownerID int64, folderID *int64, includeFolderID bool) (*paging.Page[*media.FileInfo], error) {
	if whereCond == nil {
		whereCond = Bool(true)
	}

	whereCond = whereCond.AND(FileInfo.OwnerID.EQ(Int64(ownerID)))

	if includeFolderID {
		if folderID == nil {
			whereCond = whereCond.AND(FileInfo.FolderID.IS_NULL())
		} else {
			whereCond = whereCond.AND(FileInfo.FolderID.EQ(Int64(*folderID)))
		}
	}

	whereCond = whereCond.AND(FileInfo.TrashedAt.IS_NULL())

	orderBy := []OrderByClause{FileInfo.OwnerID.ASC(), FileInfo.FolderID.ASC()}

	stmt := SELECT(FileInfo.AllColumns).
		FROM(FileInfo)

	cursorQuery := &cursorQuery{
		ID:        FileInfo.ID,
		Name:      FileInfo.Name,
		Updated:   FileInfo.UpdatedAt,
		where:     whereCond,
		orderBy:   orderBy,
		pagingOpt: pagingOpt,
	}

	page, err := runSelectSlice[model.FileInfo, media.FileInfo](ctx, cursorQuery, stmt, r.repository.dbTx)
	if err != nil {
		return nil, apperror.NewAppError(err, "repository.GetFilesInfo:runSelectSlice")
	}

	if len(page.Items) > 0 {
		lastItem := page.Items[len(page.Items)-1]
		nextCursor := &paging.Cursor{
			ID:      lastItem.ID,
			Name:    lastItem.Name,
			Updated: lastItem.UpdatedAt,
		}
		page.NextCursor = pagingOpt.CreateCursor(nextCursor)

		firstItem := page.Items[0]
		prevCursor := &paging.Cursor{
			ID:      firstItem.ID,
			Name:    firstItem.Name,
			Updated: firstItem.UpdatedAt,
		}
		page.PrevCursor = pagingOpt.CreateCursor(prevCursor)
	}

	return page, nil
}

func (r *MediaRepository) GetFilesInfo(ctx context.Context, pagingOpt *paging.Options, ownerID int64, folderID *int64) (*paging.Page[*media.FileInfo], error) {
	return r.getFilesInfo(ctx, nil, pagingOpt, ownerID, folderID, true)
}

func (r *MediaRepository) GetFilesInfoByCategory(ctx context.Context, pagingOpt *paging.Options, ownerID int64, category string) (*paging.Page[*media.FileInfo], error) {
	whereCond := FileInfo.Category.EQ(String(category))
	return r.getFilesInfo(ctx, whereCond, pagingOpt, ownerID, nil, false)
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

func (r *MediaRepository) getFolderInfo(ctx context.Context, folderID int64, onlyTrashed bool) (*media.FolderInfo, error) {
	whereCond := FolderInfo.ID.EQ(Int64(folderID))

	if onlyTrashed {
		whereCond = whereCond.AND(FolderInfo.TrashedAt.IS_NOT_NULL())
	} else {
		whereCond = whereCond.AND(FolderInfo.TrashedAt.IS_NULL())
	}

	stmt := SELECT(FolderInfo.AllColumns).
		FROM(FolderInfo).
		WHERE(whereCond)

	return runSelect[model.FolderInfo, media.FolderInfo](ctx, stmt, r.repository.dbTx)
}

func (r *MediaRepository) GetFolderInfo(ctx context.Context, folderID int64) (*media.FolderInfo, error) {
	return r.getFolderInfo(ctx, folderID, false)
}

func (r *MediaRepository) GetFolderInfoTrashed(ctx context.Context, folderID int64) (*media.FolderInfo, error) {
	return r.getFolderInfo(ctx, folderID, true)
}

func (r *MediaRepository) GetFoldersInfo(ctx context.Context, pagingOpt *paging.Options, ownerID int64, parentFolderID *int64) (*paging.Page[*media.FolderInfo], error) {
	whereCond := FolderInfo.OwnerID.EQ(Int64(ownerID))
	if parentFolderID == nil {
		whereCond = whereCond.AND(FolderInfo.ParentFolderID.IS_NULL())
	} else {
		whereCond = whereCond.AND(FolderInfo.ParentFolderID.EQ(Int64(*parentFolderID)))
	}

	whereCond = whereCond.AND(FolderInfo.TrashedAt.IS_NULL())
	orderBy := []OrderByClause{FolderInfo.OwnerID.ASC(), FolderInfo.ParentFolderID.ASC()}

	stmt := SELECT(FolderInfo.AllColumns).
		FROM(FolderInfo)

	cursorQuery := &cursorQuery{
		ID:        FolderInfo.ID,
		Name:      FolderInfo.Name,
		Updated:   FolderInfo.UpdatedAt,
		where:     whereCond,
		orderBy:   orderBy,
		pagingOpt: pagingOpt,
	}

	page, err := runSelectSlice[model.FolderInfo, media.FolderInfo](ctx, cursorQuery, stmt, r.repository.dbTx)
	if err != nil {
		return nil, apperror.NewAppError(err, "repository.GetFoldersInfo:runSelectSlice")
	}

	if len(page.Items) > 0 {
		lastItem := page.Items[len(page.Items)-1]
		nextCursor := &paging.Cursor{
			ID:      lastItem.ID,
			Name:    lastItem.Name,
			Updated: lastItem.UpdatedAt,
		}
		page.NextCursor = pagingOpt.CreateCursor(nextCursor)

		firstItem := page.Items[0]
		prevCursor := &paging.Cursor{
			ID:      firstItem.ID,
			Name:    firstItem.Name,
			Updated: firstItem.UpdatedAt,
		}
		page.PrevCursor = pagingOpt.CreateCursor(prevCursor)
	}

	return page, nil
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

func (r *MediaRepository) TrashFolderInfo(ctx context.Context, folderID int64) error {
	nestedFolders, withStmt := r.getNestedFoldersCTE(folderID)
	
	nowUTC := TimestampT(time.Now().UTC())
	
	stmt := withStmt(
		FolderInfo.UPDATE().
			SET(FolderInfo.TrashedAt.SET(nowUTC)).
			JOIN(FileInfo, FileInfo.FolderID.EQ(FolderInfo.ID)).
			WHERE(
				FolderInfo.ID.IN(
					SELECT(FolderInfo.ID.From(nestedFolders)).FROM(nestedFolders),
				).AND(FolderInfo.TrashedAt.IS_NULL()),
			).
			SET(FileInfo.TrashedAt.SET(nowUTC)).
			WHERE(FileInfo.TrashedAt.IS_NULL()),
	)

	return runUpdateOrDelete(ctx, stmt, r.repository.dbTx)
}

func (r *MediaRepository) getNestedFoldersCTE(folderID int64) (CommonTableExpression, func(Statement) Statement) {
	nestedFolders := CTE("nested_folders")

	return nestedFolders, WITH_RECURSIVE(
		nestedFolders.AS(
			SELECT(
				FolderInfo.ID,
			).FROM(
				FolderInfo,
			).WHERE(
				FolderInfo.ID.EQ(Int64(folderID)).AND(FolderInfo.TrashedAt.IS_NULL()),
			).UNION(
				SELECT(
					FolderInfo.ID,
				).FROM(
					FolderInfo.
						INNER_JOIN(nestedFolders, FolderInfo.ID.From(nestedFolders).EQ(FolderInfo.ParentFolderID)),
				).WHERE(
					FolderInfo.TrashedAt.IS_NULL(),
				),
			),
		),
	)
}
