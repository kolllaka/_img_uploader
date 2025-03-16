package img

import (
	"image"
	"image/png"
	"io"
	"os"

	_ "image/gif"
	_ "image/jpeg"

	_ "golang.org/x/image/webp"

	"github.com/nfnt/resize"
)

type Transform interface {
	Resize(file io.Reader, width, height uint) (image.Image, error)
	SavePng(path string, m image.Image) error
}

type transform struct {
}

func NewImage() Transform {
	return &transform{}
}

func (i *transform) Resize(file io.Reader, width, height uint) (image.Image, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	m := resize.Resize(width, height, img, resize.Lanczos3)

	return m, nil
}

func (i *transform) SavePng(path string, m image.Image) error {
	outFile, err := os.Create(path + ".png")
	if err != nil {
		return err
	}
	defer outFile.Close()

	if err := png.Encode(outFile, m); err != nil {
		return err
	}

	return nil
}
