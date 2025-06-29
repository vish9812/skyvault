package sharing

import (
	"skyvault/pkg/apperror"
	"skyvault/pkg/utils"
	"time"
)

type ShareConfig struct {
	ID           string
	OwnerID      string
	FileID       *string // Only one of file or folder must be set
	FolderID     *string
	PasswordHash *string
	MaxDownloads *int64
	ExpiresAt    *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Recipients   []*ShareRecipient
}

func NewShareConfig(
	ownerID string,
	fileID *string,
	folderID *string,
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

	id, err := utils.ID()
	if err != nil {
		return nil, apperror.NewAppError(err, "sharing.NewShareConfig:ID")
	}

	now := time.Now().UTC()
	return &ShareConfig{
		ID:           id,
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
func (s *ShareConfig) ValidateAccess(ownerID string) error {
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
