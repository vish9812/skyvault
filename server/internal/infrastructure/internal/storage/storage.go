package storage

import (
	"context"
	"errors"
	"math"
	"math/rand/v2"
	"os"
	"path/filepath"
	"skyvault/internal/domain/media"
	"skyvault/pkg/appconfig"
	"skyvault/pkg/apperror"
	"skyvault/pkg/utils"
)

type Storage struct {
	app          *appconfig.App
	LocalStorage media.Storage
}

func NewStorage(app *appconfig.App) *Storage {
	return &Storage{
		app:          app,
		LocalStorage: NewLocalStorage(app),
	}
}

func (s *Storage) Cleanup() error {
	return nil
}

// Health checks the health of the file storage
func (s *Storage) Health(ctx context.Context) error {
	// Create a random owner directory
	baseDir := filepath.Join(s.app.Config.Server.DataDir, localStorageBaseDir)

	// Use a random int64 value from [max-100, max)
	randomOwnerID := int64(math.MaxInt64 - rand.IntN(100))
	ownerDir := getOwnerDir(baseDir, randomOwnerID)
	err := os.MkdirAll(ownerDir, os.ModePerm)
	if err != nil {
		return apperror.NewAppError(err, "s.Health:MkdirAll").WithMetadata("owner_dir", ownerDir)
	}

	// Create a random file in the owner directory
	randomFile := utils.RandomName()
	randomFilePath := getFilePath(ownerDir, randomFile)
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
	err = os.RemoveAll(ownerDir)
	if err != nil {
		return apperror.NewAppError(err, "s.Health:RemoveAll")
	}

	return nil
}
