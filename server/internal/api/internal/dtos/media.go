package dtos

import "time"

type GetFilesInfoRes struct {
	ID        int64     `json:"id"`
	FolderID  *int64    `json:"folderId,omitempty"`
	OwnerID   int64     `json:"ownerId"`
	Name      string    `json:"name"`
	SizeBytes int64     `json:"sizeBytes"`
	MimeType  string    `json:"mimeType"`
	Extension *string   `json:"extension,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
