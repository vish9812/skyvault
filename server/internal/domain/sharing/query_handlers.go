package sharing

import (
	"context"
	"skyvault/pkg/apperror"
	"skyvault/pkg/paging"
	"time"
)

var _ Queries = (*QueryHandlers)(nil)

type QueryHandlers struct {
	repository Repository
}

func NewQueryHandlers(repository Repository) *QueryHandlers {
	return &QueryHandlers{repository: repository}
}

//--------------------------------
// Contacts
//--------------------------------

func (h *QueryHandlers) GetContact(ctx context.Context, query *GetContactQuery) (*Contact, error) {
	contact, err := h.repository.GetContact(ctx, query.OwnerID, query.ContactID)
	if err != nil {
		return nil, apperror.NewAppError(err, "sharing.QueryHandlers.GetContact:GetContact")
	}

	return contact, nil
}

func (h *QueryHandlers) GetContacts(ctx context.Context, query *GetContactsQuery) (*paging.Page[*Contact], error) {
	contacts, err := h.repository.GetContacts(ctx, query.PagingOpt, query.OwnerID, query.SearchTerm)
	if err != nil {
		return nil, apperror.NewAppError(err, "sharing.QueryHandlers.GetContacts:GetContacts")
	}

	return contacts, nil
}

//--------------------------------
// Contact Groups
//--------------------------------

func (h *QueryHandlers) GetContactGroup(ctx context.Context, query *GetContactGroupQuery) (*ContactGroup, error) {
	group, err := h.repository.GetContactGroup(ctx, query.OwnerID, query.GroupID)
	if err != nil {
		return nil, apperror.NewAppError(err, "sharing.QueryHandlers.GetContactGroup:GetContactGroup")
	}

	return group, nil
}

func (h *QueryHandlers) GetContactGroups(ctx context.Context, query *GetContactGroupsQuery) (*paging.Page[*ContactGroup], error) {
	groups, err := h.repository.GetContactGroups(ctx, query.PagingOpt, query.OwnerID, query.SearchTerm)
	if err != nil {
		return nil, apperror.NewAppError(err, "sharing.QueryHandlers.GetContactGroups:GetContactGroups")
	}

	return groups, nil
}

func (h *QueryHandlers) GetContactGroupMembers(ctx context.Context, query *GetContactGroupMembersQuery) (*paging.Page[*Contact], error) {
	members, err := h.repository.GetContactGroupMembers(ctx, query.PagingOpt, query.OwnerID, query.GroupID)
	if err != nil {
		return nil, apperror.NewAppError(err, "sharing.QueryHandlers.GetContactGroupMembers:GetContactGroupMembers")
	}

	return members, nil
}

//--------------------------------
// Sharing
//--------------------------------

func (h *QueryHandlers) GetShareConfig(ctx context.Context, query *GetShareConfigQuery) (*ShareConfig, error) {
	config, err := h.repository.GetShareConfig(ctx, query.ShareID)
	if err != nil {
		return nil, apperror.NewAppError(err, "sharing.QueryHandlers.GetShareConfig:GetShareConfig")
	}

	return config, nil
}

func (h *QueryHandlers) ValidateShareAccess(ctx context.Context, query *ValidateShareAccessQuery) error {
	config, err := h.repository.GetShareConfig(ctx, query.ShareID)
	if err != nil {
		return apperror.NewAppError(err, "sharing.QueryHandlers.ValidateShareAccess:GetShareConfig")
	}

	// Check expiry
	if config.ExpiresAt != nil && config.ExpiresAt.Before(time.Now().UTC()) {
		return apperror.NewAppError(apperror.ErrSharingExpired, "sharing.QueryHandlers.ValidateShareAccess:Expired")
	}

	// Check password if set
	if config.PasswordHash != nil {
		if query.Password == nil {
			return apperror.NewAppError(apperror.ErrSharingInvalidPassword, "sharing.QueryHandlers.ValidateShareAccess:PasswordRequired")
		}

		// TODO: Implement password validation
		// if !validatePassword(*query.Password, *config.PasswordHash) {
		// 	return apperror.NewAppError(apperror.ErrSharingInvalidPassword, "sharing.QueryHandlers.ValidateShareAccess:InvalidPassword")
		// }
	}

	// Get recipient
	recipient, err := h.repository.GetShareRecipient(ctx, query.ShareID, query.RecipientID)
	if err != nil {
		return apperror.NewAppError(err, "sharing.QueryHandlers.ValidateShareAccess:GetShareRecipient")
	}

	// Check max downloads
	if config.MaxDownloads != nil && recipient.DownloadsCount >= *config.MaxDownloads {
		return apperror.NewAppError(apperror.ErrSharingMaxDownloadsReached, "sharing.QueryHandlers.ValidateShareAccess:MaxDownloadsReached")
	}

	// Record access and increment downloads
	err = h.repository.RecordShareAccess(ctx, query.ShareID, "TODO: Get IP from context")
	if err != nil {
		return apperror.NewAppError(err, "sharing.QueryHandlers.ValidateShareAccess:RecordShareAccess")
	}

	recipient.IncrementDownloads()
	err = h.repository.UpdateShareRecipient(ctx, recipient)
	if err != nil {
		return apperror.NewAppError(err, "sharing.QueryHandlers.ValidateShareAccess:UpdateShareRecipient")
	}

	return nil
}
