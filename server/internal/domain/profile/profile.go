package profile

import (
	"skyvault/pkg/apperror"
	"skyvault/pkg/utils"
	"time"
)

type Profile struct {
	ID           string
	Email        string
	FullName     string
	Avatar       []byte
	StorageQuota int64
	StorageUsed  int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewProfile(email, fullName string, quota int64) (*Profile, error) {
	id, err := utils.ID()
	if err != nil {
		return nil, apperror.NewAppError(err, "profile.NewProfile:ID")
	}

	now := time.Now().UTC()
	return &Profile{
		ID:           id,
		Email:        email,
		FullName:     fullName,
		Avatar:       nil,
		StorageQuota: quota,
		StorageUsed:  0,
		CreatedAt:    now,
		UpdatedAt:    now,
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

// CanAllocate checks if the profile has enough storage quota to allocate the requested bytes
func (p *Profile) CanAllocate(bytes int64) bool {
	return (p.StorageQuota - p.StorageUsed) >= bytes
}
