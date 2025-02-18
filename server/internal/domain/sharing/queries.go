package sharing

import (
	"context"
	"skyvault/pkg/paging"
)

type Queries interface {
	//--------------------------------
	// Contacts
	//--------------------------------

	// App Errors:
	// - ErrCommonNoData
	GetContact(ctx context.Context, query *GetContactQuery) (*Contact, error)

	GetContacts(ctx context.Context, query *GetContactsQuery) (*paging.Page[*Contact], error)

	//--------------------------------
	// Contact Groups
	//--------------------------------

	// App Errors:
	// - ErrCommonNoData
	GetContactGroup(ctx context.Context, query *GetContactGroupQuery) (*ContactGroup, error)

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

type GetContactQuery struct {
	OwnerID   int64
	ContactID int64
}

type GetContactsQuery struct {
	OwnerID    int64
	SearchTerm *string
	PagingOpt  *paging.Options
}

//--------------------------------
// Contact Groups
//--------------------------------

type GetContactGroupQuery struct {
	OwnerID int64
	GroupID int64
}

type GetContactGroupsQuery struct {
	OwnerID    int64
	SearchTerm *string
	PagingOpt  *paging.Options
}

type GetContactGroupMembersQuery struct {
	OwnerID   int64
	GroupID   int64
	PagingOpt *paging.Options
}

//--------------------------------
// Sharing
//--------------------------------

type GetShareConfigQuery struct {
	ShareID int64
}

type GetShareConfigsQuery struct {
	OwnerID      int64
	ResourceType *ResourceType
	PagingOpt    *paging.Options
}

type ValidateShareAccessQuery struct {
	ShareID     int64
	RecipientID int64
	Password    *string
}
