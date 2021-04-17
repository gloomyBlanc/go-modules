package pnm

import (
	"errors"
	"image"
	"image/color"
)

var (
	errBadPBMSample = errors.New("pnm: PBM画像のサンプル値が不正です")
)

func (d *pnmDecoder) pbmReadRaster() (image.Image, error) {
	var (
		i, j, k int
		b       byte
		err     error
	)

	img := image.NewGray(image.Rect(0, 0, d.h.width, d.h.height))
	for i = 0; i < d.h.height; i++ {
		for j = 0; j < d.h.width; {
			b, err = d.reader.ReadByte()
			if err != nil {
				return nil, errBadPBMSample
			}
			switch d.h.magicNumber {
			case "P1":
				if !isWhiteSpece(b) {
					img.SetGray(j, i, color.Gray{255 * (b - '0')})
					j += 1
				}
			case "P4":
				for k = 0; k < 8; k++ {
					img.SetGray(j+k, i, color.Gray{255 * ((b >> (7 - k)) & 1)})
				}
				j += 8
			}
		}
	}
	return img, nil
}
