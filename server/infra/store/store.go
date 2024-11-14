package store

import (
	"skyvault/domain/auth"
	"skyvault/infra/store/db_store"
)

var _ auth.Repo = &db_store.AuthRepo{}

type Store struct {
	*db_store.DBStore
}

func NewStore(dbStore *db_store.DBStore) *Store {
	return &Store{
		DBStore: dbStore,
	}
}
