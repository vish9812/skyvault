package profile

import "errors"

var ErrProfileAlreadyExists = errors.New("profile already exists")

type Profile struct {
	ID        int64
	Email     string
	FirstName string
	LastName  string
}

func NewProfile() *Profile {
	return &Profile{}
}
