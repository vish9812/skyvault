package sharing

import (
	"context"
	"database/sql"
	"skyvault/internal/domain/media"
	"skyvault/pkg/appconfig"
	"skyvault/pkg/apperror"
	"skyvault/pkg/common"
	"skyvault/pkg/utils"
)

var _ Commands = (*CommandHandlers)(nil)

type CommandHandlers struct {
	app             *appconfig.App
	repository      Repository
	mediaRepository media.Repository
}

func NewCommandHandlers(app *appconfig.App, repository Repository, mediaRepository media.Repository) *CommandHandlers {
	return &CommandHandlers{app: app, repository: repository, mediaRepository: mediaRepository}
}

//--------------------------------
// Contacts
//--------------------------------

func (h *CommandHandlers) CreateContact(ctx context.Context, cmd *CreateContactCommand) (*Contact, error) {
	profileID := common.GetProfileIDFromContext(ctx)
	contact := NewContact(profileID, cmd.Email, cmd.Name)

	contact, err := h.repository.CreateContact(ctx, contact)
	if err != nil {
		return nil, apperror.NewAppError(err, "sharing.CommandHandlers.CreateContact:CreateContact")
	}

	return contact, nil
}

func (h *CommandHandlers) UpdateContact(ctx context.Context, cmd *UpdateContactCommand) error {
	profileID := common.GetProfileIDFromContext(ctx)
	err := h.repository.UpdateContact(ctx, profileID, cmd.ContactID, cmd.Email, cmd.Name)
	if err != nil {
		return apperror.NewAppError(err, "sharing.CommandHandlers.UpdateContact:UpdateContact")
	}

	return nil
}

func (h *CommandHandlers) DeleteContact(ctx context.Context, cmd *DeleteContactCommand) error {
	profileID := common.GetProfileIDFromContext(ctx)
	err := h.repository.DeleteContact(ctx, profileID, cmd.ContactID)
	if err != nil {
		return apperror.NewAppError(err, "sharing.CommandHandlers.DeleteContact:DeleteContact")
	}

	return nil
}

//--------------------------------
// Contact Groups
//--------------------------------

func (h *CommandHandlers) CreateContactGroup(ctx context.Context, cmd *CreateContactGroupCommand) (*ContactGroup, error) {
	profileID := common.GetProfileIDFromContext(ctx)
	group := NewContactGroup(profileID, cmd.Name)

	group, err := h.repository.CreateContactGroup(ctx, group)
	if err != nil {
		return nil, apperror.NewAppError(err, "sharing.CommandHandlers.CreateContactGroup:CreateContactGroup")
	}

	return group, nil
}

func (h *CommandHandlers) RenameContactGroup(ctx context.Context, cmd *RenameContactGroupCommand) error {
	profileID := common.GetProfileIDFromContext(ctx)
	err := h.repository.UpdateContactGroup(ctx, profileID, cmd.GroupID, cmd.NewName)
	if err != nil {
		return apperror.NewAppError(err, "sharing.CommandHandlers.RenameContactGroup:UpdateContactGroup")
	}

	return nil
}

func (h *CommandHandlers) DeleteContactGroup(ctx context.Context, cmd *DeleteContactGroupCommand) error {
	profileID := common.GetProfileIDFromContext(ctx)
	err := h.repository.DeleteContactGroup(ctx, profileID, cmd.GroupID)
	if err != nil {
		return apperror.NewAppError(err, "sharing.CommandHandlers.DeleteContactGroup:DeleteContactGroup")
	}

	return nil
}

func (h *CommandHandlers) AddContactToGroup(ctx context.Context, cmd *AddContactToGroupCommand) error {
	profileID := common.GetProfileIDFromContext(ctx)
	err := h.repository.AddContactToGroup(ctx, profileID, cmd.GroupID, cmd.ContactID)
	if err != nil {
		return apperror.NewAppError(err, "sharing.CommandHandlers.AddContactToGroup:AddContactToGroup")
	}

	return nil
}

func (h *CommandHandlers) RemoveContactFromGroup(ctx context.Context, cmd *RemoveContactFromGroupCommand) error {
	profileID := common.GetProfileIDFromContext(ctx)
	err := h.repository.RemoveContactFromGroup(ctx, profileID, cmd.GroupID, cmd.ContactID)
	if err != nil {
		return apperror.NewAppError(err, "sharing.CommandHandlers.RemoveContactFromGroup:RemoveContactFromGroup")
	}

	return nil
}

//--------------------------------
// Sharing
//--------------------------------

