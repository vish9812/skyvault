//lint:file-ignore ST1001 Using dot import to make SQL queries more readable
package repository

import (
	"context"
	"database/sql"
	"time"

	"skyvault/internal/domain/sharing"
	"skyvault/internal/infrastructure/internal/repository/internal/gen_jet/skyvault/public/model"
	. "skyvault/internal/infrastructure/internal/repository/internal/gen_jet/skyvault/public/table"
	"skyvault/pkg/apperror"
	"skyvault/pkg/paging"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/jinzhu/copier"
)

var _ sharing.Repository = (*SharingRepository)(nil)

type SharingRepository struct {
	repository *Repository
}

func NewSharingRepository(repo *Repository) *SharingRepository {
	return &SharingRepository{repository: repo}
}

func (r *SharingRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.repository.db.BeginTx(ctx, nil)
}

func (r *SharingRepository) WithTx(ctx context.Context, tx *sql.Tx) sharing.Repository {
	return &SharingRepository{repository: r.repository.withTx(ctx, tx)}
}

//--------------------------------
// Contacts
//--------------------------------

func (r *SharingRepository) CreateContact(ctx context.Context, contact *sharing.Contact) (*sharing.Contact, error) {
	dbModel := new(model.Contact)
	err := copier.Copy(dbModel, contact)
	if err != nil {
		return nil, apperror.NewAppError(err, "repository.CreateContact:copier.Copy")
	}

	stmt := Contact.INSERT(Contact.MutableColumns).MODEL(dbModel).RETURNING(Contact.AllColumns)

	return runInsert[model.Contact, sharing.Contact](ctx, stmt, r.repository.dbTx)
}

func (r *SharingRepository) GetContacts(ctx context.Context, ownerID int64, pagingOpt *paging.Options, searchTerm *string) (*paging.Page[*sharing.Contact], error) {
	whereCond := Contact.OwnerID.EQ(Int64(ownerID))

	if searchTerm != nil && *searchTerm != "" {
		searchPattern := "%" + *searchTerm + "%"
		whereCond = whereCond.AND(
			ILIKE(Contact.Name, String(searchPattern)).
				OR(ILIKE(Contact.Email, String(searchPattern))),
		)
	}

	orderBy := []OrderByClause{Contact.Name.ASC(), Contact.Email.ASC()}

	stmt := SELECT(Contact.AllColumns).
		FROM(Contact)

	cursorQuery := &cursorQuery{
		ID:        Contact.ID,
		Name:      Contact.Name,
		Updated:   Contact.UpdatedAt,
		where:     whereCond,
		orderBy:   orderBy,
		pagingOpt: pagingOpt,
	}

	return runSelectSlice[model.Contact, sharing.Contact](ctx, cursorQuery, stmt, r.repository.dbTx)
}

func (r *SharingRepository) UpdateContact(ctx context.Context, ownerID, contactID int64, email, name string) error {
	now := time.Now().UTC()
	contact := model.Contact{
		Email:     email,
		Name:      name,
		UpdatedAt: now,
	}

	stmt := Contact.UPDATE(Contact.Email, Contact.Name, Contact.UpdatedAt).
		MODEL(contact).
		WHERE(
			Contact.ID.EQ(Int64(contactID)).
				AND(Contact.OwnerID.EQ(Int64(ownerID))),
		)

	return runUpdateOrDelete(ctx, stmt, r.repository.dbTx)
}

func (r *SharingRepository) DeleteContact(ctx context.Context, ownerID, contactID int64) error {
	stmt := Contact.DELETE().
		WHERE(
			Contact.ID.EQ(Int64(contactID)).
				AND(Contact.OwnerID.EQ(Int64(ownerID))),
		)

	return runUpdateOrDelete(ctx, stmt, r.repository.dbTx)
}

//--------------------------------
// Contact Groups
//--------------------------------

