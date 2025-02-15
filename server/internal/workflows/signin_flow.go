package workflows

import (
	"context"
	"skyvault/internal/domain/auth"
	"skyvault/internal/domain/profile"
	"skyvault/pkg/apperror"
)

type SignInFlow struct {
	authCommands    auth.Commands
	authQueries     auth.Queries
	profileCommands profile.Commands
	profileQueries  profile.Queries
}

func NewSignInFlow(authCommands auth.Commands, authQueries auth.Queries, profileCommands profile.Commands, profileQueries profile.Queries) *SignInFlow {
	return &SignInFlow{
		authCommands:    authCommands,
		authQueries:     authQueries,
		profileCommands: profileCommands,
		profileQueries:  profileQueries,
	}
}

type SignInReq struct {
	Provider       auth.Provider
	ProviderUserID string
	Password       *string
}

type SignInRes struct {
	Profile *profile.Profile
	Token   string
}

// App Errors:
// - ErrCommonNoData
// - ErrCommonNoAccess
// - ErrCommonInvalidValue
// - ErrAuthWrongProvider
// - ErrAuthInvalidCredentials
func (f *SignInFlow) Run(ctx context.Context, req *SignInReq) (*SignInRes, error) {
	// 1. Get auth by provider and provider user id
	// 2. Get profile by profile id
	// 3. Signin to get token
	// 4. Return response

	au, err := f.authQueries.GetByProvider(ctx, &auth.GetByProviderQuery{
		Provider:       req.Provider,
		ProviderUserID: req.ProviderUserID,
	})
	if err != nil {
		return nil, apperror.NewAppError(err, "SignInFlow.Run:authQueries.GetByProvider")
	}

	pro, err := f.profileQueries.Get(ctx, &profile.GetQuery{ID: au.ProfileID})
	if err != nil {
		return nil, apperror.NewAppError(err, "SignInFlow.Run:profileQueries.Get")
	}

	token, err := f.authCommands.SignIn(ctx, &auth.SignInCommand{
		ProfileID:    au.ProfileID,
		Provider:     au.Provider,
		Password:     req.Password,
		PasswordHash: au.PasswordHash,
	})
	if err != nil {
		return nil, apperror.NewAppError(err, "SignInFlow.Run:authCommands.SignIn")
	}

	return &SignInRes{
		Profile: pro,
		Token:   token,
	}, nil
}
