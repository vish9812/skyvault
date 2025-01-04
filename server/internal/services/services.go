package services

import (
	"skyvault/internal/domain/auth"
	"skyvault/internal/domain/profile"
)

func NewAuthSvc(authRepo auth.Repo, profileRepo profile.Repo) IAuthSvc {
	svc := newAuthSvc(authRepo, profileRepo)
	return newAuthSvcValidator(svc)
}
