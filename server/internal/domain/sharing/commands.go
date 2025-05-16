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
	ContactID string
	Name      string
	Email     string
}

type DeleteContactCommand struct {
	ContactID string
}

//--------------------------------
// Contact Groups
//--------------------------------

type CreateContactGroupCommand struct {
	Name string
}

type RenameContactGroupCommand struct {
	GroupID string
	NewName string
}

type DeleteContactGroupCommand struct {
	GroupID string
}

type AddContactToGroupCommand struct {
	GroupID   string
	ContactID string
}

type RemoveContactFromGroupCommand struct {
	GroupID   string
	ContactID string
}

//--------------------------------
// Sharing
//--------------------------------

// Only one of FileID or FolderID must be set.
type CreateShareCommand struct {
	FileID       *string
	FolderID     *string
	Recipients   []*ShareRecipientInput
	Password     *string
	MaxDownloads *int64
	ExpiresAt    *time.Time
}

// Only one of ContactID, ContactGroupID or Email can be set.
// If SaveAsContact is true, then Email and Name must be set. Otherwise, Name is ignored.
type ShareRecipientInput struct {
	ContactID      *string
	ContactGroupID *string
	Email          *string
	Name           *string
	SaveAsContact  bool
}

type UpdateShareExpiryCommand struct {
	ShareID      string
	MaxDownloads *int64
	ExpiresAt    *time.Time
}

type UpdateSharePasswordCommand struct {
	ShareID  string
	Password *string
}

type DeleteShareCommand struct {
	ShareID string
}

// Only one of ContactID, ContactGroupID or Email can be set.
// If SaveAsContact is true, then Email and Name must be set. Otherwise, Name is ignored.
type AddShareRecipientCommand struct {
	ShareID string
	*ShareRecipientInput
}

type RemoveShareRecipientCommand struct {
	ShareID     string
	RecipientID string
}
