package profile

import (
	"skyvault/pkg/apperror"
	"skyvault/pkg/utils"
	"time"
)

type Profile struct {
	ID                string
	Email             string
	FullName          string
	Avatar            []byte
	StorageQuotaBytes int64
	StorageUsedBytes  int64
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func NewProfile(email, fullName string) (*Profile, error) {
	id, err := utils.ID()
	if err != nil {
		return nil, apperror.NewAppError(err, "profile.NewProfile:ID")
	}

	now := time.Now().UTC()
	return &Profile{
		ID:                id,
		Email:             email,
		FullName:          fullName,
		Avatar:            nil,
		StorageQuotaBytes: 0, // Set by caller (e.g., signup workflow)
		StorageUsedBytes:  0,
		CreatedAt:         now,
		UpdatedAt:         now,
	}, nil
}

// App Errors:
// - ErrCommonNoAccess
func (p *Profile) ValidateAccess(accessedByID string) error {
	if p.ID != accessedByID {
		return apperror.ErrCommonNoAccess
	}
	return nil
}

// GetAvailableStorage returns the available storage in bytes
func (p *Profile) GetAvailableStorage() int64 {
	available := p.StorageQuotaBytes - p.StorageUsedBytes
	if available < 0 {
		return 0
	}
	return available
}

// CanAllocate checks if the profile has enough storage quota to allocate the requested bytes
func (p *Profile) CanAllocate(bytes int64) bool {
	return p.GetAvailableStorage() >= bytes
}

// SetQuota sets the storage quota for this profile
func (p *Profile) SetQuota(quotaBytes int64) {
	p.StorageQuotaBytes = quotaBytes
}
