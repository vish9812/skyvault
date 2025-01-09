package common

import (
	"context"
	"database/sql"
	"errors"
)

var ErrDBNoRows = errors.New("failed to find any row")
var ErrDBUniqueConstraint = errors.New("unique constraint")

type RepoTx[TRepo any] interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
	WithTx(ctx context.Context, tx *sql.Tx) TRepo
}
