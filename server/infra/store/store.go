package store

import (
	"skyvault/domain/auth"
	"skyvault/infra/store/db_store"
)

var _ auth.Repo = &BaseAuthRepo{}

type Store struct {
	db *db_store.DBStore
}

func NewStore(dbStore *db_store.DBStore) *Store {
	return &Store{
		db: dbStore,
	}
}

type BaseAuthRepo struct {
	db_store.AuthRepo
}
