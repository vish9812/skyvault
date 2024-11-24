package domain

import "skyvault/domain/auth"

type IStore interface {
	NewAuthRepo() auth.Repo
}

type Repo interface {
	auth.Repo
}
