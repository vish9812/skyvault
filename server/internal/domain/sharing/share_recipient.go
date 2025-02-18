package sharing

import (
	"skyvault/pkg/apperror"
	"strings"
	"time"
)

type RecipientType string

const (
	RecipientTypeEmail RecipientType = "email"
	RecipientTypeGroup RecipientType = "group"
)

type ShareRecipient struct {
	ID             int64
	ShareConfigID  int64
	RecipientType  RecipientType
	RecipientID    *int64 // null for direct email shares, references contact_group(id) for group shares
	Email          string
	DownloadsCount int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewShareRecipient(
	shareConfigID int64,
	recipientType RecipientType,
	recipientID *int64,
	email string,
) (*ShareRecipient, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "sharing.NewShareRecipient:Email")
	}

	switch recipientType {
	case RecipientTypeEmail:
		if recipientID != nil {
			return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "sharing.NewShareRecipient:RecipientID")
		}
	case RecipientTypeGroup:
		if recipientID == nil {
			return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "sharing.NewShareRecipient:RecipientID")
		}
	default:
		return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "sharing.NewShareRecipient:RecipientType")
	}

	now := time.Now().UTC()
	return &ShareRecipient{
		ShareConfigID:  shareConfigID,
		RecipientType:  recipientType,
		RecipientID:    recipientID,
		Email:          email,
		DownloadsCount: 0,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

func (r *ShareRecipient) IncrementDownloads() {
	r.DownloadsCount++
	r.UpdatedAt = time.Now().UTC()
}