func (h *CommandHandlers) CreateShare(ctx context.Context, cmd *CreateShareCommand) (*ShareConfig, error) {
	profileID := common.GetProfileIDFromContext(ctx)

	if cmd.FileID != nil {
		fileInfo, err := h.mediaRepository.GetFileInfo(ctx, *cmd.FileID)
		if err != nil {
			return nil, apperror.NewAppError(err, "sharing.CommandHandlers.CreateShare:GetFileInfo")
		}

		err = fileInfo.ValidateAccess(profileID)
		if err != nil {
			return nil, apperror.NewAppError(err, "sharing.CommandHandlers.CreateShare:File.ValidateAccess")
		}
	}

	if cmd.FolderID != nil {
		folderInfo, err := h.mediaRepository.GetFolderInfo(ctx, *cmd.FolderID)
		if err != nil {
			return nil, apperror.NewAppError(err, "sharing.CommandHandlers.CreateShare:GetFolderInfo")
		}

		err = folderInfo.ValidateAccess(profileID)
		if err != nil {
			return nil, apperror.NewAppError(err, "sharing.CommandHandlers.CreateShare:Folder.ValidateAccess")
		}
	}

	config, err := NewShareConfig(
		profileID,
		cmd.FileID,
		cmd.FolderID,
		cmd.Password,
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

	// Create share config
	config, err = repoTx.CreateShareConfig(ctx, config)
	if err != nil {
		return nil, apperror.NewAppError(err, "sharing.CommandHandlers.CreateShare:CreateShareConfig")
	}

	// Create share recipients
	for _, r := range cmd.Recipients {
		recipient, err := addShareRecipient(ctx, repoTx, profileID, config.ID, r)
		if err != nil {
			return nil, apperror.NewAppError(err, "sharing.CommandHandlers.CreateShare:addShareRecipient")
		}
		config.Recipients = append(config.Recipients, recipient)
	}

	err = tx.Commit()
	if err != nil {
		return nil, apperror.NewAppError(err, "sharing.CommandHandlers.CreateShare:Commit")
	}

	return config, nil
}

func (h *CommandHandlers) UpdateShareExpiry(ctx context.Context, cmd *UpdateShareExpiryCommand) error {
	profileID := common.GetProfileIDFromContext(ctx)
	err := h.repository.UpdateShareExpiry(ctx, profileID, cmd.ShareID, cmd.MaxDownloads, cmd.ExpiresAt)
	if err != nil {
		return apperror.NewAppError(err, "sharing.CommandHandlers.UpdateShareExpiry:UpdateShareExpiry")
	}

	return nil
}

func (h *CommandHandlers) UpdateSharePassword(ctx context.Context, cmd *UpdateSharePasswordCommand) error {
	profileID := common.GetProfileIDFromContext(ctx)
	pwdHash, err := utils.HashPassword(*cmd.Password)
	if err != nil {
		return apperror.NewAppError(err, "sharing.CommandHandlers.UpdateSharePassword:HashPassword")
	}

	err = h.repository.UpdateSharePassword(ctx, profileID, cmd.ShareID, &pwdHash)
	if err != nil {
		return apperror.NewAppError(err, "sharing.CommandHandlers.UpdateSharePassword:UpdateSharePassword")
	}

	return nil
}

func (h *CommandHandlers) DeleteShare(ctx context.Context, cmd *DeleteShareCommand) error {
	profileID := common.GetProfileIDFromContext(ctx)
	err := h.repository.DeleteShareConfig(ctx, profileID, cmd.ShareID)
	if err != nil {
		return apperror.NewAppError(err, "sharing.CommandHandlers.DeleteShare:DeleteShareConfig")
	}

	return nil
}

func (h *CommandHandlers) AddShareRecipient(ctx context.Context, cmd *AddShareRecipientCommand) (*ShareRecipient, error) {
	var tx *sql.Tx
	repoTx := h.repository
	if cmd.SaveAsContact {
		tx, err := repoTx.BeginTx(ctx)
		if err != nil {
			return nil, apperror.NewAppError(err, "sharing.CommandHandlers.AddShareRecipient:BeginTx")
		}
		defer tx.Rollback()
		repoTx = repoTx.WithTx(ctx, tx)
	}

	profileID := common.GetProfileIDFromContext(ctx)
	recipient, err := addShareRecipient(ctx, repoTx, profileID, cmd.ShareID, cmd.ShareRecipientInput)
	if err != nil {
		return nil, apperror.NewAppError(err, "sharing.CommandHandlers.AddShareRecipient:addShareRecipient")
	}

	if tx != nil {
		err = tx.Commit()
		if err != nil {
			return nil, apperror.NewAppError(err, "sharing.CommandHandlers.AddShareRecipient:Commit")
		}
	}
	return recipient, nil
}

// TODO: When the Block feature is implemented,
// add a check to see if the recipient has blocked the sender.
// Don't inform the sender if the recipient has blocked them.

func addShareRecipient(ctx context.Context, repoTx Repository, profileID, shareConfigID int64, recInput *ShareRecipientInput) (*ShareRecipient, error) {
	var nonContactEmail *string
	var contactID *int64
	if recInput.SaveAsContact {
		contact := NewContact(profileID, *recInput.Email, *recInput.Name)
		contact, err := repoTx.CreateContact(ctx, contact)
		if err != nil {
			return nil, apperror.NewAppError(err, "sharing.addShareRecipient:CreateContact")
		}
		contactID = &contact.ID
	} else {
		nonContactEmail = recInput.Email
	}

	recipient := NewShareRecipient(shareConfigID, contactID, recInput.ContactGroupID, nonContactEmail)

	_, err := repoTx.CreateShareRecipient(ctx, profileID, recipient)
	if err != nil {
		return nil, apperror.NewAppError(err, "sharing.addShareRecipient:CreateShareRecipient")
	}

	return recipient, nil
}

func (h *CommandHandlers) RemoveShareRecipient(ctx context.Context, cmd *RemoveShareRecipientCommand) error {
	profileID := common.GetProfileIDFromContext(ctx)
	err := h.repository.DeleteShareRecipient(ctx, profileID, cmd.ShareID, cmd.RecipientID)
	if err != nil {
		return apperror.NewAppError(err, "sharing.CommandHandlers.RemoveShareRecipient:DeleteShareRecipient")
	}

	return nil
}
