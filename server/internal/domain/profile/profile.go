package profile

type Profile struct {
	ID       int64
	Email    string
	FullName string
}

func NewProfile() *Profile {
	return &Profile{}
}
