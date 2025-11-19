package storage

import (
	"cmp"
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
	"skyvault/pkg/common"
	"slices"
	"strings"
	"time"
)

const localStorageBaseDir = "uploads"
const chunksDir = "chunks"

var _ media.Storage = (*LocalStorage)(nil)

type LocalStorage struct {
	app     *appconfig.App
	baseDir string
}

func NewLocalStorage(app *appconfig.App) *LocalStorage {
	// Ensure the base directory exists
	baseDir := filepath.Join(app.Config.Server.DataDir, localStorageBaseDir)
	err := os.MkdirAll(baseDir, 0700)
	if err != nil {
		app.Logger.Fatal().Err(err).Str("base_dir", baseDir).Msg("Failed to create base directory for local storage")
	}

	// Ensure chunks directory exists
	chunksPath := filepath.Join(baseDir, chunksDir)
	err = os.MkdirAll(chunksPath, 0700)
	if err != nil {
		app.Logger.Fatal().Err(err).Str("chunks_dir", chunksPath).Msg("Failed to create chunks directory for local storage")
	}

	return &LocalStorage{app: app, baseDir: baseDir}
}

func (s *LocalStorage) write(dirPath, filePath string, data io.Reader, maxSizeMB int64) error {
	// Create the directory
	err := os.MkdirAll(dirPath, 0700)
	if err != nil {
		return apperror.NewAppError(err, "storage.LocalStorage.write:MkdirAll")
	}

	// Check if the file already exists
	if _, err := os.Stat(filePath); err == nil {
		return apperror.NewAppError(apperror.ErrCommonDuplicateData, "storage.LocalStorage.write:Stat.Duplicate")
	} else if !errors.Is(err, fs.ErrNotExist) {
		return apperror.NewAppError(err, "storage.LocalStorage.write:Stat")
	}

	// Create the file
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return apperror.NewAppError(err, "storage.LocalStorage.write:OpenFile")
	}
	defer f.Close()

	// Copy data with an extra MB to check if data is greater than max size.
	written, err := io.CopyN(f, data, (maxSizeMB+1)*common.BytesPerMB)
	if err != nil && !errors.Is(err, io.EOF) {
		// Clean up on copy error
		errRemove := removeFile(filePath)
		if errRemove != nil {
			if !errors.Is(errRemove, apperror.ErrCommonNoData) {
				return apperror.NewAppError(fmt.Errorf("%w: %w", errRemove, err), "storage.LocalStorage.write:CopyN:removeFile")
			}
		}

		return apperror.NewAppError(err, "storage.LocalStorage.write:CopyN")
	}

	// Check size limit
	if written > maxSizeMB*common.BytesPerMB {
		var sizeErr error = apperror.ErrCommonInvalidValue
		// Clean up the written file if it is greater than max size
		errRemove := removeFile(filePath)
		if errRemove != nil {
			return apperror.NewAppError(fmt.Errorf("%w: %w", errRemove, sizeErr), "storage.LocalStorage.write:MaxSizeExceeded:removeFile")
		}

		return apperror.NewAppError(sizeErr, "storage.LocalStorage.write:MaxSizeExceeded").WithMetadata("written_mb", written/common.BytesPerMB).WithMetadata("max_size_mb", maxSizeMB)
	}

	return nil
}

func (s *LocalStorage) SaveFile(ctx context.Context, file io.ReadSeeker, name, ownerID string) error {
	ownerDirPath := getOwnerDirPath(s.baseDir, ownerID)
	filePath := getFilePath(ownerDirPath, name)

	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return apperror.NewAppError(err, "storage.LocalStorage.SaveFile:Seek").WithMetadata("file_path", filePath)
	}

	err = s.write(ownerDirPath, filePath, file, media.MaxDirectUploadSizeMB)
	if err != nil {
		return apperror.NewAppError(err, "storage.LocalStorage.SaveFile:write").WithMetadata("file_path", filePath)
	}

	return nil
}

