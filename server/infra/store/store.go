package store

import (
	"skyvault/domain"
	"skyvault/infra/store/db_store"
)

var _ domain.IStore = &Store{}

type Store struct {
	*db_store.DBStore
}

func NewStore(dbStore *db_store.DBStore) *Store {
	return &Store{
		DBStore: dbStore,
	}
}
