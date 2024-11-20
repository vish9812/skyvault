package store

import (
	"skyvault/infra/store/db_store"
)

type Store struct {
	*db_store.DBStore
}

func NewStore(dbStore *db_store.DBStore) *Store {
	return &Store{
		DBStore: dbStore,
	}
}
