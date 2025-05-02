package sharing

import (
	"context"
	"skyvault/internal/domain/internal"
	"skyvault/pkg/paging"
	"time"
)

type Repository interface {
	internal.RepositoryTx[Repository]

	//--------------------------------
	// Contacts
	//--------------------------------

	// App Errors:
	// - ErrCommonDuplicateData
	CreateContact(ctx context.Context, contact *Contact) (*Contact, error)

	GetContacts(ctx context.Context, ownerID int64, pagingOpt *paging.Options, searchTerm *string) (*paging.Page[*Contact], error)

	// App Errors:
	// - ErrCommonNoData
	UpdateContact(ctx context.Context, ownerID, contactID int64, email, name string) error

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

	GetContactGroups(ctx context.Context, ownerID int64, pagingOpt *paging.Options, searchTerm *string) (*paging.Page[*ContactGroup], error)

	GetContactGroupMembers(ctx context.Context, ownerID int64, pagingOpt *paging.Options, groupID int64) (*paging.Page[*Contact], error)

	// App Errors:
	// - ErrCommonNoData
	UpdateContactGroup(ctx context.Context, ownerID, groupID int64, name string) error

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
	GetShareConfig(ctx context.Context, ownerID, shareID int64) (*ShareConfig, error)

	// App Errors:
	// - ErrCommonNoData
	GetShareConfigByCustomID(ctx context.Context, ownerID int64, customID string) (*ShareConfig, error)

	// App Errors:
	// - ErrCommonNoData
	UpdateShareExpiry(ctx context.Context, ownerID, shareID int64, maxDownloads *int64, expiresAt *time.Time) error

	// App Errors:
	// - ErrCommonNoData
	UpdateSharePassword(ctx context.Context, ownerID, shareID int64, password *string) error

	// App Errors:
	// - ErrCommonNoData
	DeleteShareConfig(ctx context.Context, ownerID, shareID int64) error

	// App Errors:
	// - ErrCommonDuplicateData
	CreateShareRecipient(ctx context.Context, ownerID int64, recipient *ShareRecipient) (*ShareRecipient, error)

	// App Errors:
	// - ErrCommonNoData
	DeleteShareRecipient(ctx context.Context, ownerID, shareID, recipientID int64) error

	// App Errors:
	// - ErrCommonNoData
	GetShareRecipientByEmail(ctx context.Context, shareID int64, email string) (*ShareRecipient, error)
}
