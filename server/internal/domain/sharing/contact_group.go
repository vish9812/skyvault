package sharing

import (
	"time"
)

type ContactGroup struct {
	ID        int64
	OwnerID   int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewContactGroup(ownerID int64, name string) *ContactGroup {
	now := time.Now().UTC()
	return &ContactGroup{
		OwnerID:   ownerID,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
