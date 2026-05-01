package dtos

import (
	"encoding/base64"
	"fmt"
	"skyvault/pkg/apperror"
	"skyvault/pkg/paging"
	"time"
)

// Boundary sanity caps for upload-session inputs. The domain sanitizer enforces business rules;
// these guard the API edge so an obviously malformed payload is rejected before any downstream work.
const (
	maxUploadSessionChunks  = 100_000 // ~1TB at 10MB/chunk; well past any realistic upload
	maxUploadSessionMimeLen = 255     // RFC 6838 caps type/subtype at 127 each, 255 is generous
)

type GetFileInfo struct {
	ID            string    `json:"id" copier:"must,nopanic"`
	OwnerID       string    `json:"ownerId" copier:"must,nopanic"`
	FolderID      *string   `json:"folderId,omitempty"`
	Name          string    `json:"name" copier:"must,nopanic"`
	Size          int64     `json:"size" copier:"must,nopanic"`
	Extension     *string   `json:"extension,omitempty"`
	MimeType      string    `json:"mimeType" copier:"must,nopanic"`
	Category      string    `json:"category" copier:"must,nopanic"`
	PreviewBase64 *string   `json:"previewBase64"`
	CreatedAt     time.Time `json:"createdAt" copier:"must,nopanic"`
	UpdatedAt     time.Time `json:"updatedAt" copier:"must,nopanic"`
}

func (r *GetFileInfo) Preview(preview []byte) {
	if len(preview) > 0 {
		previewBase64 := base64.StdEncoding.EncodeToString(preview)
		r.PreviewBase64 = &previewBase64
	}
}

type GetFolderContent struct {
	FilePage   *paging.Page[*GetFileInfo]   `json:"filePage" copier:"must,nopanic"`
	FolderPage *paging.Page[*GetFolderInfo] `json:"folderPage" copier:"must,nopanic"`
}

type GetFolderInfo struct {
	ID             string     `json:"id" copier:"must,nopanic"`
	OwnerID        string     `json:"ownerId" copier:"must,nopanic"`
	Name           string     `json:"name" copier:"must,nopanic"`
	ParentFolderID *string    `json:"parentFolderId,omitempty"`
	CreatedAt      time.Time  `json:"createdAt" copier:"must,nopanic"`
	UpdatedAt      time.Time  `json:"updatedAt" copier:"must,nopanic"`
	Ancestors      []BaseInfo `json:"ancestors" copier:"nopanic"`
}

type CreateUploadSessionRequest struct {
	FolderID    *string `json:"folderId,omitempty"`
	FileName    string  `json:"fileName"`
	FileSize    int64   `json:"fileSize"`
	MimeType    string  `json:"mimeType"`
	TotalChunks int64   `json:"totalChunks"`
}

// Validate performs boundary-level sanity checks on the request. The domain sanitizer
// re-checks file name (path-traversal stripping) and the positive-size invariants, so
// this only adds API-edge bounds (chunk count sanity, mime length).
func (r *CreateUploadSessionRequest) Validate() error {
	if r.FileSize <= 0 {
		return fmt.Errorf("%w: fileSize must be positive", apperror.ErrCommonInvalidValue)
	}
	if r.TotalChunks <= 0 || r.TotalChunks > maxUploadSessionChunks {
		return fmt.Errorf("%w: totalChunks out of range", apperror.ErrCommonInvalidValue)
	}
	if r.FileName == "" {
		return fmt.Errorf("%w: fileName is required", apperror.ErrCommonInvalidValue)
	}
	if len(r.MimeType) > maxUploadSessionMimeLen {
		return fmt.Errorf("%w: mimeType too long", apperror.ErrCommonInvalidValue)
	}
	return nil
}

type CreateUploadSessionResponse struct {
	UploadID    string    `json:"uploadId" copier:"must,nopanic"`
	FileName    string    `json:"fileName" copier:"must,nopanic"`
	FileSize    int64     `json:"fileSize" copier:"must,nopanic"`
	TotalChunks int64     `json:"totalChunks" copier:"must,nopanic"`
	ExpiresAt   time.Time `json:"expiresAt" copier:"must,nopanic"`
}
