package sharing

import (
	"skyvault/pkg/apperror"
	"strings"
	"time"
)

type ContactGroup struct {
	ID        int64
	OwnerID   int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewContactGroup(ownerID int64, name string) (*ContactGroup, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, apperror.NewAppError(apperror.ErrCommonInvalidValue, "sharing.NewContactGroup:Name")
	}

	now := time.Now().UTC()
	return &ContactGroup{
		OwnerID:   ownerID,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (g *ContactGroup) Rename(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return apperror.NewAppError(apperror.ErrCommonInvalidValue, "sharing.ContactGroup.Rename:Name")
	}

	g.Name = name
	g.UpdatedAt = time.Now().UTC()
	return nil
}

// App Errors:
// - ErrCommonNoAccess
func (g *ContactGroup) ValidateAccess(ownerID int64) error {
	if g.OwnerID != ownerID {
		return apperror.NewAppError(apperror.ErrCommonNoAccess, "sharing.ContactGroup.ValidateAccess")
	}
	return nil
}
