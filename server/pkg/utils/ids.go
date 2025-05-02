package utils

import (
	"github.com/google/uuid"
)

func UUID() string {
	return uuid.New().String()
}

func Ptr[T any](v T) *T {
	return &v
}
