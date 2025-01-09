package services

import (
	"skyvault/internal/domain/auth"
	"skyvault/internal/domain/profile"
)

func NewAuthSvc(authRepo auth.Repo, profileRepo profile.Repo, jwt *auth.JWT) IAuthSvc {
	svc := newAuthSvc(authRepo, profileRepo, jwt)
	return newAuthSvcValidator(svc)
}
