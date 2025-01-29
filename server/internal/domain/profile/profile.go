package profile

import (
	"fmt"
	"skyvault/pkg/apperror"
	"skyvault/pkg/utils"
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

func NewProfile(email, fullName string) (*Profile, error) {
	if fullName == "" {
		return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "profile.NewProfile:FullName")
	}

	if err := utils.IsValidEmail(email); err != nil {
		return nil, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err), "profile.NewProfile:Email:Invalid")
	}

	return &Profile{
		Email:     email,
		FullName:  fullName,
		Avatar:    nil,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
