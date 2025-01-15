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
	"skyvault/pkg/common"

	"github.com/rs/zerolog/log"
)

var _ media.Storage = (*Local)(nil)

type Local struct {
	app     *common.App
	baseDir string
}

func NewLocal(app *common.App) *Local {
	// Ensure the base directory exists
	baseDir := filepath.Join(app.Config.APP_DATA_FOLDER, "uploads")
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
		return common.NewAppErr(fmt.Errorf("failed to create owner directory: %w", err), "SaveFile")
	}

	// Create the destination file
	savePath := filepath.Join(ownerDir, name)
	f, err := os.Create(savePath)
	if err != nil {
		return common.NewAppErr(fmt.Errorf("failed to create file: %w", err), "SaveFile")
	}
	defer f.Close()

	// Write the file
	_, err = io.Copy(f, file)
	if err != nil {
		// Delete the file if writing fails
		if e := ls.DeleteFile(ctx, name, ownerID); e != nil {
			return common.NewAppErr(fmt.Errorf("%w: %w", e, err), "SaveFile")
		}
		return common.NewAppErr(fmt.Errorf("failed to write to file: %w", err), "SaveFile")
	}

	return nil
}

func (ls *Local) DeleteFile(ctx context.Context, name string, ownerID int64) error {
	savePath := filepath.Join(ls.baseDir, fmt.Sprintf("%d", ownerID), name)
	if err := os.Remove(savePath); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}

		return common.NewAppErr(fmt.Errorf("failed to delete file: %w", err), "DeleteFile")
	}
	return nil
}

func (ls *Local) OpenFile(ctx context.Context, name string, ownerID int64) (io.ReadSeekCloser, error) {
	savePath := filepath.Join(ls.baseDir, fmt.Sprintf("%d", ownerID), name)
	f, err := os.Open(savePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, common.NewAppErr(media.ErrFileNotFound, "OpenFile")
		}

		return nil, common.NewAppErr(fmt.Errorf("failed to open file: %w", err), "OpenFile")
	}
	return f, nil
}
