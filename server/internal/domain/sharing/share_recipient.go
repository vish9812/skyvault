package sharing

import (
	"time"
)

type ShareRecipient struct {
	ID             int64
	ShareConfigID  int64
	ContactID      *int64 // Only one of ContactID, ContactGroupID, or Email can be set
	ContactGroupID *int64
	Email          *string // Email which we don't want to save as a contact
	DownloadsCount int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewShareRecipient(
	shareConfigID int64,
	contactID *int64,
	contactGroupID *int64,
	nonContactEmail *string,
) *ShareRecipient {
	now := time.Now().UTC()
	return &ShareRecipient{
		ShareConfigID:  shareConfigID,
		ContactID:      contactID,
		ContactGroupID: contactGroupID,
		Email:          nonContactEmail,
		DownloadsCount: 0,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// func (r *ShareRecipient) IncrementDownloads() {
// 	r.DownloadsCount++
// 	r.UpdatedAt = time.Now().UTC()
// }
