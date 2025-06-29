package dtos

type SystemConfigDTO struct {
	MaxUploadSizeMB       int64 `json:"maxUploadSizeMB"`
	MaxDirectUploadSizeMB int64 `json:"maxDirectUploadSizeMB"`
	MaxChunkSizeMB        int64 `json:"maxChunkSizeMB"`
}