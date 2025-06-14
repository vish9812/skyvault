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
	"sort"
	"strconv"
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
	err := os.MkdirAll(baseDir, 0750)
	if err != nil {
		app.Logger.Fatal().Err(err).Str("base_dir", baseDir).Msg("Failed to create base directory for local storage")
	}

	// Ensure chunks directory exists
	chunksPath := filepath.Join(baseDir, chunksDir)
	err = os.MkdirAll(chunksPath, 0750)
	if err != nil {
		app.Logger.Fatal().Err(err).Str("chunks_dir", chunksPath).Msg("Failed to create chunks directory for local storage")
	}

	return &LocalStorage{app: app, baseDir: baseDir}
}

func (s *LocalStorage) SaveFile(ctx context.Context, file io.ReadSeeker, name, ownerID string) error {
	// Create the owner directory
	ownerDir := getOwnerDir(s.baseDir, ownerID)
	err := os.MkdirAll(ownerDir, 0750)
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
	// f, err := os.Create(savePath)
	f, err := os.OpenFile(savePath, os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		return apperror.NewAppError(err, "LocalStorage.SaveFile:Create").WithMetadata("save_path", savePath)
	}
	defer f.Close()

	// Write the file
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return apperror.NewAppError(err, "LocalStorage.SaveFile:Seek").WithMetadata("save_path", savePath)
	}

	// Copy the file with an extra MB to check if file is greater than max size.
	written, err := io.CopyN(f, file, media.MaxDirectUploadSizeMB)
	if err != nil && !errors.Is(err, io.EOF) {
		// Clean up the partially written file if copying fails
		if delErr := s.DeleteFile(ctx, name, ownerID); delErr != nil {
			if !errors.Is(delErr, apperror.ErrCommonNoData) {
				return apperror.NewAppError(fmt.Errorf("%w: %w", delErr, err), "LocalStorage.SaveFile:CopyN:DeleteFile")
			}
		}
		return apperror.NewAppError(err, "LocalStorage.SaveFile:CopyN").WithMetadata("save_path", savePath)
	}

	if written > media.MaxDirectUploadSizeMB {
		var sizeErr error = apperror.ErrMediaFileSizeLimitExceeded
		// Clean up the written file if it is greater than max size
		if delErr := s.DeleteFile(ctx, name, ownerID); delErr != nil {
			sizeErr = apperror.NewAppError(fmt.Errorf("%w: %w", delErr, sizeErr), "LocalStorage.SaveFile:MaxSizeExceeded:DeleteFile")
		}
		return apperror.NewAppError(sizeErr, "LocalStorage.SaveFile:MaxSizeExceeded").WithMetadata("save_path", savePath).WithMetadata("written_mb", written/media.BytesPerMB).WithMetadata("max_direct_upload_size_mb", media.MaxDirectUploadSizeMB/media.BytesPerMB)
	}

	return nil
}

func (s *LocalStorage) SaveChunk(ctx context.Context, reader io.Reader, uploadID string, chunkIndex int64, ownerID string) error {
	// Create chunks directory for this upload
	chunksPath := getChunksDir(s.baseDir, ownerID, uploadID)
	err := os.MkdirAll(chunksPath, 0750)
	if err != nil {
		return apperror.NewAppError(err, "LocalStorage.SaveChunk:MkdirAll").WithMetadata("chunks_path", chunksPath)
	}

	chunkPath := getChunkPath(chunksPath, chunkIndex)

	// Create the chunk file
	f, err := os.OpenFile(chunkPath, os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		return apperror.NewAppError(err, "LocalStorage.SaveChunk:Create").WithMetadata("chunk_path", chunkPath)
	}
	defer f.Close()

	// Copy chunk data with max chunk size limit for safety
	_, err = io.CopyN(f, reader, media.MaxChunkSizeMB)
	if err != nil && !errors.Is(err, io.EOF) {
		// Clean up the chunk file if copying fails
		os.Remove(chunkPath)
		return apperror.NewAppError(err, "LocalStorage.SaveChunk:CopyN").WithMetadata("chunk_path", chunkPath)
	}

	return nil
}

