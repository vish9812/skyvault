package profile

var _ Queries = (*QueriesSanitizer)(nil)

type QueriesSanitizer struct {
	Queries
}

func NewQueriesSanitizer(queries Queries) Queries {
	return &QueriesSanitizer{Queries: queries}
}
