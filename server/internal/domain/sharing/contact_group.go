package sharing

import (
	"skyvault/pkg/apperror"
	"skyvault/pkg/utils"
	"time"
)

type ContactGroup struct {
	ID        string
	OwnerID   string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewContactGroup(ownerID string, name string) (*ContactGroup, error) {
	id, err := utils.ID()
	if err != nil {
		return nil, apperror.NewAppError(err, "sharing.NewContactGroup:ID")
	}

	now := time.Now().UTC()
	return &ContactGroup{
		ID:        id,
		OwnerID:   ownerID,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
