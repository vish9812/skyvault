package store_file

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"skyvault/internal/domain/media"
	"skyvault/pkg/appconfig"
	"skyvault/pkg/apperror"
	"skyvault/pkg/utils"
)

type Store struct {
	app *appconfig.App
	FS  media.Storage
}

func NewStore(app *appconfig.App) *Store {
	return &Store{
		app: app,
		FS:  NewLocal(app),
	}
}

func (s *Store) Cleanup() error {
	return nil
}

// Health checks the health of the file storage
func (s *Store) Health(ctx context.Context) error {
	// Create a random owner directory
	baseDir := filepath.Join(s.app.Config.Server.DataDir, "uploads")
	randomOwner := utils.RandomName()
	randomOwnerDir := filepath.Join(baseDir, randomOwner)
	err := os.MkdirAll(randomOwnerDir, os.ModePerm)
	if err != nil {
		return apperror.NewAppError(err, "s.Health:MkdirAll")
	}

	// Create a random file in the owner directory
	randomFile := utils.RandomName()
	randomFilePath := filepath.Join(randomOwnerDir, randomFile)
	err = os.WriteFile(randomFilePath, []byte("test"), os.ModePerm)
	if err != nil {
		return apperror.NewAppError(err, "s.Health:WriteFile")
	}
	randomFileContent, err := os.ReadFile(randomFilePath)
	if err != nil {
		return apperror.NewAppError(err, "s.Health:ReadFile")
	}
	if string(randomFileContent) != "test" {
		return apperror.NewAppError(errors.New("health check failed: file content is not correct"), "s.Health:ReadFile")
	}

	// Cleanup
	err = os.Remove(randomFilePath)
	if err != nil {
		return apperror.NewAppError(err, "s.Health:Remove")
	}
	err = os.RemoveAll(randomOwnerDir)
	if err != nil {
		return apperror.NewAppError(err, "s.Health:RemoveAll")
	}

	return nil
}
