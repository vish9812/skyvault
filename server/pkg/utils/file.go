package utils

import "strings"

func CleanFileName(name string) string {
	// Remove any path separators to prevent directory traversal
	name = strings.ReplaceAll(name, "/", "")
	name = strings.ReplaceAll(name, "\\", "")

	// Remove any null bytes that could be used to truncate strings
	name = strings.ReplaceAll(name, "\x00", "")

	// Trim spaces from start/end
	name = strings.TrimSpace(name)

	return name
}
