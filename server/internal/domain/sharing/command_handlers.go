package sharing

import (
	"context"
	"skyvault/pkg/appconfig"
	"skyvault/pkg/apperror"
)

var _ Commands = (*CommandHandlers)(nil)

type CommandHandlers struct {
	app        *appconfig.App
	repository Repository
}

func NewCommandHandlers(app *appconfig.App, repository Repository) *CommandHandlers {
	return &CommandHandlers{app: app, repository: repository}
}

//--------------------------------
// Contacts
//--------------------------------

func (h *CommandHandlers) CreateContact(ctx context.Context, cmd *CreateContactCommand) (*Contact, error) {
	contact, err := NewContact(cmd.OwnerID, cmd.Email, cmd.Name)
	if err != nil {
		return nil, apperror.NewAppError(err, "sharing.CommandHandlers.CreateContact:NewContact")
	}

	contact, err = h.repository.CreateContact(ctx, contact)
	if err != nil {
		return nil, apperror.NewAppError(err, "sharing.CommandHandlers.CreateContact:CreateContact")
	}

	return contact, nil
}

func (h *CommandHandlers) UpdateContact(ctx context.Context, cmd *UpdateContactCommand) error {
	contact, err := h.repository.GetContact(ctx, cmd.OwnerID, cmd.ContactID)
	if err != nil {
		return apperror.NewAppError(err, "sharing.CommandHandlers.UpdateContact:GetContact")
	}

	contact.UpdateName(cmd.Name)

	err = h.repository.UpdateContact(ctx, contact)
	if err != nil {
		return apperror.NewAppError(err, "sharing.CommandHandlers.UpdateContact:UpdateContact")
	}

	return nil
}

func (h *CommandHandlers) DeleteContact(ctx context.Context, cmd *DeleteContactCommand) error {
	err := h.repository.DeleteContact(ctx, cmd.OwnerID, cmd.ContactID)
	if err != nil {
		return apperror.NewAppError(err, "sharing.CommandHandlers.DeleteContact:DeleteContact")
	}

	return nil
}

//--------------------------------
// Contact Groups
//--------------------------------

func (h *CommandHandlers) CreateContactGroup(ctx context.Context, cmd *CreateContactGroupCommand) (*ContactGroup, error) {
	group, err := NewContactGroup(cmd.OwnerID, cmd.Name)
	if err != nil {
		return nil, apperror.NewAppError(err, "sharing.CommandHandlers.CreateContactGroup:NewContactGroup")
	}

	group, err = h.repository.CreateContactGroup(ctx, group)
	if err != nil {
		return nil, apperror.NewAppError(err, "sharing.CommandHandlers.CreateContactGroup:CreateContactGroup")
	}

	return group, nil
}

func (h *CommandHandlers) RenameContactGroup(ctx context.Context, cmd *RenameContactGroupCommand) error {
	group, err := h.repository.GetContactGroup(ctx, cmd.OwnerID, cmd.GroupID)
	if err != nil {
		return apperror.NewAppError(err, "sharing.CommandHandlers.RenameContactGroup:GetContactGroup")
	}

	err = group.Rename(cmd.NewName)
	if err != nil {
		return apperror.NewAppError(err, "sharing.CommandHandlers.RenameContactGroup:Rename")
	}

	err = h.repository.UpdateContactGroup(ctx, group)
	if err != nil {
		return apperror.NewAppError(err, "sharing.CommandHandlers.RenameContactGroup:UpdateContactGroup")
	}

	return nil
}

func (h *CommandHandlers) DeleteContactGroup(ctx context.Context, cmd *DeleteContactGroupCommand) error {
	err := h.repository.DeleteContactGroup(ctx, cmd.OwnerID, cmd.GroupID)
	if err != nil {
		return apperror.NewAppError(err, "sharing.CommandHandlers.DeleteContactGroup:DeleteContactGroup")
	}

	return nil
}

func (h *CommandHandlers) AddContactToGroup(ctx context.Context, cmd *AddContactToGroupCommand) error {
	err := h.repository.AddContactToGroup(ctx, cmd.OwnerID, cmd.GroupID, cmd.ContactID)
	if err != nil {
		return apperror.NewAppError(err, "sharing.CommandHandlers.AddContactToGroup:AddContactToGroup")
	}

	return nil
}

func (h *CommandHandlers) RemoveContactFromGroup(ctx context.Context, cmd *RemoveContactFromGroupCommand) error {
	err := h.repository.RemoveContactFromGroup(ctx, cmd.OwnerID, cmd.GroupID, cmd.ContactID)
	if err != nil {
		return apperror.NewAppError(err, "sharing.CommandHandlers.RemoveContactFromGroup:RemoveContactFromGroup")
	}

	return nil
}

//--------------------------------
// Sharing
//--------------------------------

