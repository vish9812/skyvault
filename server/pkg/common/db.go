package common

import (
	"context"
	"database/sql"
)

type RepoTx[TRepo any] interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
	WithTx(ctx context.Context, tx *sql.Tx) TRepo
}