func (r *SharingRepository) CreateContactGroup(ctx context.Context, group *sharing.ContactGroup) (*sharing.ContactGroup, error) {
	dbModel := new(model.ContactGroup)
	err := copier.Copy(dbModel, group)
	if err != nil {
		return nil, apperror.NewAppError(err, "repository.CreateContactGroup:copier.Copy")
	}

	stmt := ContactGroup.INSERT(ContactGroup.MutableColumns).MODEL(dbModel).RETURNING(ContactGroup.AllColumns)

	return runInsert[model.ContactGroup, sharing.ContactGroup](ctx, stmt, r.repository.dbTx)
}

func (r *SharingRepository) GetContactGroup(ctx context.Context, ownerID, groupID int64) (*sharing.ContactGroup, error) {
	stmt := SELECT(ContactGroup.AllColumns).
		FROM(ContactGroup).
		WHERE(ContactGroup.ID.EQ(Int64(groupID)).
			AND(ContactGroup.OwnerID.EQ(Int64(ownerID))),
		)

	return runSelect[model.ContactGroup, sharing.ContactGroup](ctx, stmt, r.repository.dbTx)
}

func (r *SharingRepository) GetContactGroups(ctx context.Context, ownerID int64, pagingOpt *paging.Options, searchTerm *string) (*paging.Page[*sharing.ContactGroup], error) {
	whereCond := ContactGroup.OwnerID.EQ(Int64(ownerID))

	if searchTerm != nil && *searchTerm != "" {
		searchPattern := "%" + *searchTerm + "%"
		whereCond = whereCond.AND(ContactGroup.Name.LIKE(String(searchPattern)))
	}

	orderBy := []OrderByClause{ContactGroup.Name.ASC()}

	stmt := SELECT(ContactGroup.AllColumns).
		FROM(ContactGroup)

	cursorQuery := &cursorQuery{
		ID:        ContactGroup.ID,
		Name:      ContactGroup.Name,
		Updated:   ContactGroup.UpdatedAt,
		where:     whereCond,
		orderBy:   orderBy,
		pagingOpt: pagingOpt,
	}

	return runSelectSlice[model.ContactGroup, sharing.ContactGroup](ctx, cursorQuery, stmt, r.repository.dbTx)
}

func (r *SharingRepository) GetContactGroupMembers(ctx context.Context, ownerID int64, pagingOpt *paging.Options, groupID int64) (*paging.Page[*sharing.Contact], error) {
	whereCond := ContactGroupMember.GroupID.EQ(Int64(groupID)).
		AND(Contact.OwnerID.EQ(Int64(ownerID)))

	orderBy := []OrderByClause{Contact.Name.ASC(), Contact.Email.ASC()}

	stmt := SELECT(Contact.AllColumns).
		FROM(Contact.
			INNER_JOIN(ContactGroupMember, ContactGroupMember.ContactID.EQ(Contact.ID)))

	cursorQuery := &cursorQuery{
		ID:        Contact.ID,
		Name:      Contact.Name,
		Updated:   Contact.UpdatedAt,
		where:     whereCond,
		orderBy:   orderBy,
		pagingOpt: pagingOpt,
	}

	return runSelectSlice[model.Contact, sharing.Contact](ctx, cursorQuery, stmt, r.repository.dbTx)
}

func (r *SharingRepository) UpdateContactGroup(ctx context.Context, ownerID, groupID int64, name string) error {
	now := time.Now().UTC()
	group := model.ContactGroup{
		Name:      name,
		UpdatedAt: now,
	}

	stmt := ContactGroup.UPDATE(ContactGroup.Name, ContactGroup.UpdatedAt).
		MODEL(group).
		WHERE(
			ContactGroup.ID.EQ(Int64(groupID)).
				AND(ContactGroup.OwnerID.EQ(Int64(ownerID))),
		)

	return runUpdateOrDelete(ctx, stmt, r.repository.dbTx)
}

func (r *SharingRepository) DeleteContactGroup(ctx context.Context, ownerID, groupID int64) error {
	stmt := ContactGroup.DELETE().
		WHERE(
			ContactGroup.ID.EQ(Int64(groupID)).
				AND(ContactGroup.OwnerID.EQ(Int64(ownerID))),
		)

	return runUpdateOrDelete(ctx, stmt, r.repository.dbTx)
}

