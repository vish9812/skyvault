package sharing

import (
	"skyvault/pkg/apperror"
	"time"
)

type ResourceType string

const (
	ResourceTypeFile   ResourceType = "file"
	ResourceTypeFolder ResourceType = "folder"
)

type ShareConfig struct {
	ID           int64
	OwnerID      int64
	ResourceType ResourceType
	ResourceID   int64
	PasswordHash *string
	MaxDownloads *int
	ExpiresAt    *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Recipients   []*ShareRecipient
}

func NewShareConfig(
	ownerID int64,
	resourceType ResourceType,
	resourceID int64,
	passwordHash *string,
	maxDownloads *int,
	expiresAt *time.Time,
) (*ShareConfig, error) {
	switch resourceType {
	case ResourceTypeFile, ResourceTypeFolder:
		// valid
	default:
		return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "sharing.NewShareConfig:ResourceType")
	}

	if maxDownloads != nil && *maxDownloads <= 0 {
		return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "sharing.NewShareConfig:MaxDownloads")
	}

	now := time.Now().UTC()
	if expiresAt != nil && expiresAt.Before(now) {
		return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "sharing.NewShareConfig:ExpiresAt")
	}

	return &ShareConfig{
		OwnerID:      ownerID,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		PasswordHash: passwordHash,
		MaxDownloads: maxDownloads,
		ExpiresAt:    expiresAt,
		CreatedAt:    now,
		UpdatedAt:    now,
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

func (s *ShareConfig) UpdateConfig(
	passwordHash *string,
	maxDownloads *int,
	expiresAt *time.Time,
) error {
	if maxDownloads != nil && *maxDownloads <= 0 {
		return apperror.NewAppError(apperror.ErrCommonInvalidValue, "sharing.ShareConfig.UpdateConfig:MaxDownloads")
	}

	now := time.Now().UTC()
	if expiresAt != nil && expiresAt.Before(now) {
		return apperror.NewAppError(apperror.ErrCommonInvalidValue, "sharing.ShareConfig.UpdateConfig:ExpiresAt")
	}

	s.PasswordHash = passwordHash
	s.MaxDownloads = maxDownloads
	s.ExpiresAt = expiresAt
	s.UpdatedAt = now
	return nil
}