func (s *LocalStorage) FinalizeChunkedUpload(ctx context.Context, uploadID string, fileName string, ownerID string) error {
	chunksPath := getChunksDir(s.baseDir, ownerID, uploadID)

	// List all chunk files
	chunkFiles, err := filepath.Glob(filepath.Join(chunksPath, "chunk_*"))
	if err != nil {
		return apperror.NewAppError(err, "LocalStorage.FinalizeChunkedUpload:Glob").WithMetadata("chunks_path", chunksPath)
	}

	if len(chunkFiles) == 0 {
		return apperror.NewAppError(apperror.ErrCommonNoData, "LocalStorage.FinalizeChunkedUpload:NoChunks").WithMetadata("chunks_path", chunksPath)
	}

	// Sort chunk files by index
	sort.Slice(chunkFiles, func(i, j int) bool {
		indexI := extractChunkIndex(chunkFiles[i])
		indexJ := extractChunkIndex(chunkFiles[j])
		return indexI < indexJ
	})

	// Create the owner directory
	ownerDir := getOwnerDir(s.baseDir, ownerID)
	err = os.MkdirAll(ownerDir, 0750)
	if err != nil {
		return apperror.NewAppError(err, "LocalStorage.FinalizeChunkedUpload:MkdirAll").WithMetadata("owner_dir", ownerDir)
	}

	finalPath := getFilePath(ownerDir, fileName)

	// Check if the file already exists
	if _, err := os.Stat(finalPath); err == nil {
		return apperror.NewAppError(apperror.ErrCommonDuplicateData, "LocalStorage.FinalizeChunkedUpload:Stat").WithMetadata("final_path", finalPath)
	}

	// Create the final file
	finalFile, err := os.OpenFile(finalPath, os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		return apperror.NewAppError(err, "LocalStorage.FinalizeChunkedUpload:Create").WithMetadata("final_path", finalPath)
	}
	defer finalFile.Close()

	// Combine all chunks
	var totalSize int64
	maxSize := s.app.Config.Media.MaxSizeMB * media.BytesPerMB

	for _, chunkFile := range chunkFiles {
		chunk, err := os.Open(chunkFile)
		if err != nil {
			// Clean up on error
			os.Remove(finalPath)
			return apperror.NewAppError(err, "LocalStorage.FinalizeChunkedUpload:OpenChunk").WithMetadata("chunk_file", chunkFile)
		}

		written, err := io.Copy(finalFile, chunk)
		chunk.Close()

		if err != nil {
			// Clean up on error
			os.Remove(finalPath)
			return apperror.NewAppError(err, "LocalStorage.FinalizeChunkedUpload:CopyChunk").WithMetadata("chunk_file", chunkFile)
		}

		totalSize += written

		// Check size limit
		if totalSize > maxSize {
			// Clean up on error
			os.Remove(finalPath)
			return apperror.NewAppError(apperror.ErrMediaFileSizeLimitExceeded, "LocalStorage.FinalizeChunkedUpload:SizeExceeded").WithMetadata("total_size_mb", totalSize/media.BytesPerMB).WithMetadata("max_size_mb", maxSize/media.BytesPerMB)
		}
	}

	return nil
}

func (s *LocalStorage) CleanupChunks(ctx context.Context, uploadID string, ownerID string) error {
	chunksPath := getChunksDir(s.baseDir, ownerID, uploadID)

	// Remove the entire chunks directory for this upload
	err := os.RemoveAll(chunksPath)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return apperror.NewAppError(err, "LocalStorage.CleanupChunks:RemoveAll").WithMetadata("chunks_path", chunksPath)
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
		return apperror.NewAppError(err, "LocalStorage.CleanupOrphanedChunks:Walk").WithMetadata("chunks_base_path", chunksBasePath)
	}

	return nil
}

func (s *LocalStorage) DeleteFile(ctx context.Context, name, ownerID string) error {
	deletePath := getFilePath(getOwnerDir(s.baseDir, ownerID), name)
	if err := os.Remove(deletePath); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return apperror.NewAppError(fmt.Errorf("%w: %w", apperror.ErrCommonNoData, err), "LocalStorage.DeleteFile:Remove").WithMetadata("delete_path", deletePath)
		}

		return apperror.NewAppError(err, "LocalStorage.DeleteFile:Remove").WithMetadata("delete_path", deletePath)
	}
	return nil
}

func (s *LocalStorage) OpenFile(ctx context.Context, name, ownerID string) (io.ReadSeekCloser, error) {
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

func getOwnerDir(baseDir string, ownerID string) string {
	return filepath.Join(baseDir, ownerID)
}

func getFilePath(ownerDir string, name string) string {
	return filepath.Join(ownerDir, name)
}

func getChunksDir(baseDir string, ownerID string, uploadID string) string {
	return filepath.Join(baseDir, chunksDir, ownerID, uploadID)
}

func getChunkPath(chunksDir string, chunkIndex int64) string {
	return filepath.Join(chunksDir, fmt.Sprintf("chunk_%d", chunkIndex))
}

func extractChunkIndex(chunkPath string) int64 {
	fileName := filepath.Base(chunkPath)
	parts := strings.Split(fileName, "_")
	if len(parts) != 2 {
		return 0
	}

	index, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0
	}

	return index
}
