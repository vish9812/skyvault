package sharing

import (
	"context"
	"time"
)

type Commands interface {
	//--------------------------------
	// Contacts
	//--------------------------------

	// App Errors:
	// - ErrCommonDuplicateData
	// - ErrCommonInvalidValue
	CreateContact(ctx context.Context, cmd *CreateContactCommand) (*Contact, error)

	// App Errors:
	// - ErrCommonNoData
	UpdateContact(ctx context.Context, cmd *UpdateContactCommand) error

	// App Errors:
	// - ErrCommonNoData
	DeleteContact(ctx context.Context, cmd *DeleteContactCommand) error

	//--------------------------------
	// Contact Groups
	//--------------------------------

	// App Errors:
	// - ErrCommonDuplicateData
	// - ErrCommonInvalidValue
	CreateContactGroup(ctx context.Context, cmd *CreateContactGroupCommand) (*ContactGroup, error)

	// App Errors:
	// - ErrCommonNoData
	// - ErrCommonInvalidValue
	RenameContactGroup(ctx context.Context, cmd *RenameContactGroupCommand) error

	// App Errors:
	// - ErrCommonNoData
	DeleteContactGroup(ctx context.Context, cmd *DeleteContactGroupCommand) error

	// App Errors:
	// - ErrCommonNoData
	// - ErrCommonDuplicateData
	AddContactToGroup(ctx context.Context, cmd *AddContactToGroupCommand) error

	// App Errors:
	// - ErrCommonNoData
	RemoveContactFromGroup(ctx context.Context, cmd *RemoveContactFromGroupCommand) error

	//--------------------------------
	// Sharing
	//--------------------------------

	// App Errors:
	// - ErrCommonInvalidValue
	// - ErrCommonNoData (if resource not found)
	CreateShare(ctx context.Context, cmd *CreateShareCommand) (*ShareConfig, error)

	// App Errors:
	// - ErrCommonNoData
	// - ErrCommonNoAccess
	UpdateShareConfig(ctx context.Context, cmd *UpdateShareConfigCommand) error

	// App Errors:
	// - ErrCommonNoData
	// - ErrCommonNoAccess
	DeleteShare(ctx context.Context, cmd *DeleteShareCommand) error

	// App Errors:
	// - ErrCommonNoData
	// - ErrCommonNoAccess
	// - ErrCommonDuplicateData
	AddShareRecipient(ctx context.Context, cmd *AddShareRecipientCommand) error

	// App Errors:
	// - ErrCommonNoData
	// - ErrCommonNoAccess
	RemoveShareRecipient(ctx context.Context, cmd *RemoveShareRecipientCommand) error
}

//--------------------------------
// Contacts
//--------------------------------

type CreateContactCommand struct {
	OwnerID int64
	Email   string
	Name    *string
}

type UpdateContactCommand struct {
	OwnerID   int64
	ContactID int64
	Name      *string
}

type DeleteContactCommand struct {
	OwnerID   int64
	ContactID int64
}

//--------------------------------
// Contact Groups
//--------------------------------

type CreateContactGroupCommand struct {
	OwnerID int64
	Name    string
}

type RenameContactGroupCommand struct {
	OwnerID int64
	GroupID int64
	NewName string
}

type DeleteContactGroupCommand struct {
	OwnerID int64
	GroupID int64
}

type AddContactToGroupCommand struct {
	OwnerID   int64
	GroupID   int64
	ContactID int64
}

type RemoveContactFromGroupCommand struct {
	OwnerID   int64
	GroupID   int64
	ContactID int64
}

//--------------------------------
// Sharing
//--------------------------------

type CreateShareCommand struct {
	OwnerID      int64
	ResourceType ResourceType
	ResourceID   int64
	Recipients   []ShareRecipientInput
	PasswordHash *string
	MaxDownloads *int
	ExpiresAt    *time.Time
}

type ShareRecipientInput struct {
	Type        RecipientType
	RecipientID *int64 // For groups
	Email       string // For direct email shares
}

type UpdateShareConfigCommand struct {
	OwnerID      int64
	ShareID      int64
	PasswordHash *string
	MaxDownloads *int
	ExpiresAt    *time.Time
}

type DeleteShareCommand struct {
	OwnerID int64
	ShareID int64
}

type AddShareRecipientCommand struct {
	OwnerID     int64
	ShareID     int64
	Type        RecipientType
	RecipientID *int64
	Email       string
}

type RemoveShareRecipientCommand struct {
	OwnerID     int64
	ShareID     int64
	RecipientID int64
}
