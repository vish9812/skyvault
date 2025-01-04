package services

import (
	"context"
	"skyvault/internal/domain/auth"
	"skyvault/internal/domain/profile"
	"skyvault/pkg/common"

	"github.com/jinzhu/copier"
)

type AuthSvc struct {
	authRepo    auth.Repo
	profileRepo profile.Repo
}

func newAuthSvc(authRepo auth.Repo, profileRepo profile.Repo) IAuthSvc {
	return &AuthSvc{
		authRepo:    authRepo,
		profileRepo: profileRepo,
	}
}

type SignUpReq struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type SignUpResp struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (s *AuthSvc) SignUp(ctx context.Context, req *SignUpReq) (*SignUpResp, error) {
	// Check if profile already exists
	_, err := s.profileRepo.GetByEmail(ctx, req.Email)
	if err == nil {
		return nil, common.NewValidationError(profile.ErrProfileAlreadyExists)
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create profile
	pro := profile.NewProfile()
	err = copier.Copy(pro, req)
	if err != nil {
		return nil, err
	}

	tx, err := s.profileRepo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	profileRepoTx := s.profileRepo.WithTx(ctx, tx)
	authRepoTx := s.authRepo.WithTx(ctx, tx)

	pro, err = profileRepoTx.Create(ctx, pro)
	if err != nil {
		return nil, err
	}

	// Create auth
	au := auth.NewAuth(pro.ID)
	au.ProviderUserID = pro.Email
	au.PasswordHash = &passwordHash

	_, err = authRepoTx.Create(ctx, au)
	if err != nil {
		return nil, err
	}

	var resp *SignUpResp
	err = copier.Copy(resp, pro)
	if err != nil {
		return nil, err
	}

	tx.Commit()

	return resp, nil
}
