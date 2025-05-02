package sharing

import (
	"context"
	"skyvault/pkg/apperror"
	"skyvault/pkg/paging"
	"skyvault/pkg/validate"
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
		if n, err := validate.Name(*query.SearchTerm); err != nil {
			return nil, apperror.NewAppError(err, "sharing.QueriesSanitizer.GetContacts:SearchTerm")
		} else {
			query.SearchTerm = &n
		}
	}

	return s.Queries.GetContacts(ctx, query)
}

func (s *QueriesSanitizer) GetContactGroups(ctx context.Context, query *GetContactGroupsQuery) (*paging.Page[*ContactGroup], error) {
	if query.SearchTerm != nil {
		if n, err := validate.Name(*query.SearchTerm); err != nil {
			return nil, apperror.NewAppError(err, "sharing.QueriesSanitizer.GetContactGroups:SearchTerm")
		} else {
			query.SearchTerm = &n
		}
	}

	return s.Queries.GetContactGroups(ctx, query)
}

func (s *QueriesSanitizer) ValidateShareAccess(ctx context.Context, query *ValidateShareAccessQuery) error {
	if query.Email != nil {
		if m, err := validate.Email(*query.Email); err != nil {
			return apperror.NewAppError(err, "sharing.QueriesSanitizer.ValidateShareAccess:Email")
		} else {
			query.Email = &m
		}
	}

	if query.Password != nil {
		if p, err := validate.PasswordLen(*query.Password); err != nil {
			return apperror.NewAppError(err, "sharing.QueriesSanitizer.ValidateShareAccess:Password")
		} else {
			query.Password = &p
		}
	}

	return s.Queries.ValidateShareAccess(ctx, query)
}