func (r *SharingRepository) AddContactToGroup(ctx context.Context, ownerID, groupID, contactID int64) error {
	// First verify that both the group and contact belong to the owner
	groupStmt := SELECT(ContactGroup.ID).
		FROM(ContactGroup).
		WHERE(
			ContactGroup.ID.EQ(Int64(groupID)).
				AND(ContactGroup.OwnerID.EQ(Int64(ownerID))),
		)

	var groupExists bool
	err := groupStmt.QueryContext(ctx, r.repository.dbTx, &groupExists)
	if err != nil {
		return apperror.NewAppError(err, "repository.AddContactToGroup:VerifyGroup")
	}
	if !groupExists {
		return apperror.NewAppError(apperror.ErrCommonNoData, "repository.AddContactToGroup:GroupNotFound")
	}

	contactStmt := SELECT(Contact.ID).
		FROM(Contact).
		WHERE(
			Contact.ID.EQ(Int64(contactID)).
				AND(Contact.OwnerID.EQ(Int64(ownerID))),
		)

	var contactExists bool
	err = contactStmt.QueryContext(ctx, r.repository.dbTx, &contactExists)
	if err != nil {
		return apperror.NewAppError(err, "repository.AddContactToGroup:VerifyContact")
	}
	if !contactExists {
		return apperror.NewAppError(apperror.ErrCommonNoData, "repository.AddContactToGroup:ContactNotFound")
	}

	// Now add the contact to the group
	now := time.Now().UTC()
	member := model.ContactGroupMember{
		GroupID:   groupID,
		ContactID: contactID,
		CreatedAt: now,
	}

	stmt := ContactGroupMember.INSERT(
		ContactGroupMember.GroupID,
		ContactGroupMember.ContactID,
		ContactGroupMember.CreatedAt,
	).MODEL(member)

	return runInsertNoReturn(ctx, stmt, r.repository.dbTx)
}

func (r *SharingRepository) RemoveContactFromGroup(ctx context.Context, ownerID, groupID, contactID int64) error {
	// First verify that the group belongs to the owner
	groupStmt := SELECT(ContactGroup.ID).
		FROM(ContactGroup).
		WHERE(
			ContactGroup.ID.EQ(Int64(groupID)).
				AND(ContactGroup.OwnerID.EQ(Int64(ownerID))),
		)

	var groupExists bool
	err := groupStmt.QueryContext(ctx, r.repository.dbTx, &groupExists)
	if err != nil {
		return apperror.NewAppError(err, "repository.RemoveContactFromGroup:VerifyGroup")
	}
	if !groupExists {
		return apperror.NewAppError(apperror.ErrCommonNoData, "repository.RemoveContactFromGroup:GroupNotFound")
	}

	// Now remove the contact from the group
	stmt := ContactGroupMember.DELETE().
		WHERE(
			ContactGroupMember.GroupID.EQ(Int64(groupID)).
				AND(ContactGroupMember.ContactID.EQ(Int64(contactID))),
		)

	return runUpdateOrDelete(ctx, stmt, r.repository.dbTx)
}

//--------------------------------
// Share Configs
//--------------------------------

func (r *SharingRepository) CreateShareConfig(ctx context.Context, config *sharing.ShareConfig) (*sharing.ShareConfig, error) {
	dbModel := new(model.ShareConfig)
	err := copier.Copy(dbModel, config)
	if err != nil {
		return nil, apperror.NewAppError(err, "repository.CreateShareConfig:copier.Copy")
	}

	stmt := ShareConfig.INSERT(
		ShareConfig.CustomID,
		ShareConfig.OwnerID,
		ShareConfig.FileID,
		ShareConfig.FolderID,
		ShareConfig.PasswordHash,
		ShareConfig.MaxDownloads,
		ShareConfig.ExpiresAt,
		ShareConfig.CreatedAt,
		ShareConfig.UpdatedAt,
	).MODEL(dbModel).RETURNING(ShareConfig.AllColumns)

	return runInsert[model.ShareConfig, sharing.ShareConfig](ctx, stmt, r.repository.dbTx)
}

