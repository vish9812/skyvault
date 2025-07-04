package utils

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"golang.org/x/image/draw"
)

var (
	ErrUnsupportedImageFormat = errors.New("unsupported image format")
)

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
	draw.CatmullRom.Scale(resized, resized.Bounds(), img, img.Bounds(), draw.Over, nil)

	buf := bytes.NewBuffer(nil)
	if format == "jpeg" {
		err = jpeg.Encode(buf, resized, &jpeg.Options{Quality: 100})
	} else {
		err = png.Encode(buf, resized)
	}

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// SampleImage generates a small test image
func SampleImage(buf *bytes.Buffer) error {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	err := png.Encode(buf, img)
	return err
}
