package profile

import (
	"context"
	"skyvault/pkg/common"
)

type Repo interface {
	common.RepoTx[Repo]
	Create(ctx context.Context, pro *Profile) (*Profile, error)
	Get(ctx context.Context, id int64) (*Profile, error)
	GetByEmail(ctx context.Context, email string) (*Profile, error)
}
