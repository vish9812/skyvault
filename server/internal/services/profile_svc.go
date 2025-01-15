package services

import (
	"context"
	"skyvault/internal/domain/profile"
)

type IProfileSvc interface {
	Get(ctx context.Context, id int64) (*profile.Profile, error)
}

type ProfileSvc struct {
	profileRepo profile.Repo
}

func NewProfileSvc(profileRepo profile.Repo) *ProfileSvc {
	return &ProfileSvc{profileRepo: profileRepo}
}

func (s *ProfileSvc) Get(ctx context.Context, id int64) (*profile.Profile, error) {
	return s.profileRepo.Get(ctx, id)
}
