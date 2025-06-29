package sharing

import (
	"context"
	"skyvault/pkg/paging"
)

type Queries interface {
	//--------------------------------
	// Contacts
	//--------------------------------

	GetContacts(ctx context.Context, query *GetContactsQuery) (*paging.Page[*Contact], error)

	//--------------------------------
	// Contact Groups
	//--------------------------------

	GetContactGroups(ctx context.Context, query *GetContactGroupsQuery) (*paging.Page[*ContactGroup], error)

	GetContactGroupMembers(ctx context.Context, query *GetContactGroupMembersQuery) (*paging.Page[*Contact], error)

	//--------------------------------
	// Sharing
	//--------------------------------

	// App Errors:
	// - ErrCommonNoData
	GetShareConfig(ctx context.Context, query *GetShareConfigQuery) (*ShareConfig, error)

	// App Errors:
	// - ErrCommonNoData
	// - ErrCommonNoAccess
	// - ErrSharingExpired
	// - ErrSharingMaxDownloadsReached
	ValidateShareAccess(ctx context.Context, query *ValidateShareAccessQuery) error
}

//--------------------------------
// Contacts
//--------------------------------

type GetContactsQuery struct {
	SearchTerm *string
	PagingOpt  *paging.Options
}

//--------------------------------
// Contact Groups
//--------------------------------

type GetContactGroupsQuery struct {
	SearchTerm *string
	PagingOpt  *paging.Options
}

type GetContactGroupMembersQuery struct {
	GroupID   string
	PagingOpt *paging.Options
}

//--------------------------------
// Sharing
//--------------------------------

type GetShareConfigQuery struct {
	ShareID string
}

type ValidateShareAccessQuery struct {
	ShareID  string
	Email    *string
	Password *string
}
