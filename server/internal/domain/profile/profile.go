package profile

import (
	"skyvault/pkg/apperror"
	"time"
)

type Profile struct {
	ID        int64
	Email     string
	FullName  string
	Avatar    []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewProfile(email, fullName string) *Profile {
	now := time.Now().UTC()
	return &Profile{
		Email:     email,
		FullName:  fullName,
		Avatar:    nil,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// App Errors:
// - ErrCommonNoAccess
func (p *Profile) ValidateAccess(accessedByID int64) error {
	if p.ID != accessedByID {
		return apperror.ErrCommonNoAccess
	}
	return nil
}
