package workflows

import (
	"context"
	"skyvault/internal/domain/auth"
	"skyvault/internal/domain/profile"
	"skyvault/pkg/appconfig"
	"skyvault/pkg/apperror"
)

type SignUpFlow struct {
	app               *appconfig.App
	authCommands      auth.Commands
	authRepository    auth.Repository
	profileCommands   profile.Commands
	profileRepository profile.Repository
}

func NewSignUpFlow(app *appconfig.App, authCommands auth.Commands, authRepository auth.Repository, profileCommands profile.Commands, profileRepository profile.Repository) *SignUpFlow {
	return &SignUpFlow{
		app:               app,
		authCommands:      authCommands,
		authRepository:    authRepository,
		profileCommands:   profileCommands,
		profileRepository: profileRepository,
	}
}

type SignUpReq struct {
	Email          string
	FullName       string
	Provider       auth.Provider
	ProviderUserID string
	Password       *string
}

type SignUpRes struct {
	Profile *profile.Profile
	Token   string
}

// App Errors:
// - ErrCommonInvalidValue
// - ErrCommonDuplicateData
func (f *SignUpFlow) Run(ctx context.Context, req *SignUpReq) (*SignUpRes, error) {
	// 1. Create profile with transaction
	// 2. Create auth with transaction
	// 3. Generate JWT
	// 4. Commit transaction
	// 5. Return response

	tx, err := f.profileRepository.BeginTx(ctx)
	if err != nil {
		return nil, apperror.NewAppError(err, "SignUpFlow.Run:BeginTx")
	}
	defer tx.Rollback()

	profileRepoTx := f.profileRepository.WithTx(ctx, tx)
	authRepoTx := f.authRepository.WithTx(ctx, tx)

	profileCmdTx := f.profileCommands.WithTxRepository(ctx, profileRepoTx)
	authCmdTx := f.authCommands.WithTxRepository(ctx, authRepoTx)

	pro, err := profileCmdTx.Create(ctx, &profile.CreateCommand{
		Email:    req.Email,
		FullName: req.FullName,
	})
	if err != nil {
		return nil, apperror.NewAppError(err, "SignUpFlow.Run:Create")
	}

	// Set default storage quota for new user
	quotaBytes := f.app.Config.Storage.DefaultQuotaMB * 1024 * 1024 // Convert MB to bytes
	pro.SetQuota(quotaBytes)
	err = profileRepoTx.Update(ctx, pro)
	if err != nil {
		return nil, apperror.NewAppError(err, "SignUpFlow.Run:UpdateQuota")
	}

	token, err := authCmdTx.SignUp(ctx, &auth.SignUpCommand{
		ProfileID:      pro.ID,
		Provider:       req.Provider,
		ProviderUserID: req.ProviderUserID,
		Password:       req.Password,
	})
	if err != nil {
		return nil, apperror.NewAppError(err, "SignUpFlow.Run:SignUp")
	}

	tx.Commit()

	return &SignUpRes{
		Profile: pro,
		Token:   token,
	}, nil
}
