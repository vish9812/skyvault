package profile

import (
	"fmt"
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

// App Errors:
// - ErrCommonInvalidValue (when quota < 0)
func NewProfile(email, fullName string, quota int64) (*Profile, error) {
	if quota < 0 {
		return nil, apperror.NewAppError(fmt.Errorf("%w: quota must be non-negative", apperror.ErrCommonInvalidValue), "profile.NewProfile:Quota").
			WithMetadata("quota", quota)
	}

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
