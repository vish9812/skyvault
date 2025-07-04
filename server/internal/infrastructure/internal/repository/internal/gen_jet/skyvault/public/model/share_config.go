//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"github.com/google/uuid"
	"time"
)

type ShareConfig struct {
	ID           uuid.UUID `sql:"primary_key"`
	OwnerID      uuid.UUID
	FileID       *uuid.UUID
	FolderID     *uuid.UUID
	PasswordHash *string
	MaxDownloads *int64
	ExpiresAt    *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