func (h *CommandHandlers) CreateShare(ctx context.Context, cmd *CreateShareCommand) (*ShareConfig, error) {
	config, err := NewShareConfig(
		cmd.OwnerID,
		cmd.ResourceType,
		cmd.ResourceID,
		cmd.PasswordHash,
		cmd.MaxDownloads,
		cmd.ExpiresAt,
	)
	if err != nil {
		return nil, apperror.NewAppError(err, "sharing.CommandHandlers.CreateShare:NewShareConfig")
	}

	tx, err := h.repository.BeginTx(ctx)
	if err != nil {
		return nil, apperror.NewAppError(err, "sharing.CommandHandlers.CreateShare:BeginTx")
	}
	defer tx.Rollback()

	repoTx := h.repository.WithTx(ctx, tx)

	config, err = repoTx.CreateShareConfig(ctx, config)
	if err != nil {
		return nil, apperror.NewAppError(err, "sharing.CommandHandlers.CreateShare:CreateShareConfig")
	}

	for _, r := range cmd.Recipients {
		recipient, err := NewShareRecipient(
			config.ID,
			r.Type,
			r.RecipientID,
			r.Email,
		)
		if err != nil {
			return nil, apperror.NewAppError(err, "sharing.CommandHandlers.CreateShare:NewShareRecipient")
		}

		recipient, err = repoTx.CreateShareRecipient(ctx, recipient)
		if err != nil {
			return nil, apperror.NewAppError(err, "sharing.CommandHandlers.CreateShare:CreateShareRecipient")
		}

		config.Recipients = append(config.Recipients, recipient)
	}

	err = tx.Commit()
	if err != nil {
		return nil, apperror.NewAppError(err, "sharing.CommandHandlers.CreateShare:Commit")
	}

	return config, nil
}

func (h *CommandHandlers) UpdateShareConfig(ctx context.Context, cmd *UpdateShareConfigCommand) error {
	config, err := h.repository.GetShareConfig(ctx, cmd.ShareID)
	if err != nil {
		return apperror.NewAppError(err, "sharing.CommandHandlers.UpdateShareConfig:GetShareConfig")
	}

	err = config.ValidateAccess(cmd.OwnerID)
	if err != nil {
		return apperror.NewAppError(err, "sharing.CommandHandlers.UpdateShareConfig:ValidateAccess")
	}

	err = config.UpdateConfig(cmd.PasswordHash, cmd.MaxDownloads, cmd.ExpiresAt)
	if err != nil {
		return apperror.NewAppError(err, "sharing.CommandHandlers.UpdateShareConfig:UpdateConfig")
	}

	err = h.repository.UpdateShareConfig(ctx, config)
	if err != nil {
		return apperror.NewAppError(err, "sharing.CommandHandlers.UpdateShareConfig:UpdateShareConfig")
	}

	return nil
}

func (h *CommandHandlers) DeleteShare(ctx context.Context, cmd *DeleteShareCommand) error {
	config, err := h.repository.GetShareConfig(ctx, cmd.ShareID)
	if err != nil {
		return apperror.NewAppError(err, "sharing.CommandHandlers.DeleteShare:GetShareConfig")
	}

	err = config.ValidateAccess(cmd.OwnerID)
	if err != nil {
		return apperror.NewAppError(err, "sharing.CommandHandlers.DeleteShare:ValidateAccess")
	}

	err = h.repository.DeleteShareConfig(ctx, cmd.OwnerID, cmd.ShareID)
	if err != nil {
		return apperror.NewAppError(err, "sharing.CommandHandlers.DeleteShare:DeleteShareConfig")
	}

	return nil
}

func (h *CommandHandlers) AddShareRecipient(ctx context.Context, cmd *AddShareRecipientCommand) error {
	config, err := h.repository.GetShareConfig(ctx, cmd.ShareID)
	if err != nil {
		return apperror.NewAppError(err, "sharing.CommandHandlers.AddShareRecipient:GetShareConfig")
	}

	err = config.ValidateAccess(cmd.OwnerID)
	if err != nil {
		return apperror.NewAppError(err, "sharing.CommandHandlers.AddShareRecipient:ValidateAccess")
	}

	recipient, err := NewShareRecipient(
		config.ID,
		cmd.Type,
		cmd.RecipientID,
		cmd.Email,
	)
	if err != nil {
		return apperror.NewAppError(err, "sharing.CommandHandlers.AddShareRecipient:NewShareRecipient")
	}

	_, err = h.repository.CreateShareRecipient(ctx, recipient)
	if err != nil {
		return apperror.NewAppError(err, "sharing.CommandHandlers.AddShareRecipient:CreateShareRecipient")
	}

	return nil
}

func (h *CommandHandlers) RemoveShareRecipient(ctx context.Context, cmd *RemoveShareRecipientCommand) error {
	config, err := h.repository.GetShareConfig(ctx, cmd.ShareID)
	if err != nil {
		return apperror.NewAppError(err, "sharing.CommandHandlers.RemoveShareRecipient:GetShareConfig")
	}

	err = config.ValidateAccess(cmd.OwnerID)
	if err != nil {
		return apperror.NewAppError(err, "sharing.CommandHandlers.RemoveShareRecipient:ValidateAccess")
	}

	err = h.repository.DeleteShareRecipient(ctx, cmd.ShareID, cmd.RecipientID)
	if err != nil {
		return apperror.NewAppError(err, "sharing.CommandHandlers.RemoveShareRecipient:DeleteShareRecipient")
	}

	return nil
}
