package dtos

type BaseInfo struct {
	ID   string `json:"id" copier:"must,nopanic"`
	Name string `json:"name" copier:"must,nopanic"`
}
