package sharing

import (
	"context"
	"skyvault/pkg/apperror"
	"skyvault/pkg/validate"
	"time"
)

var _ Commands = (*CommandsSanitizer)(nil)

type CommandsSanitizer struct {
	Commands
}

func NewCommandsSanitizer(commands Commands) *CommandsSanitizer {
	return &CommandsSanitizer{Commands: commands}
}

//--------------------------------
// Contacts
//--------------------------------

func (s *CommandsSanitizer) CreateContact(ctx context.Context, cmd *CreateContactCommand) (*Contact, error) {
	if e, err := validate.Email(cmd.Email); err != nil {
		return nil, apperror.NewAppError(err, "sharing.CommandsSanitizer.CreateContact:Email")
	} else {
		cmd.Email = e
	}

	if n, err := validate.Name(cmd.Name); err != nil {
		return nil, apperror.NewAppError(err, "sharing.CommandsSanitizer.CreateContact:Name")
	} else {
		cmd.Name = n
	}

	return s.Commands.CreateContact(ctx, cmd)
}

func (s *CommandsSanitizer) UpdateContact(ctx context.Context, cmd *UpdateContactCommand) error {
	if e, err := validate.Email(cmd.Email); err != nil {
		return apperror.NewAppError(err, "sharing.CommandsSanitizer.UpdateContact:Email")
	} else {
		cmd.Email = e
	}

	if n, err := validate.Name(cmd.Name); err != nil {
		return apperror.NewAppError(err, "sharing.CommandsSanitizer.UpdateContact:Name")
	} else {
		cmd.Name = n
	}

	return s.Commands.UpdateContact(ctx, cmd)
}

//--------------------------------
// Contact Groups
//--------------------------------

func (s *CommandsSanitizer) CreateContactGroup(ctx context.Context, cmd *CreateContactGroupCommand) (*ContactGroup, error) {
	if n, err := validate.Name(cmd.Name); err != nil {
		return nil, apperror.NewAppError(err, "sharing.CommandsSanitizer.CreateContactGroup:Name")
	} else {
		cmd.Name = n
	}

	return s.Commands.CreateContactGroup(ctx, cmd)
}

func (s *CommandsSanitizer) RenameContactGroup(ctx context.Context, cmd *RenameContactGroupCommand) error {
	if n, err := validate.Name(cmd.NewName); err != nil {
		return apperror.NewAppError(err, "sharing.CommandsSanitizer.RenameContactGroup:NewName")
	} else {
		cmd.NewName = n
	}

	return s.Commands.RenameContactGroup(ctx, cmd)
}

//--------------------------------
// Sharing
//--------------------------------

func (s *CommandsSanitizer) CreateShare(ctx context.Context, cmd *CreateShareCommand) (*ShareConfig, error) {
	if (cmd.FileID != nil && cmd.FolderID != nil) ||
		(cmd.FileID == nil && cmd.FolderID == nil) {
		return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "sharing.CommandsSanitizer.CreateShare:EitherFileIDOrFolderID")
	}

	if cmd.MaxDownloads != nil {
		if *cmd.MaxDownloads < 1 {
			return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "sharing.CommandsSanitizer.CreateShare:MaxDownloads")
		}
	}

	if cmd.ExpiresAt != nil {
		if time.Now().UTC().After(*cmd.ExpiresAt) {
			return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "sharing.CommandsSanitizer.CreateShare:ExpiresAt")
		}
	}

	if cmd.Password != nil {
		if p, err := validate.PasswordLen(*cmd.Password); err != nil {
			return nil, apperror.NewAppError(err, "sharing.CommandsSanitizer.CreateShare:Password")
		} else {
			cmd.Password = &p
		}
	}

	if len(cmd.Recipients) == 0 {
		return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "sharing.CommandsSanitizer.CreateShare:Recipients.Empty")
	}

	for _, recipient := range cmd.Recipients {
		if validShareRecipient(recipient) {
			return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "sharing.CommandsSanitizer.CreateShare:Recipient")
		}

		if recipient.Email != nil {
			if e, err := validate.Email(*recipient.Email); err != nil {
				return nil, apperror.NewAppError(err, "sharing.CommandsSanitizer.CreateShare:Email")
			} else {
				recipient.Email = &e
			}

			if recipient.SaveAsContact {
				if n, err := validate.Name(*recipient.Name); err != nil {
					return nil, apperror.NewAppError(err, "sharing.CommandsSanitizer.CreateShare:Name")
				} else {
					recipient.Name = &n
				}
			}
		}
	}

	return s.Commands.CreateShare(ctx, cmd)
}

func (s *CommandsSanitizer) UpdateShareExpiry(ctx context.Context, cmd *UpdateShareExpiryCommand) error {
	if cmd.MaxDownloads != nil {
		if *cmd.MaxDownloads < 1 {
			return apperror.NewAppError(apperror.ErrCommonInvalidValue, "sharing.CommandsSanitizer.UpdateShareExpiry:MaxDownloads")
		}
	}

	if cmd.ExpiresAt != nil {
		if time.Now().UTC().After(*cmd.ExpiresAt) {
			return apperror.NewAppError(apperror.ErrCommonInvalidValue, "sharing.CommandsSanitizer.UpdateShareExpiry:ExpiresAt")
		}
	}

	return s.Commands.UpdateShareExpiry(ctx, cmd)
}

func (s *CommandsSanitizer) UpdateSharePassword(ctx context.Context, cmd *UpdateSharePasswordCommand) error {
	if cmd.Password != nil {
		if p, err := validate.PasswordLen(*cmd.Password); err != nil {
			return apperror.NewAppError(err, "sharing.CommandsSanitizer.UpdateSharePassword:Password")
		} else {
			cmd.Password = &p
		}
	}

	return s.Commands.UpdateSharePassword(ctx, cmd)
}

func (s *CommandsSanitizer) AddShareRecipient(ctx context.Context, cmd *AddShareRecipientCommand) (*ShareRecipient, error) {
	if validShareRecipient(cmd.ShareRecipientInput) {
		return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "sharing.CommandsSanitizer.AddShareRecipient:Recipient")
	}

	if cmd.Email != nil {
		if e, err := validate.Email(*cmd.Email); err != nil {
			return nil, apperror.NewAppError(err, "sharing.CommandsSanitizer.AddShareRecipient:Email")
		} else {
			cmd.Email = &e
		}

		if cmd.SaveAsContact {
			if cmd.Name != nil {
				if n, err := validate.Name(*cmd.Name); err != nil {
					return nil, apperror.NewAppError(err, "sharing.CommandsSanitizer.AddShareRecipient:Name")
				} else {
					cmd.Name = &n
				}
			}
		}
	}

	return s.Commands.AddShareRecipient(ctx, cmd)
}

// Only one of the contactID, contactGroupID, or email can be set.
func validShareRecipient(recipient *ShareRecipientInput) bool {
	if recipient.SaveAsContact {
		if recipient.Email == nil || recipient.Name == nil || recipient.ContactID != nil || recipient.ContactGroupID != nil {
			return false
		}
	} else {
		if (recipient.ContactID == nil && recipient.ContactGroupID == nil && recipient.Email == nil) ||
			(recipient.ContactID != nil && (recipient.ContactGroupID != nil || recipient.Email != nil)) ||
			(recipient.ContactGroupID != nil && (recipient.ContactID != nil || recipient.Email != nil)) ||
			(recipient.Email != nil && (recipient.ContactID != nil || recipient.ContactGroupID != nil)) {
			return false
		}
	}
	return true
}
