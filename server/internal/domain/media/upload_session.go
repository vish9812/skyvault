package media

import (
	"fmt"
	"skyvault/pkg/apperror"
	"skyvault/pkg/common"
	"skyvault/pkg/utils"
	"time"
)

const (
	// UploadSessionTTLHours is the number of hours an upload session remains valid
	UploadSessionTTLHours = 24
)

// UploadSession represents a tracked chunked upload session.
// It reserves storage quota upfront to prevent quota bypass attacks.
type UploadSession struct {
	ID             string
	OwnerID        string
	FolderID       *string // null if uploading to root folder
	FileName       string
	FileSize       int64 // Total size of the final file in bytes
	MimeType       string
	TotalChunks    int64
	QuotaAllocated int64 // Amount of quota reserved for this upload
	UploadedBytes  int64 // Actual bytes uploaded so far (cumulative across all chunks)
	CreatedAt      time.Time
	ExpiresAt      time.Time
}

// App Errors:
// - ErrCommonNoAccess
// - ErrCommonInvalidValue
func NewUploadSession(ownerID string, parentFolder *FolderInfo, fileName string, fileSize int64, mimeType string, totalChunks int64) (*UploadSession, error) {
	var folderID *string
	if parentFolder != nil {
		if err := parentFolder.ValidateAccess(ownerID); err != nil {
			return nil, apperror.NewAppError(err, "media.NewUploadSession:ValidateParentAccess")
		}
		folderID = &parentFolder.ID
	}

	// Validate file size
	if fileSize <= 0 {
		return nil, apperror.NewAppError(fmt.Errorf("%w: file size must be positive", apperror.ErrCommonInvalidValue), "media.NewUploadSession:InvalidFileSize").
			WithMetadata("file_size", fileSize)
	}

	// Validate total chunks
	if totalChunks <= 0 {
		return nil, apperror.NewAppError(fmt.Errorf("%w: total chunks must be positive", apperror.ErrCommonInvalidValue), "media.NewUploadSession:InvalidTotalChunks").
			WithMetadata("total_chunks", totalChunks)
	}

	// Validate chunks don't exceed maximum allowed based on direct upload size limit
	// maxTotalChunks = MaxDirectUploadSizeMB / MaxChunkSizeMB
	maxTotalChunks := int64((MaxDirectUploadSizeMB * common.BytesPerMB) / (MaxChunkSizeMB * common.BytesPerMB))
	if totalChunks > maxTotalChunks {
		return nil, apperror.NewAppError(fmt.Errorf("%w: too many chunks", apperror.ErrCommonInvalidValue), "media.NewUploadSession:TooManyChunks").
			WithMetadata("total_chunks", totalChunks).
			WithMetadata("max_total_chunks", maxTotalChunks)
	}

	// Validate file name
	if fileName == "" {
		return nil, apperror.NewAppError(fmt.Errorf("%w: file name cannot be empty", apperror.ErrCommonInvalidValue), "media.NewUploadSession:EmptyFileName")
	}

	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	id, err := utils.ID()
	if err != nil {
		return nil, apperror.NewAppError(err, "media.NewUploadSession:ID")
	}

	now := time.Now()
	expiresAt := now.Add(UploadSessionTTLHours * time.Hour)

	return &UploadSession{
		ID:             id,
		OwnerID:        ownerID,
		FolderID:       folderID,
		FileName:       fileName,
		FileSize:       fileSize,
		MimeType:       mimeType,
		TotalChunks:    totalChunks,
		QuotaAllocated: fileSize, // Reserve the full file size upfront
		UploadedBytes:  0,        // No bytes uploaded yet
		CreatedAt:      now,
		ExpiresAt:      expiresAt,
	}, nil
}

// IsExpired checks if the upload session has expired
func (s *UploadSession) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// App Errors:
// - ErrCommonNoAccess
func (s *UploadSession) ValidateAccess(ownerID string) error {
	if s.OwnerID != ownerID {
		return apperror.NewAppError(apperror.ErrCommonNoAccess, "media.UploadSession.ValidateAccess:OwnerMismatch").
			WithMetadata("session_owner_id", s.OwnerID).
			WithMetadata("requester_id", ownerID)
	}
	return nil
}

// App Errors:
// - ErrCommonInvalidValue
func (s *UploadSession) ValidateExpiration() error {
	if s.IsExpired() {
		return apperror.NewAppError(fmt.Errorf("%w: upload session expired", apperror.ErrCommonInvalidValue), "media.UploadSession.ValidateExpiration:Expired").
			WithMetadata("session_id", s.ID).
			WithMetadata("expires_at", s.ExpiresAt).
			WithMetadata("current_time", time.Now())
	}
	return nil
}
