package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"skyvault/internal/domain/media"
	"skyvault/pkg/appconfig"
	"skyvault/pkg/apperror"
)

const localStorageBaseDir = "uploads"

var _ media.Storage = (*LocalStorage)(nil)

type LocalStorage struct {
	app     *appconfig.App
	baseDir string
}

func NewLocalStorage(app *appconfig.App) *LocalStorage {
	// Ensure the base directory exists
	baseDir := filepath.Join(app.Config.Server.DataDir, localStorageBaseDir)
	err := os.MkdirAll(baseDir, os.ModePerm)
	if err != nil {
		app.Logger.Fatal().Err(err).Str("base_dir", baseDir).Msg("Failed to create base directory for local storage")
	}

	return &LocalStorage{app: app, baseDir: baseDir}
}

func (s *LocalStorage) SaveFile(ctx context.Context, file io.ReadSeeker, name string, ownerID int64) error {
	// Create the owner directory
	ownerDir := getOwnerDir(s.baseDir, ownerID)
	err := os.MkdirAll(ownerDir, os.ModePerm)
	if err != nil {
		return apperror.NewAppError(err, "LocalStorage.SaveFile:MkdirAll").WithMetadata("owner_dir", ownerDir)
	}

	savePath := getFilePath(ownerDir, name)

	// Check if the file already exists
	if _, err := os.Stat(savePath); err == nil {
		return apperror.NewAppError(apperror.ErrCommonDuplicateData, "LocalStorage.SaveFile:Stat").WithMetadata("save_path", savePath)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return apperror.NewAppError(err, "LocalStorage.SaveFile:Stat").WithMetadata("save_path", savePath)
	}

	// Create the destination file
	f, err := os.Create(savePath)
	if err != nil {
		return apperror.NewAppError(err, "LocalStorage.SaveFile:Create").WithMetadata("save_path", savePath)
	}
	defer f.Close()

	// Write the file
	maxSize := s.app.Config.Media.MaxSizeMB * media.BytesPerMB
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return apperror.NewAppError(err, "LocalStorage.SaveFile:Seek").WithMetadata("save_path", savePath)
	}

	// Copy the file with an extra MB to check if file is greater than max size.
	written, err := io.CopyN(f, file, maxSize+media.BytesPerMB)
	if err != nil && !errors.Is(err, io.EOF) {
		// Clean up the partially written file if copying fails
		if delErr := s.DeleteFile(ctx, name, ownerID); delErr != nil {
			if !errors.Is(delErr, apperror.ErrCommonNoData) {
				return apperror.NewAppError(fmt.Errorf("%w: %w", delErr, err), "LocalStorage.SaveFile:CopyN:DeleteFile")
			}
		}
		return apperror.NewAppError(err, "LocalStorage.SaveFile:CopyN").WithMetadata("save_path", savePath)
	}

	if written > maxSize {
		var sizeErr error = apperror.ErrMediaFileSizeLimitExceeded
		// Clean up the written file if it is greater than max size
		if delErr := s.DeleteFile(ctx, name, ownerID); delErr != nil {
			sizeErr = apperror.NewAppError(fmt.Errorf("%w: %w", delErr, sizeErr), "LocalStorage.SaveFile:MaxSizeExceeded:DeleteFile")
		}
		return apperror.NewAppError(sizeErr, "LocalStorage.SaveFile:MaxSizeExceeded").WithMetadata("save_path", savePath).WithMetadata("written_mb", written/media.BytesPerMB).WithMetadata("max_size_mb", maxSize/media.BytesPerMB)
	}

	return nil
}

func (s *LocalStorage) DeleteFile(ctx context.Context, name string, ownerID int64) error {
	deletePath := getFilePath(getOwnerDir(s.baseDir, ownerID), name)
	if err := os.Remove(deletePath); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonNoData, err), "LocalStorage.DeleteFile:Remove").WithMetadata("delete_path", deletePath)
		}

		return apperror.NewAppError(err, "LocalStorage.DeleteFile:Remove").WithMetadata("delete_path", deletePath)
	}
	return nil
}

func (s *LocalStorage) OpenFile(ctx context.Context, name string, ownerID int64) (io.ReadSeekCloser, error) {
	openPath := getFilePath(getOwnerDir(s.baseDir, ownerID), name)
	f, err := os.Open(openPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonNoData, err), "LocalStorage.OpenFile:Open").WithMetadata("open_path", openPath)
		}

		return nil, apperror.NewAppError(err, "LocalStorage.OpenFile:Open").WithMetadata("open_path", openPath)
	}
	return f, nil
}

func getOwnerDir(baseDir string, ownerID int64) string {
	return filepath.Join(baseDir, fmt.Sprintf("%d", ownerID))
}

func getFilePath(ownerDir string, name string) string {
	return filepath.Join(ownerDir, name)
}
