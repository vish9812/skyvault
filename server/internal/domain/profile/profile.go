package profile

import (
	"skyvault/pkg/apperror"
	"skyvault/pkg/utils"
	"time"
)

type Profile struct {
	ID        string
	Email     string
	FullName  string
	Avatar    []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewProfile(email, fullName string) (*Profile, error) {
	id, err := utils.ID()
	if err != nil {
		return nil, apperror.NewAppError(err, "profile.NewProfile:ID")
	}

	now := time.Now().UTC()
	return &Profile{
		ID:        id,
		Email:     email,
		FullName:  fullName,
		Avatar:    nil,
		CreatedAt: now,
		UpdatedAt: now,
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
