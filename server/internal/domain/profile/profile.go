package profile

import "errors"

var ErrProfileAlreadyExists = errors.New("profile already exists")

type Profile struct {
	ID       int64
	Email    string
	FullName string
}

func NewProfile() *Profile {
	return &Profile{}
}
