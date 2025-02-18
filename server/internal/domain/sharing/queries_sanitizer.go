package sharing

import (
	"context"
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

func (s *QueriesSanitizer) GetContacts(ctx context.Context, query *GetContactsQuery) (*paging.Page[*Contact], error) {
	if query.SearchTerm != nil {
		trimmed := strings.TrimSpace(*query.SearchTerm)
		if trimmed == "" {
			query.SearchTerm = nil
		} else {
			query.SearchTerm = &trimmed
		}
	}

	return s.Queries.GetContacts(ctx, query)
}

func (s *QueriesSanitizer) GetContactGroups(ctx context.Context, query *GetContactGroupsQuery) (*paging.Page[*ContactGroup], error) {
	if query.SearchTerm != nil {
		trimmed := strings.TrimSpace(*query.SearchTerm)
		if trimmed == "" {
			query.SearchTerm = nil
		} else {
			query.SearchTerm = &trimmed
		}
	}

	return s.Queries.GetContactGroups(ctx, query)
}
