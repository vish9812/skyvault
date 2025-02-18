package sharing

import (
	"context"
	"skyvault/pkg/apperror"
	"strings"
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
	cmd.Email = strings.TrimSpace(strings.ToLower(cmd.Email))
	if cmd.Email == "" {
		return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "sharing.CommandsSanitizer.CreateContact:Email")
	}

	if cmd.Name != nil {
		trimmedName := strings.TrimSpace(*cmd.Name)
		if trimmedName == "" {
			cmd.Name = nil
		} else {
			cmd.Name = &trimmedName
		}
	}

	return s.Commands.CreateContact(ctx, cmd)
}

func (s *CommandsSanitizer) UpdateContact(ctx context.Context, cmd *UpdateContactCommand) error {
	if cmd.Name != nil {
		trimmedName := strings.TrimSpace(*cmd.Name)
		if trimmedName == "" {
			cmd.Name = nil
		} else {
			cmd.Name = &trimmedName
		}
	}

	return s.Commands.UpdateContact(ctx, cmd)
}

//--------------------------------
// Contact Groups
//--------------------------------

func (s *CommandsSanitizer) CreateContactGroup(ctx context.Context, cmd *CreateContactGroupCommand) (*ContactGroup, error) {
	cmd.Name = strings.TrimSpace(cmd.Name)
	if cmd.Name == "" {
		return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "sharing.CommandsSanitizer.CreateContactGroup:Name")
	}

	return s.Commands.CreateContactGroup(ctx, cmd)
}

func (s *CommandsSanitizer) RenameContactGroup(ctx context.Context, cmd *RenameContactGroupCommand) error {
	cmd.NewName = strings.TrimSpace(cmd.NewName)
	if cmd.NewName == "" {
		return apperror.NewAppError(apperror.ErrCommonInvalidValue, "sharing.CommandsSanitizer.RenameContactGroup:NewName")
	}

	return s.Commands.RenameContactGroup(ctx, cmd)
}

//--------------------------------
// Sharing
//--------------------------------

func (s *CommandsSanitizer) CreateShare(ctx context.Context, cmd *CreateShareCommand) (*ShareConfig, error) {
	switch cmd.ResourceType {
	case ResourceTypeFile, ResourceTypeFolder:
		// valid
	default:
		return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "sharing.CommandsSanitizer.CreateShare:ResourceType")
	}

	if len(cmd.Recipients) == 0 {
		return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "sharing.CommandsSanitizer.CreateShare:Recipients")
	}

	for i, r := range cmd.Recipients {
		switch r.Type {
		case RecipientTypeEmail:
			if r.RecipientID != nil {
				return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "sharing.CommandsSanitizer.CreateShare:RecipientID")
			}
			r.Email = strings.TrimSpace(strings.ToLower(r.Email))
			if r.Email == "" {
				return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "sharing.CommandsSanitizer.CreateShare:Email")
			}
		case RecipientTypeGroup:
			if r.RecipientID == nil {
				return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "sharing.CommandsSanitizer.CreateShare:RecipientID")
			}
		default:
			return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "sharing.CommandsSanitizer.CreateShare:RecipientType")
		}
		cmd.Recipients[i] = r
	}

	if cmd.MaxDownloads != nil && *cmd.MaxDownloads <= 0 {
		return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "sharing.CommandsSanitizer.CreateShare:MaxDownloads")
	}

	return s.Commands.CreateShare(ctx, cmd)
}

func (s *CommandsSanitizer) AddShareRecipient(ctx context.Context, cmd *AddShareRecipientCommand) error {
	switch cmd.Type {
	case RecipientTypeEmail:
		if cmd.RecipientID != nil {
			return apperror.NewAppError(apperror.ErrCommonInvalidValue, "sharing.CommandsSanitizer.AddShareRecipient:RecipientID")
		}
		cmd.Email = strings.TrimSpace(strings.ToLower(cmd.Email))
		if cmd.Email == "" {
			return apperror.NewAppError(apperror.ErrCommonInvalidValue, "sharing.CommandsSanitizer.AddShareRecipient:Email")
		}
	case RecipientTypeGroup:
		if cmd.RecipientID == nil {
			return apperror.NewAppError(apperror.ErrCommonInvalidValue, "sharing.CommandsSanitizer.AddShareRecipient:RecipientID")
		}
	default:
		return apperror.NewAppError(apperror.ErrCommonInvalidValue, "sharing.CommandsSanitizer.AddShareRecipient:RecipientType")
	}

	return s.Commands.AddShareRecipient(ctx, cmd)
}
