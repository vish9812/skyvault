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
	// - ErrCommonInvalidValue
	// - ErrCommonNoData
	// - ErrCommonNoAccess
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
	UpdateShareExpiry(ctx context.Context, cmd *UpdateShareExpiryCommand) error

	// App Errors:
	// - ErrCommonNoData
	// - ErrCommonNoAccess
	UpdateSharePassword(ctx context.Context, cmd *UpdateSharePasswordCommand) error

	// App Errors:
	// - ErrCommonNoData
	// - ErrCommonNoAccess
	DeleteShare(ctx context.Context, cmd *DeleteShareCommand) error

	// App Errors:
	// - ErrCommonNoData
	// - ErrCommonNoAccess
	// - ErrCommonDuplicateData
	AddShareRecipient(ctx context.Context, cmd *AddShareRecipientCommand) (*ShareRecipient, error)

	// App Errors:
	// - ErrCommonNoData
	// - ErrCommonNoAccess
	RemoveShareRecipient(ctx context.Context, cmd *RemoveShareRecipientCommand) error
}

//--------------------------------
// Contacts
//--------------------------------

type CreateContactCommand struct {
	Email string
	Name  string
}

type UpdateContactCommand struct {
	ContactID int64
	Name      string
	Email     string
}

type DeleteContactCommand struct {
	ContactID int64
}

//--------------------------------
// Contact Groups
//--------------------------------

type CreateContactGroupCommand struct {
	Name string
}

type RenameContactGroupCommand struct {
	GroupID int64
	NewName string
}

type DeleteContactGroupCommand struct {
	GroupID int64
}

type AddContactToGroupCommand struct {
	GroupID   int64
	ContactID int64
}

type RemoveContactFromGroupCommand struct {
	GroupID   int64
	ContactID int64
}

//--------------------------------
// Sharing
//--------------------------------

// Only one of FileID or FolderID must be set.
type CreateShareCommand struct {
	FileID       *int64
	FolderID     *int64
	Recipients   []*ShareRecipientInput
	Password     *string
	MaxDownloads *int64
	ExpiresAt    *time.Time
}

// Only one of ContactID, ContactGroupID or Email can be set.
// If SaveAsContact is true, then Email and Name must be set. Otherwise, Name is ignored.
type ShareRecipientInput struct {
	ContactID      *int64
	ContactGroupID *int64
	Email          *string
	Name           *string
	SaveAsContact  bool
}

type UpdateShareExpiryCommand struct {
	ShareID      int64
	MaxDownloads *int64
	ExpiresAt    *time.Time
}

type UpdateSharePasswordCommand struct {
	ShareID  int64
	Password *string
}

type DeleteShareCommand struct {
	ShareID int64
}

// Only one of ContactID, ContactGroupID or Email can be set.
// If SaveAsContact is true, then Email and Name must be set. Otherwise, Name is ignored.
type AddShareRecipientCommand struct {
	ShareID int64
	*ShareRecipientInput
}

type RemoveShareRecipientCommand struct {
	ShareID     int64
	RecipientID int64
}
