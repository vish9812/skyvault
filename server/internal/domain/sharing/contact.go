package sharing

import (
	"skyvault/pkg/apperror"
	"strings"
	"time"
)

type Contact struct {
	ID        int64
	OwnerID   int64
	Email     string
	Name      *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewContact(ownerID int64, email string, name *string) (*Contact, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "sharing.NewContact:Email")
	}

	if name != nil {
		trimmedName := strings.TrimSpace(*name)
		if trimmedName == "" {
			name = nil
		} else {
			name = &trimmedName
		}
	}

	now := time.Now().UTC()
	return &Contact{
		OwnerID:   ownerID,
		Email:     email,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (c *Contact) UpdateName(name *string) {
	if name != nil {
		trimmedName := strings.TrimSpace(*name)
		if trimmedName == "" {
			name = nil
		} else {
			name = &trimmedName
		}
	}

	c.Name = name
	c.UpdatedAt = time.Now().UTC()
}
