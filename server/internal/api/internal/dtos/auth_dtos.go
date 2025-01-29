package dtos

type SignUpReq struct {
	Email          string  `json:"email"`
	FullName       string  `json:"fullName"`
	Password       *string `json:"password"`
	Provider       string  `json:"provider"`
}

type SignUpRes struct {
	Token   string         `json:"token" copier:"must,nopanic"`
	Profile *GetProfileRes `json:"profile" copier:"must,nopanic"`
}

type SignInReq struct {
	Provider       string  `json:"provider"`
	ProviderUserID string  `json:"providerUserId"`
	Password       *string `json:"password"`
}

type SignInRes struct {
	Token   string         `json:"token" copier:"must,nopanic"`
	Profile *GetProfileRes `json:"profile" copier:"must,nopanic"`
}