func (r *SharingRepository) GetShareConfig(ctx context.Context, ownerID, shareID int64) (*sharing.ShareConfig, error) {
	stmt := SELECT(ShareConfig.AllColumns).
		FROM(ShareConfig).
		WHERE(
			ShareConfig.ID.EQ(Int64(shareID)).
				AND(ShareConfig.OwnerID.EQ(Int64(ownerID))),
		)

	config, err := runSelect[model.ShareConfig, sharing.ShareConfig](ctx, stmt, r.repository.dbTx)
	if err != nil {
		return nil, err
	}

	// Get recipients
	recipients, err := r.getShareRecipients(ctx, config.ID)
	if err != nil {
		return nil, err
	}

	config.Recipients = recipients
	return config, nil
}

func (r *SharingRepository) GetShareConfigByCustomID(ctx context.Context, ownerID int64, customID string) (*sharing.ShareConfig, error) {
	stmt := SELECT(ShareConfig.AllColumns).
		FROM(ShareConfig).
		WHERE(
			ShareConfig.CustomID.EQ(String(customID)).
				AND(ShareConfig.OwnerID.EQ(Int64(ownerID))),
		)

	config, err := runSelect[model.ShareConfig, sharing.ShareConfig](ctx, stmt, r.repository.dbTx)
	if err != nil {
		return nil, err
	}

	// Get recipients
	recipients, err := r.getShareRecipients(ctx, config.ID)
	if err != nil {
		return nil, err
	}

	config.Recipients = recipients
	return config, nil
}

func (r *SharingRepository) getShareRecipients(ctx context.Context, shareConfigID int64) ([]*sharing.ShareRecipient, error) {
	stmt := SELECT(ShareRecipient.AllColumns).
		FROM(ShareRecipient).
		WHERE(ShareRecipient.ShareConfigID.EQ(Int64(shareConfigID)))

	var dbRecipients []*model.ShareRecipient
	err := stmt.QueryContext(ctx, r.repository.dbTx, &dbRecipients)
	if err != nil {
		return nil, apperror.NewAppError(err, "repository.getShareRecipients:QueryContext")
	}

	recipients := make([]*sharing.ShareRecipient, len(dbRecipients))
	for i, dbRecipient := range dbRecipients {
		recipient := new(sharing.ShareRecipient)
		err = copier.Copy(recipient, dbRecipient)
		if err != nil {
			return nil, apperror.NewAppError(err, "repository.getShareRecipients:copier.Copy")
		}
		recipients[i] = recipient
	}

	return recipients, nil
}

func (r *SharingRepository) UpdateShareExpiry(ctx context.Context, ownerID, shareID int64, maxDownloads *int, expiresAt *time.Time) error {
	now := time.Now().UTC()
	config := model.ShareConfig{
		MaxDownloads: maxDownloads,
		ExpiresAt:    expiresAt,
		UpdatedAt:    now,
	}

	stmt := ShareConfig.UPDATE(ShareConfig.MaxDownloads, ShareConfig.ExpiresAt, ShareConfig.UpdatedAt).
		MODEL(config).
		WHERE(
			ShareConfig.ID.EQ(Int64(shareID)).
				AND(ShareConfig.OwnerID.EQ(Int64(ownerID))),
		)

	return runUpdateOrDelete(ctx, stmt, r.repository.dbTx)
}

func (r *SharingRepository) UpdateSharePassword(ctx context.Context, ownerID, shareID int64, passwordHash *string) error {
	now := time.Now().UTC()
	config := model.ShareConfig{
		PasswordHash: passwordHash,
		UpdatedAt:    now,
	}

	stmt := ShareConfig.UPDATE(ShareConfig.PasswordHash, ShareConfig.UpdatedAt).
		MODEL(config).
		WHERE(
			ShareConfig.ID.EQ(Int64(shareID)).
				AND(ShareConfig.OwnerID.EQ(Int64(ownerID))),
		)

	return runUpdateOrDelete(ctx, stmt, r.repository.dbTx)
}

