package dtos

import "time"

type GetFileInfoRes struct {
	ID            int64     `json:"id" copier:"must,nopanic"`
	OwnerID       int64     `json:"ownerId" copier:"must,nopanic"`
	FolderID      *int64    `json:"folderId,omitempty"`
	Name          string    `json:"name" copier:"must,nopanic"`
	GeneratedName string    `json:"generatedName" copier:"must,nopanic"`
	Size          int64     `json:"size" copier:"must,nopanic"`
	Extension     *string   `json:"extension,omitempty"`
	MimeType      string    `json:"mimeType" copier:"must,nopanic"`
	CreatedAt     time.Time `json:"createdAt" copier:"must,nopanic"`
	UpdatedAt     time.Time `json:"updatedAt" copier:"must,nopanic"`
}

type GetFilesInfoRes struct {
	Infos []GetFileInfoRes `json:"infos" copier:"must,nopanic"`
}

type GetFolderInfoRes struct {
	ID             int64     `json:"id" copier:"must,nopanic"`
	OwnerID        int64     `json:"ownerId" copier:"must,nopanic"`
	ParentFolderID *int64    `json:"parentFolderId,omitempty"`
	Name           string    `json:"name" copier:"must,nopanic"`
	CreatedAt      time.Time `json:"createdAt" copier:"must,nopanic"`
	UpdatedAt      time.Time `json:"updatedAt" copier:"must,nopanic"`
}

type CreateFolderReq struct {
	Name           string `json:"name"`
	ParentFolderID *int64 `json:"parentFolderId,omitempty"`
}

type GetFoldersInfoRes struct {
	Folders []GetFolderInfoRes `json:"folders" copier:"must,nopanic"`
}
