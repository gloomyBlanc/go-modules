package pnm

import (
	"image"
	"io"
)

func Encode(w io.Writer, img image.Image) error {
	return nil
}

type pnmHeader struct {
	magicNumber   string
	width, height int
	maxValue      int
}
