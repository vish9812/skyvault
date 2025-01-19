package store_db

import (
	"skyvault/internal/domain/auth"
	"skyvault/internal/domain/media"
	"skyvault/internal/domain/profile"
)

type Repo struct {
	Store   *Store
	Auth    auth.Repo
	Profile profile.Repo
	Media   media.Repo
}

func NewRepo(store *Store) *Repo {
	return &Repo{
		Store:   store,
		Auth:    NewAuthRepo(store),
		Profile: NewProfileRepo(store),
		Media:   NewMediaRepo(store),
	}
}
