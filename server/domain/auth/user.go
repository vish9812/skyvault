package auth

import "skyvault/common/utils"

type User struct {
	ID           utils.ID
	FirstName    string
	LastName     string
	Email        string
	PasswordHash string
}

func NewUser() *User {
	return &User{
		ID: utils.NewID(),
	}
}
