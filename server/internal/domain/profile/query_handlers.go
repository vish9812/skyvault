package profile

import "context"

var _ Queries = (*QueryHandlers)(nil)

type QueryHandlers struct {
	repository Repository
}

func NewQueryHandlers(repository Repository) Queries {
	return &QueryHandlers{repository: repository}
}

func (h *QueryHandlers) Get(ctx context.Context, query *GetQuery) (*Profile, error) {
	return h.repository.Get(ctx, query.ID)
}
