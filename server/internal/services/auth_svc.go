package services

import (
	"context"
	"errors"
	"fmt"
	"skyvault/internal/domain/auth"
	"skyvault/internal/domain/profile"
	"skyvault/pkg/common"

	"github.com/jinzhu/copier"
)

var ErrEmailOrPasswordIncorrect = common.NewValidationError(errors.New("email or password is incorrect"))

var _ IAuthSvc = (*AuthSvc)(nil)

type AuthSvc struct {
	authRepo    auth.Repo
	profileRepo profile.Repo
	jwt         *auth.JWT
}

func newAuthSvc(authRepo auth.Repo, profileRepo profile.Repo, jwt *auth.JWT) IAuthSvc {
	return &AuthSvc{
		authRepo:    authRepo,
		profileRepo: profileRepo,
		jwt:         jwt,
	}
}

type SignUpReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"fullName"`
}

type AuthSuccessResp struct {
	Token   string             `json:"token"`
	Profile *SignUpRespProfile `json:"profile"`
}

type SignUpRespProfile struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"fullName"`
}

func (s *AuthSvc) SignUp(ctx context.Context, req *SignUpReq) (*AuthSuccessResp, error) {
	// Check if profile already exists
	_, err := s.profileRepo.GetByEmail(ctx, req.Email)
	if err == nil {
		return nil, common.NewAppErr(common.NewValidationError(profile.ErrProfileAlreadyExists), "SignUp")
	}
	
	// Create profile
	pro := profile.NewProfile()
	err = copier.Copy(pro, req)
	if err != nil {
		return nil, common.NewAppErr(err, "SignUp")
	}
	
	tx, err := s.profileRepo.BeginTx(ctx)
	if err != nil {
		return nil, common.NewAppErr(err, "SignUp")
	}
	defer tx.Rollback()
	
	profileRepoTx := s.profileRepo.WithTx(ctx, tx)
	authRepoTx := s.authRepo.WithTx(ctx, tx)
	
	pro, err = profileRepoTx.Create(ctx, pro)
	if err != nil {
		return nil, common.NewAppErr(err, "SignUp")
	}
	
	// Create auth
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, common.NewAppErr(err, "SignUp")
	}
	au := auth.NewAuth(pro.ID)
	au.ProviderUserID = pro.Email
	au.PasswordHash = &passwordHash

	_, err = authRepoTx.Create(ctx, au)
	if err != nil {
		return nil, common.NewAppErr(fmt.Errorf("failed to create auth: %w", err), "SignUp")
	}

	resp := &AuthSuccessResp{Profile: &SignUpRespProfile{}}
	err = copier.Copy(resp.Profile, pro)
	if err != nil {
		return nil, common.NewAppErr(fmt.Errorf("failed to copy the struct profile: %w", err), "SignUp")
	}

	tx.Commit()

	// Generate JWT
	token, err := s.jwt.Generate(pro.ID, pro.Email)
	if err != nil {
		return nil, common.NewAppErr(err, "SignUp")
	}
	resp.Token = token

	return resp, nil
}

type SignInReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *AuthSvc) SignIn(ctx context.Context, req *SignInReq) (*AuthSuccessResp, error) {
	pro, err := s.profileRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, common.ErrDBNoRows) {
			return nil, common.NewAppErr(fmt.Errorf("%w: %w", ErrEmailOrPasswordIncorrect, err), "SignIn")
		}
		return nil, common.NewAppErr(fmt.Errorf("failed to get profile by email: %w", err), "SignIn")
	}

	au, err := s.authRepo.GetByProfileID(ctx, pro.ID)
	if err != nil {
		return nil, common.NewAppErr(fmt.Errorf("failed to get auth by profile id: %w", err), "SignIn")
	}

	if ok, err := auth.IsValidPassword(*au.PasswordHash, req.Password); err != nil || !ok {
		if err != nil {
			return nil, common.NewAppErr(fmt.Errorf("failed to validate password: %w", err), "SignIn")
		}
		return nil, common.NewAppErr(ErrEmailOrPasswordIncorrect, "SignIn")
	}

	resp := &AuthSuccessResp{Profile: &SignUpRespProfile{}}
	err = copier.Copy(resp.Profile, pro)
	if err != nil {
		return nil, common.NewAppErr(fmt.Errorf("failed to copy the struct profile: %w", err), "SignIn")
	}

	// Generate JWT
	token, err := s.jwt.Generate(pro.ID, pro.Email)
	if err != nil {
		return nil, common.NewAppErr(err, "SignIn")
	}
	resp.Token = token

	return resp, nil
}
