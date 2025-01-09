package auth

import (
	"context"
	"skyvault/pkg/common"
)

type Repo interface {
	common.RepoTx[Repo]
	Create(ctx context.Context, au *Auth) (*Auth, error)
	Get(ctx context.Context, id int64) (*Auth, error)
	GetByProfileID(ctx context.Context, id int64) (*Auth, error)
}
