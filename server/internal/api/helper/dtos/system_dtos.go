package dtos

type SystemConfigDTO struct {
	MaxDirectUploadSizeMB int64 `json:"maxDirectUploadSizeMB"`
	MaxChunkSizeMB        int64 `json:"maxChunkSizeMB"`
}