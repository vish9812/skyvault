package media

import "time"

type Folder struct {
	ID             int64
	Name           string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	TrashedAt      *time.Time
	ParentFolderID *int64
}
