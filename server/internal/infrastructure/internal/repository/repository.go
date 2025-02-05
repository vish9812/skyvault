package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"skyvault/internal/domain/auth"
	"skyvault/internal/domain/media"
	"skyvault/internal/domain/profile"
	"skyvault/pkg/appconfig"
	"skyvault/pkg/apperror"
	"skyvault/pkg/applog"
	"skyvault/pkg/common"
	"time"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jinzhu/copier"

	jetpg "github.com/go-jet/jet/v2/postgres"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Repository struct {
	app *appconfig.App

	// db is to be used with the standard library queries.
	// Do NOT use it with the go-jet library queries.
	// Use dbTx instead which can use both sql.db and sql.Tx interchangeably.
	db *sql.DB

	// dbTx can use both sql.DB and sql.Tx interchangeably.
	// It is to be used with the go-jet library queries.
	dbTx qrm.DB

	// Repositories
	Auth    auth.Repository
	Profile profile.Repository
	Media   media.Repository
}

func NewRepository(app *appconfig.App) *Repository {
	logger := app.Logger.With().Str("where", "NewRepository").Str("dsn", app.Config.DB.DSN).Logger()

	db := connectDatabase(logger, app.Config.DB.DSN)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := db.PingContext(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to ping the db")
	}

	logger.Info().Msg("connected to the db")

	repo := &Repository{app: app, db: db, dbTx: db}

	repo.migrateUp()

	repo.initRepositories()

	return repo
}

func (r *Repository) initRepositories() {
	r.Auth = NewAuthRepository(r)
	r.Profile = NewProfileRepository(r)
	r.Media = NewMediaRepository(r)
}

func connectDatabase(logger applog.Logger, dsn string) *sql.DB {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to open the db")
	}
	return db
}

func (r *Repository) Cleanup() error {
	return r.close()
}

func (r *Repository) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return r.db.PingContext(ctx)
}

func (r *Repository) close() error {
	if err := r.db.Close(); err != nil {
		return fmt.Errorf("failed to close the db: %w", err)
	}
	return nil
}

// TODO: No need to export this method.
func (r *Repository) WithTx(ctx context.Context, tx *sql.Tx) *Repository {
	return &Repository{app: r.app, db: r.db, dbTx: tx}
}

func (r *Repository) migrateUp() {
	// Assuming the code is run from the main.go inside the cmd folder.
	migrationPath, err := filepath.Abs("../internal/infrastructure/internal/repository/internal/migrations")
	if err != nil {
		r.app.Logger.Fatal().Err(err).Msg("failed to get the absolute migration path")
	}
	migrationDirURL := fmt.Sprintf("file://%s", migrationPath)
	logger := r.app.Logger.With().Str("where", "migrateUp").Str("migration_path", migrationDirURL).Logger()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, err := r.db.Conn(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to get the connection")
	}
	p, err := postgres.WithConnection(ctx, conn, &postgres.Config{})
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create the new postgres instance")
	}
	defer func() {
		if err := p.Close(); err != nil {
			logger.Fatal().Err(err).Msg("failed to close the postgres instance")
		}
	}()
	migrateObj, err := migrate.NewWithDatabaseInstance(migrationDirURL, "postgres", p)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create the new migrate instance")
	}

	err = migrateObj.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info().Msg("no new migration needed")
		} else {
			logger.Fatal().Err(err).Msg("failed to migrate up")
		}

		return
	}

	logger.Info().Msg("migrated db up")
}

// -----------------------------
// DB Queries runner functions
// -----------------------------

type cursorQueryOptions struct {
	ID        jetpg.ColumnInteger
	Name      jetpg.ColumnString
	Updated   jetpg.ColumnTimestamp
	where     jetpg.BoolExpression
	orderBy   []jetpg.OrderByClause
	pagingOpt *common.PagingOptions
}

