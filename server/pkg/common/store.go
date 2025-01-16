package common

import (
	"context"
	"database/sql"
	"errors"
)

var ErrNoData = errors.New("no data found")
var ErrDuplicateData = errors.New("found duplicate data")

// RepoTx is an interface that represents a repository that supports transactions
//
// Usage:
//
// - First use BeginTx to start a new transaction
//
// - Then use WithTx to get a new Repo with the transaction
//
// - Then use the new Repo to call the repo methods (Make sure to use the new Repo for all operations, not the original one)
//
// - Finally, use Commit or Rollback to commit or rollback the transaction
type RepoTx[TRepo any] interface {
	// BeginTx starts a new transaction.
	// Pass the returned transaction to WithTx to get a new Repo with the transaction
	BeginTx(ctx context.Context) (*sql.Tx, error)
	// WithTx returns a new Repo with the given transaction.
	// Make sure to use the new Repo for all operations, not the original one
	WithTx(ctx context.Context, tx *sql.Tx) TRepo
}
