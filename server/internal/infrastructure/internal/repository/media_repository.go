//lint:file-ignore ST1001 Using dot import to make SQL queries more readable
package repository

import (
	"context"
	"database/sql"
	"slices"
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
		FileInfo.AllColumns,
	).MODEL(dbModel).RETURNING(FileInfo.AllColumns)

	return runInsert[model.FileInfo, media.FileInfo](ctx, stmt, r.repository.dbTx)
}

func (r *MediaRepository) getFileInfo(ctx context.Context, fileID string, onlyTrashed bool) (*media.FileInfo, error) {
	whereCond := FileInfo.ID.EQ(UUID(UUIDStr(fileID)))

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

func (r *MediaRepository) GetFileInfo(ctx context.Context, fileID string) (*media.FileInfo, error) {
	return r.getFileInfo(ctx, fileID, false)
}

func (r *MediaRepository) GetFileInfoTrashed(ctx context.Context, fileID string) (*media.FileInfo, error) {
	return r.getFileInfo(ctx, fileID, true)
}

func (r *MediaRepository) getFileInfos(ctx context.Context, whereCond BoolExpression, pagingOpt *paging.Options, ownerID string, folderID *string, includeFolderID bool) (*paging.Page[*media.FileInfo], error) {
	if whereCond == nil {
		whereCond = Bool(true)
	}

	whereCond = whereCond.AND(FileInfo.OwnerID.EQ(UUID(UUIDStr(ownerID))))

	if includeFolderID {
		if folderID == nil {
			whereCond = whereCond.AND(FileInfo.FolderID.IS_NULL())
		} else {
			whereCond = whereCond.AND(FileInfo.FolderID.EQ(UUID(UUIDStr(*folderID))))
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
		return nil, apperror.NewAppError(err, "repository.GetFileInfos:runSelectSlice")
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

func (r *MediaRepository) GetFileInfos(ctx context.Context, pagingOpt *paging.Options, ownerID string, folderID *string) (*paging.Page[*media.FileInfo], error) {
	return r.getFileInfos(ctx, nil, pagingOpt, ownerID, folderID, true)
}

func (r *MediaRepository) GetFileInfosByCategory(ctx context.Context, pagingOpt *paging.Options, ownerID string, category media.Category) (*paging.Page[*media.FileInfo], error) {
	whereCond := FileInfo.Category.EQ(String(string(category)))
	return r.getFileInfos(ctx, whereCond, pagingOpt, ownerID, nil, false)
}

func (r *MediaRepository) UpdateFileInfo(ctx context.Context, info *media.FileInfo) error {
	stmt := FileInfo.UPDATE(FileInfo.MutableColumns).
		MODEL(info).
		WHERE(FileInfo.ID.EQ(UUID(UUIDStr(info.ID))))

	return runUpdateOrDelete(ctx, stmt, r.repository.dbTx)
}

func (r *MediaRepository) DeleteFileInfo(ctx context.Context, fileID string) error {
	stmt := FileInfo.DELETE().
		WHERE(FileInfo.ID.EQ(UUID(UUIDStr(fileID))))

	return runUpdateOrDelete(ctx, stmt, r.repository.dbTx)
}

func (r *MediaRepository) TrashFileInfos(ctx context.Context, ownerID string, fileIDs []string) error {
	inExp := make([]Expression, 0, len(fileIDs))
	for _, fileID := range fileIDs {
		inExp = append(inExp, UUID(UUIDStr(fileID)))
	}

	now := time.Now().UTC()
	fileInfo := model.FileInfo{
		TrashedAt: &now,
		UpdatedAt: now,
	}

	stmt := FileInfo.UPDATE(FileInfo.TrashedAt, FileInfo.UpdatedAt).
		MODEL(fileInfo).
		WHERE(
			FileInfo.ID.IN(inExp...).
				AND(FileInfo.TrashedAt.IS_NULL()).
				AND(FileInfo.OwnerID.EQ(UUID(UUIDStr(ownerID)))),
		)

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
		FolderInfo.AllColumns,
	).MODEL(dbModel).RETURNING(FolderInfo.AllColumns)

	return runInsert[model.FolderInfo, media.FolderInfo](ctx, stmt, r.repository.dbTx)
}

func (r *MediaRepository) getFolderInfo(ctx context.Context, folderID string, onlyTrashed bool) (*media.FolderInfo, error) {
	whereCond := FolderInfo.ID.EQ(UUID(UUIDStr(folderID)))

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

func (r *MediaRepository) GetFolderInfo(ctx context.Context, folderID string) (*media.FolderInfo, error) {
	return r.getFolderInfo(ctx, folderID, false)
}

func (r *MediaRepository) GetFolderInfoTrashed(ctx context.Context, folderID string) (*media.FolderInfo, error) {
	return r.getFolderInfo(ctx, folderID, true)
}

func (r *MediaRepository) GetFolderInfos(ctx context.Context, pagingOpt *paging.Options, ownerID string, parentFolderID *string) (*paging.Page[*media.FolderInfo], error) {
	whereCond := FolderInfo.OwnerID.EQ(UUID(UUIDStr(ownerID)))
	if parentFolderID == nil {
		whereCond = whereCond.AND(FolderInfo.ParentFolderID.IS_NULL())
	} else {
		whereCond = whereCond.AND(FolderInfo.ParentFolderID.EQ(UUID(UUIDStr(*parentFolderID))))
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
		return nil, apperror.NewAppError(err, "repository.GetFolderInfos:runSelectSlice")
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
		WHERE(FolderInfo.ID.EQ(UUID(UUIDStr(info.ID))))

	return runUpdateOrDelete(ctx, stmt, r.repository.dbTx)
}

func (r *MediaRepository) DeleteFolderInfo(ctx context.Context, folderID string) error {
	stmt := FolderInfo.DELETE().
		WHERE(FolderInfo.ID.EQ(UUID(UUIDStr(folderID))))

	return runUpdateOrDelete(ctx, stmt, r.repository.dbTx)
}

func (r *MediaRepository) TrashFolderInfos(ctx context.Context, ownerID string, folderIDs []string) error {
	nestedFoldersCTE := r.getNestedFoldersCTE(ownerID, folderIDs, false)
	trashFilesCTE := CTE("trash_files")

	now := time.Now().UTC()
	fileInfo := model.FileInfo{
		TrashedAt: &now,
		UpdatedAt: now,
	}
	folderInfo := model.FolderInfo{
		TrashedAt: &now,
		UpdatedAt: now,
	}

	// First trash all files in the folder and its sub-folders.
	// Then trash the folder and its sub-folders.
	stmt := WITH_RECURSIVE(
		nestedFoldersCTE,
		trashFilesCTE.AS(
			FileInfo.UPDATE(FileInfo.TrashedAt, FileInfo.UpdatedAt).
				MODEL(fileInfo).
				WHERE(
					FileInfo.FolderID.IN(
						SELECT(FolderInfo.ID.From(nestedFoldersCTE)).FROM(nestedFoldersCTE),
					).AND(FileInfo.TrashedAt.IS_NULL()),
				),
		),
	)(
		FolderInfo.UPDATE(FolderInfo.TrashedAt, FolderInfo.UpdatedAt).
			MODEL(folderInfo).
			WHERE(
				FolderInfo.ID.IN(
					SELECT(FolderInfo.ID.From(nestedFoldersCTE)).FROM(nestedFoldersCTE),
				),
			),
	)

	return runUpdateOrDelete(ctx, stmt, r.repository.dbTx)
}

func (r *MediaRepository) RestoreFolderInfo(ctx context.Context, ownerID string, folderID string) error {
	nestedFoldersCTE := r.getNestedFoldersCTE(ownerID, []string{folderID}, true)
	restoreFoldersCTE := CTE("restore_folders")

	now := time.Now().UTC()
	fileInfo := model.FileInfo{
		TrashedAt: nil,
		UpdatedAt: now,
	}
	folderInfo := model.FolderInfo{
		TrashedAt: nil,
		UpdatedAt: now,
	}

	// First restore the folder and its sub-folders
	// Then restore all files in those folders
	stmt := WITH_RECURSIVE(
		nestedFoldersCTE,
		restoreFoldersCTE.AS(
			FolderInfo.UPDATE(FolderInfo.TrashedAt, FolderInfo.UpdatedAt).
				MODEL(folderInfo).
				WHERE(
					FolderInfo.ID.IN(
						SELECT(FolderInfo.ID.From(nestedFoldersCTE)).FROM(nestedFoldersCTE),
					),
				),
		),
	)(
		FileInfo.UPDATE(FileInfo.TrashedAt, FileInfo.UpdatedAt).
			MODEL(fileInfo).
			WHERE(
				FileInfo.FolderID.IN(
					SELECT(FolderInfo.ID.From(nestedFoldersCTE)).FROM(nestedFoldersCTE),
				).AND(FileInfo.TrashedAt.IS_NOT_NULL()),
			),
	)

	return runUpdateOrDelete(ctx, stmt, r.repository.dbTx)
}

// getNestedFoldersCTE returns a CTE that returns all nested folders of the given folders, including the folders themselves.
func (r *MediaRepository) getNestedFoldersCTE(ownerID string, folderIDs []string, onlyTrashed bool) CommonTableExpression {
	nestedFolders := CTE("nested_folders")

	inExp := make([]Expression, 0, len(folderIDs))
	for _, folderID := range folderIDs {
		inExp = append(inExp, UUID(UUIDStr(folderID)))
	}

	trashedCond := FolderInfo.TrashedAt.IS_NULL()
	if onlyTrashed {
		trashedCond = FolderInfo.TrashedAt.IS_NOT_NULL()
	}

	return nestedFolders.AS(
		SELECT(
			FolderInfo.ID,
		).FROM(
			FolderInfo,
		).WHERE(
			FolderInfo.ID.IN(inExp...).
				AND(FolderInfo.OwnerID.EQ(UUID(UUIDStr(ownerID)))).
				AND(trashedCond),
		).UNION(
			SELECT(
				FolderInfo.ID,
			).FROM(
				FolderInfo.
					INNER_JOIN(nestedFolders, FolderInfo.ID.From(nestedFolders).EQ(FolderInfo.ParentFolderID)),
			).WHERE(
				FolderInfo.OwnerID.EQ(UUID(UUIDStr(ownerID))).
					AND(trashedCond),
			),
		),
	)
}

func (r *MediaRepository) GetDescendantFolderIDs(ctx context.Context, ownerID, folderID string) ([]string, error) {
	nestedFoldersCTE := r.getNestedFoldersCTE(ownerID, []string{folderID}, false)

	stmt := SELECT(FolderInfo.ID.From(nestedFoldersCTE)).
		FROM(nestedFoldersCTE)

	var folderIDs []string
	err := stmt.QueryContext(ctx, r.repository.dbTx, &folderIDs)
	if err != nil {
		return nil, apperror.NewAppError(err, "repository.GetDescendantFolderIDs:stmt.QueryContext")
	}

	return slices.DeleteFunc(folderIDs, func(id string) bool {
		return id == folderID
	}), nil
}
