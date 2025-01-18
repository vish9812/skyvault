package utils

import (
	"path/filepath"
	"strings"
)

// GetFileNameAndExtension returns the file name(without extension) and extension(with dot) from a given file name
func GetFileNameAndExtension(fileName string) (nameWithoutExt string, extensionWithDot string) {
	extensionWithDot = filepath.Ext(fileName)
	nameWithoutExt = strings.TrimSuffix(fileName, extensionWithDot)
	return nameWithoutExt, extensionWithDot
}
