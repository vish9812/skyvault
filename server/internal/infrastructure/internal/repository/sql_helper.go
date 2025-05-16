//lint:file-ignore ST1001 Using dot import to make SQL queries more readable
package repository

import (
	"context"
	"errors"
	"fmt"
	"skyvault/pkg/apperror"
	"skyvault/pkg/paging"
	"slices"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/jinzhu/copier"
)

type UUIDStr string

func (u UUIDStr) String() string {
	return string(u)
}

func ILIKE(lhs, rhs StringExpression) BoolExpression {
	return BoolExp(CustomExpression(lhs, Token("ILIKE"), rhs))
}

type cursorQuery struct {
	ID        ColumnString
	Name      ColumnString
	Updated   ColumnTimestamp
	where     BoolExpression
	orderBy   []OrderByClause
	pagingOpt *paging.Options
}

func (o *cursorQuery) buildClauses() error {
	o.pagingOpt.Validate()
	cursor, err := o.pagingOpt.GetCursor()
	if err != nil {
		if errors.Is(err, paging.ErrInvalidCursor) {
			return apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "repository.cursorQuery.buildClauses:GetCursor")
		}

		return apperror.NewAppError(err, "repository.cursorQuery.buildClauses:GetCursor")
	}

	o.initClauses()

	if o.pagingOpt.Direction == paging.DirectionForward {
		o.buildForwardClauses(cursor)
	} else {
		o.buildBackwardClauses(cursor)
	}

	return nil
}

func (o *cursorQuery) initClauses() {
	if o.where == nil {
		o.where = Bool(true)
	}

	if o.orderBy == nil {
		o.orderBy = []OrderByClause{}
	}
}

func (o *cursorQuery) buildForwardClauses(cursor *paging.Cursor) {
	switch o.pagingOpt.SortBy {
	case paging.SortByID:
		o.buildIDClauses(cursor, true)
	case paging.SortByName:
		o.buildNameClauses(cursor, true)
	case paging.SortByUpdated:
		o.buildUpdatedClauses(cursor, true)
	}
}

func (o *cursorQuery) buildBackwardClauses(cursor *paging.Cursor) {
	switch o.pagingOpt.SortBy {
	case paging.SortByID:
		o.buildIDClauses(cursor, false)
	case paging.SortByName:
		o.buildNameClauses(cursor, false)
	case paging.SortByUpdated:
		o.buildUpdatedClauses(cursor, false)
	}
}

func (o *cursorQuery) buildIDClauses(cursor *paging.Cursor, forward bool) {
	if o.pagingOpt.Sort == paging.SortAsc {
		if cursor == nil {
			o.orderBy = append(o.orderBy, o.ID.ASC())
			return
		}

		if forward {
			o.orderBy = append(o.orderBy, o.ID.ASC())
			o.where = o.where.AND(o.ID.GT(UUID(UUIDStr(cursor.ID))))
			return
		}

		o.orderBy = append(o.orderBy, o.ID.DESC())
		o.where = o.where.AND(o.ID.LT(UUID(UUIDStr(cursor.ID))))
		return
	}

	if cursor == nil {
		o.orderBy = append(o.orderBy, o.ID.DESC())
		return
	}

	if forward {
		o.orderBy = append(o.orderBy, o.ID.DESC())
		o.where = o.where.AND(o.ID.LT(UUID(UUIDStr(cursor.ID))))
		return
	}

	o.orderBy = append(o.orderBy, o.ID.ASC())
	o.where = o.where.AND(o.ID.GT(UUID(UUIDStr(cursor.ID))))
}