func (o *cursorQueryOptions) buildClauses() error {
	o.pagingOpt.Validate()
	cursor, err := o.pagingOpt.GetCursor()
	if err != nil {
		return apperror.NewAppError(err, "repository.cursorQueryOptions.buildClauses:GetCursor")
	}

	if o.where == nil {
		o.where = jetpg.Bool(true)
	}

	if o.orderBy == nil {
		o.orderBy = []jetpg.OrderByClause{}
	}

	switch o.pagingOpt.SortBy {
	case common.SortByID:
		if o.pagingOpt.Sort == common.SortAsc {
			o.orderBy = append(o.orderBy, o.ID.ASC())

			if cursor != nil {
				o.where = o.where.AND(o.ID.GT(jetpg.Int64(cursor.ID)))
			}

		} else {
			o.orderBy = append(o.orderBy, o.ID.DESC())

			if cursor != nil {
				o.where = o.where.AND(o.ID.LT(jetpg.Int64(cursor.ID)))
			}
		}

	case common.SortByName:
		if o.pagingOpt.Sort == common.SortAsc {
			o.orderBy = append(o.orderBy, o.Name.ASC(), o.ID.ASC())

			if cursor != nil {
				o.where = o.where.AND(
					o.Name.GT(jetpg.String(cursor.Name)).
						OR(o.Name.EQ(jetpg.String(cursor.Name)).AND(
							o.ID.GT(jetpg.Int64(cursor.ID)),
						)),
				)
			}

		} else {
			o.orderBy = append(o.orderBy, o.Name.DESC(), o.ID.DESC())

			if cursor != nil {
				o.where = o.where.AND(
					o.Name.LT(jetpg.String(cursor.Name)).
						OR(o.Name.EQ(jetpg.String(cursor.Name)).AND(
							o.ID.LT(jetpg.Int64(cursor.ID)),
						)),
				)
			}
		}

	case common.SortByUpdated:
		if o.pagingOpt.Sort == common.SortAsc {
			o.orderBy = append(o.orderBy, o.Updated.ASC(), o.ID.ASC())

			if cursor != nil {
				o.where = o.where.AND(
					o.Updated.GT(jetpg.TimestampT(cursor.Updated)).
						OR(o.Updated.EQ(jetpg.TimestampT(cursor.Updated)).AND(
							o.ID.GT(jetpg.Int64(cursor.ID)),
						)),
				)
			}

		} else {
			o.orderBy = append(o.orderBy, o.Updated.DESC(), o.ID.DESC())

			if cursor != nil {
				o.where = o.where.AND(
					o.Updated.LT(jetpg.TimestampT(cursor.Updated)).
						OR(o.Updated.EQ(jetpg.TimestampT(cursor.Updated)).AND(
							o.ID.LT(jetpg.Int64(cursor.ID)),
						)),
				)
			}
		}
	}

	return nil
}

// runSelect is to be used with Select statements that return a single row.
//
// App Errors:
// - apperror.ErrNoData
func runSelect[TDBModel any, TRes any](ctx context.Context, stmt jetpg.Statement, dbTx qrm.DB) (*TRes, error) {
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

// runSelectSlice is to be used with Select statements that return multiple rows.
//
// App Errors:
// - apperror.ErrNoData
func runSelectSlice[TDBModel any, TRes any](ctx context.Context, cursorOptions *cursorQueryOptions, stmt jetpg.SelectStatement, dbTx qrm.DB) (*common.PagedItems[*TRes], error) {
	err := cursorOptions.buildClauses()
	if err != nil {
		return nil, apperror.NewAppError(err, "repository.runSelectSlice:cursorOptions.buildClauses")
	}

	stmt = stmt.WHERE(cursorOptions.where).
		ORDER_BY(cursorOptions.orderBy...).
		LIMIT(int64(cursorOptions.pagingOpt.Limit) + 1)

	// TODO: Comment this once this function is tested.
	applog.GetLoggerFromContext(ctx).Debug().Msg(stmt.DebugSql())

	res, err := runSelect[[]*TDBModel, []*TRes](ctx, stmt, dbTx)
	if err != nil {
		return nil, apperror.NewAppError(err, "repository.runSelectSlice:runSelect")
	}
	items := *res

	pagedItems := &common.PagedItems[*TRes]{}

	if len(items) > cursorOptions.pagingOpt.Limit {
		pagedItems.HasMore = true
		if cursorOptions.pagingOpt.Direction == common.DirectionNext {
			items = items[:cursorOptions.pagingOpt.Limit]
		} else {
			items = items[1:]
		}
	}

	pagedItems.Items = items

	return pagedItems, nil
}

// runInsert is to be used with Insert statements
//
// App Errors:
// - apperror.ErrDuplicateData
func runInsert[TDBModel any, TRes any](ctx context.Context, stmt jetpg.Statement, dbTx qrm.DB) (*TRes, error) {
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

// runUpdateOrDelete is to be used with Update or Delete statements
//
// App Errors:
// - apperror.ErrNoData
func runUpdateOrDelete(ctx context.Context, stmt jetpg.Statement, dbTx qrm.DB) error {
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
