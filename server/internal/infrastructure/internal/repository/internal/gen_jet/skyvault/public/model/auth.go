//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"time"
)

type Auth struct {
	ID             int64 `sql:"primary_key"`
	ProfileID      int64
	Provider       string
	ProviderUserID string
	PasswordHash   *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
