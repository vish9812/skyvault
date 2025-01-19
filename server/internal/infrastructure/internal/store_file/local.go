package store_file

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

	"github.com/rs/zerolog/log"
)

var _ media.Storage = (*Local)(nil)

type Local struct {
	app     *appconfig.App
	baseDir string
}

func NewLocal(app *appconfig.App) *Local {
	// Ensure the base directory exists
	baseDir := filepath.Join(app.Config.Server.DataDir, "uploads")
	err := os.MkdirAll(baseDir, os.ModePerm)
	if err != nil {
		log.Fatal().Err(err).Str("base_dir", baseDir).Msg("Failed to create base directory for local storage")
	}

	return &Local{baseDir: baseDir, app: app}
}

func (ls *Local) SaveFile(ctx context.Context, file io.Reader, name string, ownerID int64) error {
	// Create the owner directory
	ownerDir := filepath.Join(ls.baseDir, fmt.Sprintf("%d", ownerID))
	err := os.MkdirAll(ownerDir, os.ModePerm)
	if err != nil {
		return apperror.NewAppError(err, "store_file.SaveFile:MkdirAll").WithMetadata("owner_dir", ownerDir)
	}

	savePath := filepath.Join(ownerDir, name)

	// Check if the file already exists
	if _, err := os.Stat(savePath); err == nil {
		return apperror.NewAppError(apperror.ErrDuplicateData, "store_file.SaveFile:Stat").WithMetadata("save_path", savePath)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return apperror.NewAppError(err, "store_file.SaveFile:Stat").WithMetadata("save_path", savePath)
	}

	// Create the destination file
	f, err := os.Create(savePath)
	if err != nil {
		return apperror.NewAppError(err, "store_file.SaveFile:Create").WithMetadata("save_path", savePath)
	}
	defer f.Close()

	// Write the file
	written, err := io.Copy(f, file)
	if err != nil {
		// Delete the file if writing fails
		if e := ls.DeleteFile(ctx, name, ownerID); e != nil {
			if !errors.Is(e, apperror.ErrNoData) {
				return apperror.NewAppError(fmt.Errorf("%w: %w", e, err), "store_file.SaveFile:DeleteFile")
			}
		}
		return apperror.NewAppError(err, "store_file.SaveFile:Copy").WithMetadata("save_path", savePath)
	}

	if written > ls.app.Config.Media.MaxSizeMB*media.BytesPerMB {
		return apperror.NewAppError(media.ErrFileSizeLimitExceeded, "store_file.SaveFile").WithMetadata("save_path", savePath)
	}

	return nil
}

func (ls *Local) DeleteFile(ctx context.Context, name string, ownerID int64) error {
	deletePath := filepath.Join(ls.baseDir, fmt.Sprintf("%d", ownerID), name)
	if err := os.Remove(deletePath); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrNoData, err), "store_file.DeleteFile:Remove").WithMetadata("delete_path", deletePath)
		}

		return apperror.NewAppError(err, "store_file.DeleteFile:Remove").WithMetadata("delete_path", deletePath)
	}
	return nil
}

func (ls *Local) OpenFile(ctx context.Context, name string, ownerID int64) (io.ReadSeekCloser, error) {
	openPath := filepath.Join(ls.baseDir, fmt.Sprintf("%d", ownerID), name)
	f, err := os.Open(openPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrNoData, err), "store_file.OpenFile:Open").WithMetadata("open_path", openPath)
		}

		return nil, apperror.NewAppError(err, "store_file.OpenFile:Open").WithMetadata("open_path", openPath)
	}
	return f, nil
}