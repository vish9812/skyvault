package sharing

import (
	"context"
	"skyvault/internal/domain/internal"
	"skyvault/pkg/paging"
)

type Repository interface {
	internal.RepositoryTx[Repository]

	//--------------------------------
	// Contacts
	//--------------------------------

	// App Errors:
	// - ErrCommonDuplicateData
	CreateContact(ctx context.Context, contact *Contact) (*Contact, error)

	// App Errors:
	// - ErrCommonNoData
	GetContact(ctx context.Context, ownerID, contactID int64) (*Contact, error)

	GetContacts(ctx context.Context, pagingOpt *paging.Options, ownerID int64, searchTerm *string) (*paging.Page[*Contact], error)

	// App Errors:
	// - ErrCommonNoData
	UpdateContact(ctx context.Context, contact *Contact) error

	// App Errors:
	// - ErrCommonNoData
	DeleteContact(ctx context.Context, ownerID, contactID int64) error

	//--------------------------------
	// Contact Groups
	//--------------------------------

	// App Errors:
	// - ErrCommonDuplicateData
	CreateContactGroup(ctx context.Context, group *ContactGroup) (*ContactGroup, error)

	// App Errors:
	// - ErrCommonNoData
	GetContactGroup(ctx context.Context, ownerID, groupID int64) (*ContactGroup, error)

	GetContactGroups(ctx context.Context, pagingOpt *paging.Options, ownerID int64, searchTerm *string) (*paging.Page[*ContactGroup], error)

	GetContactGroupMembers(ctx context.Context, pagingOpt *paging.Options, ownerID, groupID int64) (*paging.Page[*Contact], error)

	// App Errors:
	// - ErrCommonNoData
	UpdateContactGroup(ctx context.Context, group *ContactGroup) error

	// App Errors:
	// - ErrCommonNoData
	DeleteContactGroup(ctx context.Context, ownerID, groupID int64) error

	// App Errors:
	// - ErrCommonNoData
	// - ErrCommonDuplicateData
	AddContactToGroup(ctx context.Context, ownerID, groupID, contactID int64) error

	// App Errors:
	// - ErrCommonNoData
	RemoveContactFromGroup(ctx context.Context, ownerID, groupID, contactID int64) error

	//--------------------------------
	// Share Configs
	//--------------------------------

	// App Errors:
	// - ErrCommonDuplicateData
	CreateShareConfig(ctx context.Context, config *ShareConfig) (*ShareConfig, error)

	// App Errors:
	// - ErrCommonNoData
	GetShareConfig(ctx context.Context, shareID int64) (*ShareConfig, error)

	// App Errors:
	// - ErrCommonNoData
	UpdateShareConfig(ctx context.Context, config *ShareConfig) error

	// App Errors:
	// - ErrCommonNoData
	DeleteShareConfig(ctx context.Context, ownerID, shareID int64) error

	//--------------------------------
	// Share Recipients
	//--------------------------------

	// App Errors:
	// - ErrCommonDuplicateData
	CreateShareRecipient(ctx context.Context, recipient *ShareRecipient) (*ShareRecipient, error)

	// App Errors:
	// - ErrCommonNoData
	GetShareRecipient(ctx context.Context, shareID, recipientID int64) (*ShareRecipient, error)

	// App Errors:
	// - ErrCommonNoData
	UpdateShareRecipient(ctx context.Context, recipient *ShareRecipient) error

	// App Errors:
	// - ErrCommonNoData
	DeleteShareRecipient(ctx context.Context, shareID, recipientID int64) error

	//--------------------------------
	// Share Access
	//--------------------------------

	// App Errors:
	// - ErrCommonNoData
	RecordShareAccess(ctx context.Context, shareID int64, accessedFromIP string) error
}