func (s *LocalStorage) SaveChunk(ctx context.Context, chunk io.Reader, uploadID string, chunkIndex int64, ownerID string) error {
	chunksDirPath := getChunksDirPath(s.baseDir, ownerID, uploadID)
	chunkPath := getChunkPath(chunksDirPath, chunkIndex)

	err := s.write(chunksDirPath, chunkPath, chunk, media.MaxChunkSizeMB)
	if err != nil {
		return apperror.NewAppError(err, "storage.LocalStorage.SaveChunk:write").WithMetadata("chunk_path", chunkPath)
	}

	return nil
}

func (s *LocalStorage) FinalizeChunkedUpload(ctx context.Context, uploadID string, fileName string, ownerID string) error {
	chunksDirPath := getChunksDirPath(s.baseDir, ownerID, uploadID)

	// List all chunk files
	chunkFiles, err := filepath.Glob(filepath.Join(chunksDirPath, "chunk_*"))
	if err != nil {
		return apperror.NewAppError(err, "storage.LocalStorage.FinalizeChunkedUpload:Glob").WithMetadata("chunks_dir_path", chunksDirPath)
	}

	if len(chunkFiles) < 2 {
		return apperror.NewAppError(apperror.ErrCommonInvalidValue, "storage.LocalStorage.FinalizeChunkedUpload:NoChunks").WithMetadata("chunks_dir_path", chunksDirPath).WithMetadata("chunk_files_count", len(chunkFiles))
	}

	slices.SortFunc(chunkFiles, func(a, b string) int {
		return cmp.Compare(filepath.Base(a), filepath.Base(b))
	})

	// Create the owner directory
	ownerDirPath := getOwnerDirPath(s.baseDir, ownerID)
	err = os.MkdirAll(ownerDirPath, 0700)
	if err != nil {
		return apperror.NewAppError(err, "storage.LocalStorage.FinalizeChunkedUpload:MkdirAll").WithMetadata("owner_dir_path", ownerDirPath)
	}

	finalPath := getFilePath(ownerDirPath, fileName)

	// Check if the file already exists
	if _, err := os.Stat(finalPath); err == nil {
		return apperror.NewAppError(apperror.ErrCommonDuplicateData, "storage.LocalStorage.FinalizeChunkedUpload:Stat").WithMetadata("final_path", finalPath)
	}

	// Create the final file
	finalFile, err := os.OpenFile(finalPath, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return apperror.NewAppError(err, "storage.LocalStorage.FinalizeChunkedUpload:OpenFile").WithMetadata("final_path", finalPath)
	}
	defer finalFile.Close()

	// Combine all chunks
	var totalSize int64
	maxSize := int64(media.MaxFileSizeMB) * common.BytesPerMB

	for _, chunkFile := range chunkFiles {
		chunk, err := os.Open(chunkFile)
		if err != nil {
			// Clean up on error
			errRemove := removeFile(finalPath)
			if errRemove != nil {
				return apperror.NewAppError(fmt.Errorf("%w: %w", errRemove, err), "storage.LocalStorage.FinalizeChunkedUpload:OpenChunk:removeFile").WithMetadata("final_path", finalPath).WithMetadata("chunk_file", chunkFile)
			}

			return apperror.NewAppError(err, "storage.LocalStorage.FinalizeChunkedUpload:OpenChunk").WithMetadata("chunk_file", chunkFile)
		}

		written, err := io.Copy(finalFile, chunk)
		chunk.Close()

		if err != nil {
			// Clean up on error
			errRemove := removeFile(finalPath)
			if errRemove != nil {
				return apperror.NewAppError(fmt.Errorf("%w: %w", errRemove, err), "storage.LocalStorage.FinalizeChunkedUpload:CopyChunk:removeFile").WithMetadata("final_path", finalPath).WithMetadata("chunk_file", chunkFile)
			}

			return apperror.NewAppError(err, "storage.LocalStorage.FinalizeChunkedUpload:CopyChunk").WithMetadata("chunk_file", chunkFile)
		}

		// Check size limit
		if totalSize += written; totalSize > maxSize {
			// Clean up on error
			errRemove := removeFile(finalPath)
			if errRemove != nil {
				return apperror.NewAppError(fmt.Errorf("%w: %w", errRemove, err), "storage.LocalStorage.FinalizeChunkedUpload:SizeExceeded:removeFile").WithMetadata("final_path", finalPath).WithMetadata("chunk_file", chunkFile)
			}

			return apperror.NewAppError(apperror.ErrCommonInvalidValue, "storage.LocalStorage.FinalizeChunkedUpload:SizeExceeded").WithMetadata("total_size_mb", totalSize/common.BytesPerMB).WithMetadata("max_size_mb", maxSize/common.BytesPerMB)
		}
	}

	// Clean up the chunks directory
	err = s.CleanupChunks(ctx, uploadID, ownerID)
	if err != nil {
		// Just log the error, we don't want to fail the upload
		s.app.Logger.Error().Err(err).Str("upload_id", uploadID).Str("owner_id", ownerID).Msg("Failed to cleanup chunks after successful upload")
	}

	return nil
}

