package profile

import "context"

var _ Queries = (*QueryHandlers)(nil)

type QueryHandlers struct {
	repository Repository
}

func NewQueryHandlers(repository Repository) *QueryHandlers {
	return &QueryHandlers{repository: repository}
}

func (h *QueryHandlers) Get(ctx context.Context, query *GetQuery) (*Profile, error) {
	return h.repository.Get(ctx, query.ID)
}

func (h *QueryHandlers) GetByEmail(ctx context.Context, query *GetByEmailQuery) (*Profile, error) {
	return h.repository.GetByEmail(ctx, query.Email)
}
