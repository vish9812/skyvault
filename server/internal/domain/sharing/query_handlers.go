package sharing

import (
	"context"
	"skyvault/pkg/apperror"
	"skyvault/pkg/common"
	"skyvault/pkg/paging"
	"skyvault/pkg/utils"
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

func (h *QueryHandlers) GetContacts(ctx context.Context, query *GetContactsQuery) (*paging.Page[*Contact], error) {
	profileID := common.GetProfileIDFromContext(ctx)
	contacts, err := h.repository.GetContacts(ctx, profileID, query.PagingOpt, query.SearchTerm)
	if err != nil {
		return nil, apperror.NewAppError(err, "sharing.QueryHandlers.GetContacts:GetContacts")
	}

	return contacts, nil
}

//--------------------------------
// Contact Groups
//--------------------------------

func (h *QueryHandlers) GetContactGroups(ctx context.Context, query *GetContactGroupsQuery) (*paging.Page[*ContactGroup], error) {
	profileID := common.GetProfileIDFromContext(ctx)
	groups, err := h.repository.GetContactGroups(ctx, profileID, query.PagingOpt, query.SearchTerm)
	if err != nil {
		return nil, apperror.NewAppError(err, "sharing.QueryHandlers.GetContactGroups:GetContactGroups")
	}

	return groups, nil
}

func (h *QueryHandlers) GetContactGroupMembers(ctx context.Context, query *GetContactGroupMembersQuery) (*paging.Page[*Contact], error) {
	profileID := common.GetProfileIDFromContext(ctx)
	members, err := h.repository.GetContactGroupMembers(ctx, profileID, query.PagingOpt, query.GroupID)
	if err != nil {
		return nil, apperror.NewAppError(err, "sharing.QueryHandlers.GetContactGroupMembers:GetContactGroupMembers")
	}

	return members, nil
}

//--------------------------------
// Sharing
//--------------------------------

func (h *QueryHandlers) GetShareConfig(ctx context.Context, query *GetShareConfigQuery) (*ShareConfig, error) {
	profileID := common.GetProfileIDFromContext(ctx)
	config, err := h.repository.GetShareConfig(ctx, profileID, query.ShareID)
	if err != nil {
		return nil, apperror.NewAppError(err, "sharing.QueryHandlers.GetShareConfig:GetShareConfig")
	}

	return config, nil
}

func (h *QueryHandlers) ValidateShareAccess(ctx context.Context, query *ValidateShareAccessQuery) error {
	profileID := common.GetProfileIDFromContext(ctx)
	config, err := h.repository.GetShareConfigByCustomID(ctx, profileID, query.CustomID)
	if err != nil {
		return apperror.NewAppError(err, "sharing.QueryHandlers.ValidateShareAccess:GetShareConfig")
	}

	// validate expiry
	err = config.ValidateExpiry()
	if err != nil {
		return apperror.NewAppError(err, "sharing.QueryHandlers.ValidateShareAccess:ValidateExpiry")
	}

	// validate password
	if config.PasswordHash != nil {
		if query.Password == nil {
			return apperror.NewAppError(apperror.ErrSharingInvalidCredentials, "sharing.QueryHandlers.ValidateShareAccess:PasswordRequired")
		}

		ok, err := utils.SamePassword(*config.PasswordHash, *query.Password)
		if err != nil {
			return apperror.NewAppError(err, "sharing.QueryHandlers.ValidateShareAccess:SamePassword")
		}
		if !ok {
			return apperror.NewAppError(apperror.ErrSharingInvalidCredentials, "sharing.QueryHandlers.ValidateShareAccess:InvalidPassword")
		}
	}

	if query.Email == nil {
		return nil
	}

	// validate recipient
	_, err = h.repository.GetShareRecipientByEmail(ctx, config.ID, *query.Email)
	if err != nil {
		return apperror.NewAppError(err, "sharing.QueryHandlers.ValidateShareAccess:GetShareRecipientByEmail")
	}

	return nil
}