func (s *LocalStorage) CleanupChunks(ctx context.Context, uploadID string, ownerID string) error {
	chunksPath := getChunksDirPath(s.baseDir, ownerID, uploadID)

	// Remove the entire chunks directory for this upload
	err := os.RemoveAll(chunksPath)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return apperror.NewAppError(err, "storage.LocalStorage.CleanupChunks:RemoveAll").WithMetadata("chunks_path", chunksPath)
	}

	return nil
}

// CleanupOrphanedChunks removes chunk directories older than the specified duration
func (s *LocalStorage) CleanupOrphanedChunks(ctx context.Context, maxAge time.Duration) error {
	chunksBasePath := filepath.Join(s.baseDir, chunksDir)

	// Walk through all owner directories
	err := filepath.Walk(chunksBasePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip if not a directory or if it's the base chunks directory
		if !info.IsDir() || path == chunksBasePath {
			return nil
		}

		// Check if this is an upload directory (contains "upload_" prefix)
		if strings.Contains(info.Name(), "upload_") && time.Since(info.ModTime()) > maxAge {
			s.app.Logger.Info().Str("path", path).Msg("Cleaning up orphaned chunks")
			return os.RemoveAll(path)
		}

		return nil
	})

	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return apperror.NewAppError(err, "storage.LocalStorage.CleanupOrphanedChunks:Walk").WithMetadata("chunks_base_path", chunksBasePath)
	}

	return nil
}

func (s *LocalStorage) DeleteFile(ctx context.Context, name, ownerID string) error {
	deletePath := getFilePath(getOwnerDirPath(s.baseDir, ownerID), name)
	return removeFile(deletePath)
}

func (s *LocalStorage) DeleteChunk(ctx context.Context, uploadID string, chunkIndex int64, ownerID string) error {
	chunkPath := getChunkPath(getChunksDirPath(s.baseDir, ownerID, uploadID), chunkIndex)
	return removeFile(chunkPath)
}

func removeFile(path string) error {
	if err := os.Remove(path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonNoData, err), "storageremoveFile:Remove.NotExist").WithMetadata("remove_path", path)
		}

		return apperror.NewAppError(err, "storageremoveFile:Remove").WithMetadata("remove_path", path)
	}
	return nil
}

func (s *LocalStorage) OpenFile(ctx context.Context, name, ownerID string) (io.ReadSeekCloser, error) {
	openPath := getFilePath(getOwnerDirPath(s.baseDir, ownerID), name)
	f, err := os.Open(openPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonNoData, err), "storage.LocalStorage.OpenFile:Open").WithMetadata("open_path", openPath)
		}

		return nil, apperror.NewAppError(err, "storage.LocalStorage.OpenFile:Open").WithMetadata("open_path", openPath)
	}
	return f, nil
}

func getOwnerDirPath(baseDir string, ownerID string) string {
	return filepath.Join(baseDir, ownerID)
}

func getFilePath(ownerDir string, name string) string {
	return filepath.Join(ownerDir, name)
}

func getChunksDirPath(baseDir string, ownerID string, uploadID string) string {
	return filepath.Join(baseDir, chunksDir, ownerID, uploadID)
}

func getChunkPath(chunksDir string, chunkIndex int64) string {
	return filepath.Join(chunksDir, fmt.Sprintf("chunk_%d", chunkIndex))
}
