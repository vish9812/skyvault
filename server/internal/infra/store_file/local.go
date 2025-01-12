package store_file

import (
	"context"
	"fmt"
	"io"
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

func (ls *Local) CreateFile(ctx context.Context, name string, reader io.Reader, ownerID int64) error {
	// Create the owner directory
	ownerDir := filepath.Join(ls.baseDir, fmt.Sprintf("%d", ownerID))
	err := os.MkdirAll(ownerDir, os.ModePerm)
	if err != nil {
		return common.NewAppErr(fmt.Errorf("failed to create owner directory: %w", err), "CreateFile")
	}

	// Create the destination file
	savePath := filepath.Join(ls.baseDir, name)
	createdFile, err := os.Create(savePath)
	if err != nil {
		return common.NewAppErr(fmt.Errorf("failed to create file: %w", err), "CreateFile")
	}
	defer createdFile.Close()

	_, err = io.Copy(createdFile, reader)
	if err != nil {
		// Delete the file if writing fails
		if e := ls.DeleteFile(ctx, name, ownerID); e != nil {
			return common.NewAppErr(fmt.Errorf("%w: %w", e, err), "CreateFile")
		}
		return common.NewAppErr(fmt.Errorf("failed to write to file: %w", err), "CreateFile")
	}

	return nil
}

func (ls *Local) DeleteFile(ctx context.Context, name string, ownerID int64) error {
	savePath := filepath.Join(ls.baseDir, fmt.Sprintf("%d", ownerID), name)
	if err := os.Remove(savePath); err != nil {
		return common.NewAppErr(fmt.Errorf("failed to delete file: %w", err), "DeleteFile")
	}
	return nil
}

func (ls *Local) GetFile(ctx context.Context, name string, ownerID int64) (io.ReadCloser, error) {
	savePath := filepath.Join(ls.baseDir, fmt.Sprintf("%d", ownerID), name)
	file, err := os.Open(savePath)
	if err != nil {
		return nil, common.NewAppErr(fmt.Errorf("failed to open file: %w", err), "GetFile")
	}
	return file, nil
}
