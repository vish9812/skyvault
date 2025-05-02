package sharing

import (
	"time"
)

type Contact struct {
	ID        int64
	OwnerID   int64
	Email     string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewContact(ownerID int64, email string, name string) *Contact {
	now := time.Now().UTC()
	return &Contact{
		OwnerID:   ownerID,
		Email:     email,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
