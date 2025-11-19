package dtos

import "encoding/base64"

type CreateProfileRes struct {
	ID       string `json:"id" copier:"must,nopanic"`
	Email    string `json:"email" copier:"must,nopanic"`
	FullName string `json:"fullName" copier:"must,nopanic"`
}

type GetProfileRes struct {
	ID           string `json:"id" copier:"must,nopanic"`
	Email        string `json:"email" copier:"must,nopanic"`
	FullName     string `json:"fullName" copier:"must,nopanic"`
	AvatarBase64 string `json:"avatarBase64,omitempty"`
}

func (r *GetProfileRes) SetAvatarBase64(avatar []byte) {
	if len(avatar) == 0 {
		return
	}

	r.AvatarBase64 = base64.StdEncoding.EncodeToString(avatar)
}

type StorageUsageRes struct {
	UsedBytes  int64 `json:"usedBytes"`
	QuotaBytes int64 `json:"quotaBytes"`
	UsedMB     int64 `json:"usedMB"`
	QuotaMB    int64 `json:"quotaMB"`
}
