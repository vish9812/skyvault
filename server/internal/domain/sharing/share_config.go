package sharing

import (
	"skyvault/pkg/apperror"
	"skyvault/pkg/utils"
	"time"
)

type ShareConfig struct {
	ID           int64
	CustomID     string // for url path
	OwnerID      int64
	FileID       *int64 // Only one of file or folder must be set
	FolderID     *int64
	PasswordHash *string
	MaxDownloads *int64
	ExpiresAt    *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Recipients   []*ShareRecipient
}

func NewShareConfig(
	ownerID int64,
	fileID *int64,
	folderID *int64,
	password *string,
	maxDownloads *int64,
	expiresAt *time.Time,
) (*ShareConfig, error) {
	var passwordHash *string
	if password != nil {
		h, err := utils.HashPassword(*password)
		if err != nil {
			return nil, apperror.NewAppError(err, "sharing.NewShareConfig:Password")
		}
		passwordHash = &h
	}

	now := time.Now().UTC()
	return &ShareConfig{
		CustomID:     utils.UUID(),
		OwnerID:      ownerID,
		FileID:       fileID,
		FolderID:     folderID,
		PasswordHash: passwordHash,
		MaxDownloads: maxDownloads,
		ExpiresAt:    expiresAt,
		CreatedAt:    now,
		UpdatedAt:    now,
		Recipients:   []*ShareRecipient{},
	}, nil
}

// App Errors:
// - ErrCommonNoAccess
func (s *ShareConfig) ValidateAccess(ownerID int64) error {
	if s.OwnerID != ownerID {
		return apperror.NewAppError(apperror.ErrCommonNoAccess, "sharing.ShareConfig.ValidateAccess")
	}
	return nil
}

func (s *ShareConfig) ValidateExpiry() error {
	if s.ExpiresAt != nil && s.ExpiresAt.Before(time.Now().UTC()) {
		return apperror.NewAppError(apperror.ErrSharingExpired, "sharing.ShareConfig.ValidateExpiry")
	}
	return nil
}
