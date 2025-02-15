package media

import (
	"context"
	"skyvault/pkg/apperror"
	"skyvault/pkg/paging"
	"strings"
)

var _ Queries = (*QueriesSanitizer)(nil)

type QueriesSanitizer struct {
	Queries
}

func NewQueriesSanitizer(queries Queries) *QueriesSanitizer {
	return &QueriesSanitizer{Queries: queries}
}

func validateCategory(category Category) (Category, error) {
	category = Category(strings.TrimSpace(string(category)))
	switch category {
	case CategoryImages, CategoryDocuments, CategoryVideos, CategoryAudios, CategoryOthers:
		return category, nil
	default:
		return "", apperror.ErrCommonInvalidValue
	}
}

func (s *QueriesSanitizer) GetFileInfosByCategory(ctx context.Context, query *GetFileInfosByCategoryQuery) (*paging.Page[*FileInfo], error) {
	if c, err := validateCategory(query.Category); err != nil {
		return nil, err
	} else {
		query.Category = c
	}

	return s.Queries.GetFileInfosByCategory(ctx, query)
}