func (r *SharingRepository) DeleteShareConfig(ctx context.Context, ownerID, shareID int64) error {
	stmt := ShareConfig.DELETE().
		WHERE(
			ShareConfig.ID.EQ(Int64(shareID)).
				AND(ShareConfig.OwnerID.EQ(Int64(ownerID))),
		)

	return runUpdateOrDelete(ctx, stmt, r.repository.dbTx)
}

//--------------------------------
// Share Recipients
//--------------------------------

func (r *SharingRepository) CreateShareRecipient(ctx context.Context, ownerID int64, recipient *sharing.ShareRecipient) (*sharing.ShareRecipient, error) {
	// First verify that the share config belongs to the owner
	shareStmt := SELECT(ShareConfig.ID).
		FROM(ShareConfig).
		WHERE(
			ShareConfig.ID.EQ(Int64(recipient.ShareConfigID)).
				AND(ShareConfig.OwnerID.EQ(Int64(ownerID))),
		)

	var shareExists bool
	err := shareStmt.QueryContext(ctx, r.repository.dbTx, &shareExists)
	if err != nil {
		return nil, apperror.NewAppError(err, "repository.CreateShareRecipient:VerifyShare")
	}
	if !shareExists {
		return nil, apperror.NewAppError(apperror.ErrCommonNoData, "repository.CreateShareRecipient:ShareNotFound")
	}

	dbModel := new(model.ShareRecipient)
	err = copier.Copy(dbModel, recipient)
	if err != nil {
		return nil, apperror.NewAppError(err, "repository.CreateShareRecipient:copier.Copy")
	}

	stmt := ShareRecipient.INSERT(
		ShareRecipient.ShareConfigID,
		ShareRecipient.ContactID,
		ShareRecipient.GroupID,
		ShareRecipient.Email,
		ShareRecipient.DownloadsCount,
		ShareRecipient.CreatedAt,
		ShareRecipient.UpdatedAt,
	).MODEL(dbModel).RETURNING(ShareRecipient.AllColumns)

	return runInsert[model.ShareRecipient, sharing.ShareRecipient](ctx, stmt, r.repository.dbTx)
}

func (r *SharingRepository) DeleteShareRecipient(ctx context.Context, ownerID, recipientID int64) error {
	// First verify that the recipient belongs to a share config owned by the owner
	verifyStmt := SELECT(ShareConfig.ID).
		FROM(ShareConfig.
			INNER_JOIN(ShareRecipient, ShareRecipient.ShareConfigID.EQ(ShareConfig.ID))).
		WHERE(
			ShareRecipient.ID.EQ(Int64(recipientID)).
				AND(ShareConfig.OwnerID.EQ(Int64(ownerID))),
		)

	var shareExists bool
	err := verifyStmt.QueryContext(ctx, r.repository.dbTx, &shareExists)
	if err != nil {
		return apperror.NewAppError(err, "repository.DeleteShareRecipient:VerifyOwnership")
	}
	if !shareExists {
		return apperror.NewAppError(apperror.ErrCommonNoData, "repository.DeleteShareRecipient:RecipientNotFound")
	}

	stmt := ShareRecipient.DELETE().
		WHERE(ShareRecipient.ID.EQ(Int64(recipientID)))

	return runUpdateOrDelete(ctx, stmt, r.repository.dbTx)
}

func (r *SharingRepository) GetShareRecipientByEmail(ctx context.Context, shareID int64, email string) (*sharing.ShareRecipient, error) {
	stmt := SELECT(ShareRecipient.AllColumns).
		FROM(ShareRecipient).
		WHERE(
			ShareRecipient.ShareConfigID.EQ(Int64(shareID)).
				AND(ShareRecipient.Email.EQ(String(email))),
		)

	return runSelect[model.ShareRecipient, sharing.ShareRecipient](ctx, stmt, r.repository.dbTx)
}
