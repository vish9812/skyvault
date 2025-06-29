package sharing

import (
	"time"

	"skyvault/pkg/apperror"
	"skyvault/pkg/utils"
)

type Contact struct {
	ID        string
	OwnerID   string
	Email     string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewContact(ownerID string, email string, name string) (*Contact, error) {
	id, err := utils.ID()
	if err != nil {
		return nil, apperror.NewAppError(err, "sharing.NewContact:ID")
	}
	now := time.Now().UTC()
	return &Contact{
		ID:        id,
		OwnerID:   ownerID,
		Email:     email,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
