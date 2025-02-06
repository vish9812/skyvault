package utils

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"strings"

	"golang.org/x/image/draw"
)

var (
	ErrUnsupportedImageFormat = errors.New("unsupported image format")
)

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

// ScaleDownImageTo resizes an image of type jpeg or png to the given width and height
//
// App Errors:
// - ErrUnsupportedImageFormat
func ScaleDownImageTo(format string, reader io.ReadSeeker, maxWidth, maxHeight int) ([]byte, error) {
	if format != "jpeg" && format != "png" {
		return nil, ErrUnsupportedImageFormat
	}

	reader.Seek(0, io.SeekStart)

	var img image.Image
	var err error

	// img, _, err = image.Decode(reader)
	// if err != nil {
	// 	return nil, err
	// }
	if format == "jpeg" {
		img, err = jpeg.Decode(reader)
		if err != nil {
			return nil, err
		}
	} else {
		img, err = png.Decode(reader)
		if err != nil {
			return nil, err
		}
	}

	resized := image.NewRGBA(image.Rect(0, 0, maxWidth, maxHeight))
	draw.NearestNeighbor.Scale(resized, resized.Bounds(), img, img.Bounds(), draw.Over, nil)

	buf := bytes.NewBuffer(nil)
	if format == "jpeg" {
		err = jpeg.Encode(buf, resized, nil)
	} else {
		err = png.Encode(buf, resized)
	}

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