func (o *cursorQuery) buildNameClauses(cursor *paging.Cursor, forward bool) {
	if o.pagingOpt.Sort == paging.SortAsc {
		if cursor == nil {
			o.orderBy = append(o.orderBy, o.Name.ASC(), o.ID.ASC())
			return
		}

		if forward {
			o.orderBy = append(o.orderBy, o.Name.ASC(), o.ID.ASC())
			o.where = o.where.AND(
				o.Name.GT(String(cursor.Name)).
					OR(o.Name.EQ(String(cursor.Name)).AND(
						o.ID.GT(UUID(UUIDStr(cursor.ID))),
					)),
			)
			return
		}

		o.orderBy = append(o.orderBy, o.Name.DESC(), o.ID.DESC())
		o.where = o.where.AND(
			o.Name.LT(String(cursor.Name)).
				OR(o.Name.EQ(String(cursor.Name)).AND(
					o.ID.LT(UUID(UUIDStr(cursor.ID))),
				)),
		)
		return
	}

	if cursor == nil {
		o.orderBy = append(o.orderBy, o.Name.DESC(), o.ID.DESC())
		return
	}

	if forward {
		o.orderBy = append(o.orderBy, o.Name.DESC(), o.ID.DESC())
		o.where = o.where.AND(
			o.Name.LT(String(cursor.Name)).
				OR(o.Name.EQ(String(cursor.Name)).AND(
					o.ID.LT(UUID(UUIDStr(cursor.ID))),
				)),
		)
		return
	}

	o.orderBy = append(o.orderBy, o.Name.ASC(), o.ID.ASC())
	o.where = o.where.AND(
		o.Name.GT(String(cursor.Name)).
			OR(o.Name.EQ(String(cursor.Name)).AND(
				o.ID.GT(UUID(UUIDStr(cursor.ID))),
			)),
	)
}

func (o *cursorQuery) buildUpdatedClauses(cursor *paging.Cursor, forward bool) {
	if o.pagingOpt.Sort == paging.SortAsc {
		if cursor == nil {
			o.orderBy = append(o.orderBy, o.Updated.ASC(), o.ID.ASC())
			return
		}

		if forward {
			o.orderBy = append(o.orderBy, o.Updated.ASC(), o.ID.ASC())
			o.where = o.where.AND(
				o.Updated.GT(TimestampT(cursor.Updated)).
					OR(o.Updated.EQ(TimestampT(cursor.Updated)).AND(
						o.ID.GT(UUID(UUIDStr(cursor.ID))),
					)),
			)
			return
		}

		o.orderBy = append(o.orderBy, o.Updated.DESC(), o.ID.DESC())
		o.where = o.where.AND(
			o.Updated.LT(TimestampT(cursor.Updated)).
				OR(o.Updated.EQ(TimestampT(cursor.Updated)).AND(
					o.ID.LT(UUID(UUIDStr(cursor.ID))),
				)),
		)
		return
	}

	if cursor == nil {
		o.orderBy = append(o.orderBy, o.Updated.DESC(), o.ID.DESC())
		return
	}

	if forward {
		o.orderBy = append(o.orderBy, o.Updated.DESC(), o.ID.DESC())
		o.where = o.where.AND(
			o.Updated.LT(TimestampT(cursor.Updated)).
				OR(o.Updated.EQ(TimestampT(cursor.Updated)).AND(
					o.ID.LT(UUID(UUIDStr(cursor.ID))),
				)),
		)
		return
	}

	o.orderBy = append(o.orderBy, o.Updated.ASC(), o.ID.ASC())
	o.where = o.where.AND(
		o.Updated.GT(TimestampT(cursor.Updated)).
			OR(o.Updated.EQ(TimestampT(cursor.Updated)).AND(
				o.ID.GT(UUID(UUIDStr(cursor.ID))),
			)),
	)
}

// runSelect is to be used with Select statements that return a single row.
//
// App Errors:
// - ErrCommonNoData
func runSelect[TDBModel any, TRes any](ctx context.Context, stmt Statement, dbTx qrm.DB) (*TRes, error) {
	var dbModel TDBModel
	err := stmt.QueryContext(ctx, dbTx, &dbModel)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonNoData, err), "repository.runSelect:QueryContext:ErrNoRows")
		}

		return nil, apperror.NewAppError(err, "repository.runSelect:QueryContext")
	}

	var resModel TRes
	err = copier.Copy(&resModel, &dbModel)
	if err != nil {
		return nil, apperror.NewAppError(fmt.Errorf("failed to copy the db model to the res model: %w", err), "repository.runSelect:Copy")
	}

	return &resModel, nil
}

