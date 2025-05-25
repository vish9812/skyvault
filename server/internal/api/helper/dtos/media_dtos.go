package dtos

import (
	"encoding/base64"
	"skyvault/pkg/paging"
	"time"
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
	PreviewBase64 *string   `json:"preview"`
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
