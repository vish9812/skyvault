package dtos

type SignUp struct {
	Token   string         `json:"token" copier:"must,nopanic"`
	Profile *GetProfileRes `json:"profile" copier:"must,nopanic"`
}
