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
	"skyvault/pkg/utils"

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

	stmt := Contact.INSERT(Contact.AllColumns).MODEL(dbModel).RETURNING(Contact.AllColumns)

	return runInsert[model.Contact, sharing.Contact](ctx, stmt, r.repository.dbTx)
}

func (r *SharingRepository) GetContacts(ctx context.Context, ownerID string, pagingOpt *paging.Options, searchTerm *string) (*paging.Page[*sharing.Contact], error) {
	whereCond := Contact.OwnerID.EQ(UUID(UUIDStr(ownerID)))

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

func (r *SharingRepository) UpdateContact(ctx context.Context, ownerID, contactID, email, name string) error {
	now := time.Now().UTC()
	contact := model.Contact{
		Email:     email,
		Name:      name,
		UpdatedAt: now,
	}

	stmt := Contact.UPDATE(Contact.Email, Contact.Name, Contact.UpdatedAt).
		MODEL(contact).
		WHERE(
			Contact.ID.EQ(UUID(UUIDStr(contactID))).
				AND(Contact.OwnerID.EQ(UUID(UUIDStr(ownerID)))),
		)

	return runUpdateOrDelete(ctx, stmt, r.repository.dbTx)
}

func (r *SharingRepository) DeleteContact(ctx context.Context, ownerID, contactID string) error {
	stmt := Contact.DELETE().
		WHERE(
			Contact.ID.EQ(UUID(UUIDStr(contactID))).
				AND(Contact.OwnerID.EQ(UUID(UUIDStr(ownerID)))),
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

	stmt := ContactGroup.INSERT(ContactGroup.AllColumns).MODEL(dbModel).RETURNING(ContactGroup.AllColumns)

	return runInsert[model.ContactGroup, sharing.ContactGroup](ctx, stmt, r.repository.dbTx)
}

func (r *SharingRepository) GetContactGroup(ctx context.Context, ownerID, groupID string) (*sharing.ContactGroup, error) {
	stmt := SELECT(ContactGroup.AllColumns).
		FROM(ContactGroup).
		WHERE(ContactGroup.ID.EQ(UUID(UUIDStr(groupID))).
			AND(ContactGroup.OwnerID.EQ(UUID(UUIDStr(ownerID)))),
		)

	return runSelect[model.ContactGroup, sharing.ContactGroup](ctx, stmt, r.repository.dbTx)
}

func (r *SharingRepository) GetContactGroups(ctx context.Context, ownerID string, pagingOpt *paging.Options, searchTerm *string) (*paging.Page[*sharing.ContactGroup], error) {
	whereCond := ContactGroup.OwnerID.EQ(UUID(UUIDStr(ownerID)))

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

func (r *SharingRepository) GetContactGroupMembers(ctx context.Context, ownerID string, pagingOpt *paging.Options, groupID string) (*paging.Page[*sharing.Contact], error) {
	whereCond := ContactGroupMember.GroupID.EQ(UUID(UUIDStr(groupID))).
		AND(Contact.OwnerID.EQ(UUID(UUIDStr(ownerID))))

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

func (r *SharingRepository) UpdateContactGroup(ctx context.Context, ownerID, groupID string, name string) error {
	now := time.Now().UTC()
	group := model.ContactGroup{
		Name:      name,
		UpdatedAt: now,
	}

	stmt := ContactGroup.UPDATE(ContactGroup.Name, ContactGroup.UpdatedAt).
		MODEL(group).
		WHERE(
			ContactGroup.ID.EQ(UUID(UUIDStr(groupID))).
				AND(ContactGroup.OwnerID.EQ(UUID(UUIDStr(ownerID)))),
		)

	return runUpdateOrDelete(ctx, stmt, r.repository.dbTx)
}

func (r *SharingRepository) DeleteContactGroup(ctx context.Context, ownerID, groupID string) error {
	stmt := ContactGroup.DELETE().
		WHERE(
			ContactGroup.ID.EQ(UUID(UUIDStr(groupID))).
				AND(ContactGroup.OwnerID.EQ(UUID(UUIDStr(ownerID)))),
		)

	return runUpdateOrDelete(ctx, stmt, r.repository.dbTx)
}

func (r *SharingRepository) AddContactToGroup(ctx context.Context, ownerID, groupID, contactID string) error {
	validContact := SELECT(Contact.ID).
		FROM(Contact).
		WHERE(
			Contact.ID.EQ(UUID(UUIDStr(contactID))).
				AND(Contact.OwnerID.EQ(UUID(UUIDStr(ownerID)))),
		).AsTable("valid_contact")

	validContactID := Contact.ID.From(validContact)

	validGroup := SELECT(ContactGroup.ID).
		FROM(ContactGroup).
		WHERE(
			ContactGroup.ID.EQ(UUID(UUIDStr(groupID))).
				AND(ContactGroup.OwnerID.EQ(UUID(UUIDStr(ownerID)))),
		).AsTable("valid_group")

	validGroupID := ContactGroup.ID.From(validGroup)

	id, err := utils.ID()
	if err != nil {
		return apperror.NewAppError(err, "repository.AddContactToGroup:utils.ID")
	}

	stmt := ContactGroupMember.INSERT(
		ContactGroupMember.ID,
		ContactGroupMember.GroupID,
		ContactGroupMember.ContactID,
		ContactGroupMember.CreatedAt,
	).
		QUERY(SELECT(UUID(UUIDStr(id)), validGroupID, validContactID, TimestampT(time.Now().UTC())).
			FROM(validGroup, validContact),
		)

	return runInsertNoReturn(ctx, stmt, r.repository.dbTx)
}

func (r *SharingRepository) RemoveContactFromGroup(ctx context.Context, ownerID, groupID, contactID string) error {
	validContact := SELECT(Contact.ID).
		FROM(Contact).
		WHERE(
			Contact.ID.EQ(UUID(UUIDStr(contactID))).
				AND(Contact.OwnerID.EQ(UUID(UUIDStr(ownerID)))),
		).AsTable("valid_contact")

	validContactID := Contact.ID.From(validContact)

	validGroup := SELECT(ContactGroup.ID).
		FROM(ContactGroup).
		WHERE(
			ContactGroup.ID.EQ(UUID(UUIDStr(groupID))).
				AND(ContactGroup.OwnerID.EQ(UUID(UUIDStr(ownerID)))),
		).AsTable("valid_group")

	validGroupID := ContactGroup.ID.From(validGroup)

	stmt := ContactGroupMember.DELETE().
		USING(validGroup, validContact).
		WHERE(
			ContactGroupMember.GroupID.EQ(validGroupID).
				AND(ContactGroupMember.ContactID.EQ(validContactID)),
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

	stmt := ShareConfig.INSERT(ShareConfig.AllColumns).MODEL(dbModel).RETURNING(ShareConfig.AllColumns)

	return runInsert[model.ShareConfig, sharing.ShareConfig](ctx, stmt, r.repository.dbTx)
}

func (r *SharingRepository) GetShareConfig(ctx context.Context, ownerID, shareID string) (*sharing.ShareConfig, error) {
	stmt := SELECT(ShareConfig.AllColumns).
		FROM(ShareConfig).
		WHERE(
			ShareConfig.ID.EQ(UUID(UUIDStr(shareID))).
				AND(ShareConfig.OwnerID.EQ(UUID(UUIDStr(ownerID)))),
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

func (r *SharingRepository) getShareRecipients(ctx context.Context, shareConfigID string) ([]*sharing.ShareRecipient, error) {
	stmt := SELECT(ShareRecipient.AllColumns).
		FROM(ShareRecipient).
		WHERE(ShareRecipient.ShareConfigID.EQ(UUID(UUIDStr(shareConfigID))))

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

func (r *SharingRepository) UpdateShareExpiry(ctx context.Context, ownerID, shareID string, maxDownloads *int64, expiresAt *time.Time) error {
	now := time.Now().UTC()
	config := model.ShareConfig{
		MaxDownloads: maxDownloads,
		ExpiresAt:    expiresAt,
		UpdatedAt:    now,
	}

	stmt := ShareConfig.UPDATE(ShareConfig.MaxDownloads, ShareConfig.ExpiresAt, ShareConfig.UpdatedAt).
		MODEL(config).
		WHERE(
			ShareConfig.ID.EQ(UUID(UUIDStr(shareID))).
				AND(ShareConfig.OwnerID.EQ(UUID(UUIDStr(ownerID)))),
		)

	return runUpdateOrDelete(ctx, stmt, r.repository.dbTx)
}

func (r *SharingRepository) UpdateSharePassword(ctx context.Context, ownerID, shareID string, passwordHash *string) error {
	now := time.Now().UTC()
	config := model.ShareConfig{
		PasswordHash: passwordHash,
		UpdatedAt:    now,
	}

	stmt := ShareConfig.UPDATE(ShareConfig.PasswordHash, ShareConfig.UpdatedAt).
		MODEL(config).
		WHERE(
			ShareConfig.ID.EQ(UUID(UUIDStr(shareID))).
				AND(ShareConfig.OwnerID.EQ(UUID(UUIDStr(ownerID)))),
		)

	return runUpdateOrDelete(ctx, stmt, r.repository.dbTx)
}

func (r *SharingRepository) DeleteShareConfig(ctx context.Context, ownerID, shareID string) error {
	stmt := ShareConfig.DELETE().
		WHERE(
			ShareConfig.ID.EQ(UUID(UUIDStr(shareID))).
				AND(ShareConfig.OwnerID.EQ(UUID(UUIDStr(ownerID)))),
		)

	return runUpdateOrDelete(ctx, stmt, r.repository.dbTx)
}

//--------------------------------
// Share Recipients
//--------------------------------

func (r *SharingRepository) CreateShareRecipient(ctx context.Context, ownerID string, recipient *sharing.ShareRecipient) (*sharing.ShareRecipient, error) {
	// First verify that the share config belongs to the owner
	shareStmt := SELECT(ShareConfig.ID).
		FROM(ShareConfig).
		WHERE(
			ShareConfig.ID.EQ(UUID(UUIDStr(recipient.ShareConfigID))).
				AND(ShareConfig.OwnerID.EQ(UUID(UUIDStr(ownerID)))),
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

	stmt := ShareRecipient.INSERT(ShareRecipient.AllColumns).MODEL(dbModel).RETURNING(ShareRecipient.AllColumns)

	return runInsert[model.ShareRecipient, sharing.ShareRecipient](ctx, stmt, r.repository.dbTx)
}

func (r *SharingRepository) DeleteShareRecipient(ctx context.Context, ownerID, shareID, recipientID string) error {
	// First verify that the recipient belongs to a share config owned by the owner
	shareSubTable := SELECT(ShareConfig.ID).
		FROM(ShareConfig).
		WHERE(ShareConfig.ID.EQ(UUID(UUIDStr(shareID))).
			AND(ShareConfig.OwnerID.EQ(UUID(UUIDStr(ownerID)))),
		).AsTable("share_sub_table")

	shareConfigID := ShareConfig.ID.From(shareSubTable)

	stmt := ShareRecipient.DELETE().
		WHERE(ShareRecipient.ID.EQ(UUID(UUIDStr(recipientID))).
			AND(ShareRecipient.ShareConfigID.EQ(shareConfigID)),
		)

	return runUpdateOrDelete(ctx, stmt, r.repository.dbTx)
}

func (r *SharingRepository) GetShareRecipientByEmail(ctx context.Context, shareID string, email string) (*sharing.ShareRecipient, error) {
	stmt := SELECT(ShareRecipient.AllColumns).
		FROM(ShareRecipient).
		WHERE(
			ShareRecipient.ShareConfigID.EQ(UUID(UUIDStr(shareID))).
				AND(ShareRecipient.Email.EQ(String(email))),
		)

	return runSelect[model.ShareRecipient, sharing.ShareRecipient](ctx, stmt, r.repository.dbTx)
}
