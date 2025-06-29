package sharing

import (
	"time"

	"skyvault/pkg/apperror"
	"skyvault/pkg/utils"
)

type ShareRecipient struct {
	ID             string
	ShareConfigID  string
	ContactID      *string // Only one of ContactID, ContactGroupID, or Email can be set
	ContactGroupID *string
	Email          *string // Email which we don't want to save as a contact
	DownloadsCount int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewShareRecipient(
	shareConfigID string,
	contactID *string,
	contactGroupID *string,
	nonContactEmail *string,
) (*ShareRecipient, error) {
	id, err := utils.ID()
	if err != nil {
		return nil, apperror.NewAppError(err, "sharing.NewShareRecipient:ID")
	}
	now := time.Now().UTC()
	return &ShareRecipient{
		ID:             id,
		ShareConfigID:  shareConfigID,
		ContactID:      contactID,
		ContactGroupID: contactGroupID,
		Email:          nonContactEmail,
		DownloadsCount: 0,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// func (r *ShareRecipient) IncrementDownloads() {
// 	r.DownloadsCount++
// 	r.UpdatedAt = time.Now().UTC()
// }
