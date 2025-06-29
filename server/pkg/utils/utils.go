package utils

import "github.com/google/uuid"

func Ptr[T any](v T) *T {
	return &v
}

// ID generates a new time-ordered UUID(v7) string.
func ID() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

// ParseID parses a UUID string and returns uuid.UUID.
func ParseID(id string) (uuid.UUID, error) {
	if err := uuid.Validate(id); err != nil {
		return uuid.UUID{}, err
	}

	return uuid.Parse(id)
}

// ParseValidID assumes that id is already validated. So, it will return a valid uuid.UUID.
func ParseValidID(id string) uuid.UUID {
	return uuid.MustParse(id)
}

// UUID generates a new UUID(v4) string.
func UUID() string {
	return uuid.New().String()
}
