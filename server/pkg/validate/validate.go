package validate

import (
	"fmt"
	"net/mail"
	"regexp"
	"skyvault/pkg/apperror"
	"strings"

	"github.com/google/uuid"
)

const (
	MaxLen = 255
)

var uploadIDRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

func Email(email string) (string, error) {
	em := strings.TrimSpace(strings.ToLower(email))
	if em == "" || len(em) > MaxLen {
		return "", apperror.ErrCommonInvalidValue
	}

	mailObj, err := mail.ParseAddress(email)
	if err != nil {
		return "", fmt.Errorf("%w: %w", apperror.ErrCommonInvalidValue, err)
	}
	return mailObj.Address, nil
}

func FileName(name string) (string, error) {
	// Remove any path separators to prevent directory traversal
	name = strings.ReplaceAll(name, "/", "")
	name = strings.ReplaceAll(name, "\\", "")

	// Remove any null bytes that could be used to truncate strings
	name = strings.ReplaceAll(name, "\x00", "")

	return Name(name)
}

func Name(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" || len(name) > MaxLen {
		return "", apperror.ErrCommonInvalidValue
	}

	return name, nil
}

func PasswordLen(pwd string) (string, error) {
	pwd = strings.TrimSpace(pwd)
	if len(pwd) < 4 || len(pwd) > MaxLen {
		return "", apperror.ErrCommonInvalidValue
	}

	return pwd, nil
}

func UUID(uuidStr string) bool {
	return uuid.Validate(uuidStr) == nil
}

func UUIDs(uuidStrs []string) ([]string, []string) {
	validUUIDs := make([]string, 0, len(uuidStrs))
	invalidUUIDs := make([]string, 0)
	for _, uuidStr := range uuidStrs {
		if UUID(uuidStr) {
			validUUIDs = append(validUUIDs, uuidStr)
		} else {
			invalidUUIDs = append(invalidUUIDs, uuidStr)
		}
	}
	return validUUIDs, invalidUUIDs
}

// UploadID validates that an upload ID contains only safe characters
// to prevent directory traversal attacks. Upload IDs should only contain
// alphanumeric characters, underscores, and hyphens.
func UploadID(uploadID string) bool {
	if uploadID == "" || len(uploadID) > MaxLen {
		return false
	}

	return uploadIDRegex.MatchString(uploadID)
}