// runSelectSlice is to be used with Select statements that return multiple rows with pagination.
// it returns empty items if no data is found.
func runSelectSlice[TDBModel any, TRes any](ctx context.Context, cursorQuery *cursorQuery, stmt SelectStatement, dbTx qrm.DB) (*paging.Page[*TRes], error) {
	err := cursorQuery.buildClauses()
	if err != nil {
		return nil, apperror.NewAppError(err, "repository.runSelectSlice:cursorOptions.buildClauses")
	}

	stmt = stmt.WHERE(cursorQuery.where).
		ORDER_BY(cursorQuery.orderBy...).
		LIMIT(int64(cursorQuery.pagingOpt.Limit) + 1)

	// Uncomment this when you need to debug the sql statement
	// applog.GetLoggerFromContext(ctx).Debug().Msg(stmt.DebugSql())

	page := &paging.Page[*TRes]{
		Items: []*TRes{},
	}
	res, err := runSelect[[]*TDBModel, []*TRes](ctx, stmt, dbTx)
	if err != nil {
		if errors.Is(err, apperror.ErrCommonNoData) {
			return page, nil
		}

		return nil, apperror.NewAppError(err, "repository.runSelectSlice:runSelect")
	}
	items := *res

	if len(items) > cursorQuery.pagingOpt.Limit {
		page.HasMore = true
		items = items[:cursorQuery.pagingOpt.Limit]
	}

	if cursorQuery.pagingOpt.Direction != paging.DirectionForward {
		slices.Reverse(items)
	}

	page.Items = items

	return page, nil
}

// runSelectSliceAll is to be used with Select statements that return multiple rows without pagination.
func runSelectSliceAll[TDBModel any, TRes any](ctx context.Context, stmt SelectStatement, dbTx qrm.DB) ([]*TRes, error) {
	var dbModels []*TDBModel
	err := stmt.QueryContext(ctx, dbTx, &dbModels)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return []*TRes{}, nil
		}
		return nil, apperror.NewAppError(err, "repository.runSelectSliceAll:QueryContext")
	}

	var resModels []*TRes
	err = copier.Copy(&resModels, &dbModels)
	if err != nil {
		return nil, apperror.NewAppError(fmt.Errorf("failed to copy the db models to the res models: %w", err), "repository.runSelectSliceAll:Copy")
	}

	return resModels, nil
}

// runInsert is to be used with Insert statements
//
// App Errors:
// - ErrCommonDuplicateData
func runInsert[TDBModel any, TRes any](ctx context.Context, stmt Statement, dbTx qrm.DB) (*TRes, error) {
	var dbModel TDBModel
	err := stmt.QueryContext(ctx, dbTx, &dbModel)
	if err != nil {
		if apperror.Contains(err, "unique constraint") {
			return nil, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonDuplicateData, err), "repository.runInsert:QueryContext")
		}

		return nil, apperror.NewAppError(err, "repository.runInsert:QueryContext")
	}

	var resModel TRes
	err = copier.Copy(&resModel, &dbModel)
	if err != nil {
		return nil, apperror.NewAppError(fmt.Errorf("failed to copy the db model to the res model: %w", err), "repository.runInsert:Copy")
	}

	return &resModel, nil
}

// runInsertNoReturn is to be used with Insert statements that do not return a value
//
// App Errors:
// - ErrCommonDuplicateData
// - ErrCommonNoData
func runInsertNoReturn(ctx context.Context, stmt Statement, dbTx qrm.DB) error {
	res, err := stmt.ExecContext(ctx, dbTx)
	if err != nil {
		if apperror.Contains(err, "unique constraint") {
			return apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonDuplicateData, err), "repository.runInsertNoReturn:ExecContext")
		}

		return apperror.NewAppError(err, "repository.runInsertNoReturn:ExecContext")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return apperror.NewAppError(err, "repository.runInsertNoReturn:RowsAffected")
	}

	if rowsAffected == 0 {
		return apperror.NewAppError(apperror.ErrCommonNoData, "repository.runInsertNoReturn:NoRowsAffected")
	}

	return nil
}

// runUpdateOrDelete is to be used with Update or Delete statements
//
// App Errors:
// - ErrCommonNoData
func runUpdateOrDelete(ctx context.Context, stmt Statement, dbTx qrm.DB) error {
	res, err := stmt.ExecContext(ctx, dbTx)
	if err != nil {
		return apperror.NewAppError(err, "repository.runUpdateOrDelete:ExecContext")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return apperror.NewAppError(err, "repository.runUpdateOrDelete:RowsAffected")
	}

	if rowsAffected == 0 {
		return apperror.NewAppError(apperror.ErrCommonNoData, "repository.runUpdateOrDelete:RowsAffected")
	}

	return nil
}
