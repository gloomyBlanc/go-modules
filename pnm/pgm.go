package pnm

import (
	"errors"
	"image"
	"image/color"
	"strconv"
)

var (
	errBadPGMSample = errors.New("pnm: PGM画像のサンプル値が不正です")
)

func (d *pnmDecoder) pgmReadRaster() (image.Image, error) {
	var (
		i, j        int
		b           byte
		pixel       int
		readBytes   []byte
		err         error
		overFF      bool
		enSampleEnd bool
	)
	overFF = (d.maxValue < 256)
	img := image.NewGray16(image.Rect(0, 0, d.width, d.height))

	enSampleEnd = false
	for i = 0; i < d.height; i++ {
		for j = 0; j < d.width; {
			b, err = d.reader.ReadByte()
			if err != nil {
				return nil, errBadPGMSample
			}
			switch d.magicNumber {
			case "P2":
				if enSampleEnd {
					if isWhiteSpece(b) {
						pixel, err = strconv.Atoi(string(readBytes))
						if err != nil {
							return nil, errBadPGMSample
						}
						img.SetGray16(j, i,
							color.Gray16{uint16(pixel * 65536.0 / d.maxValue)},
						)
						readBytes = []byte{}
						enSampleEnd = false
						j += 1
					} else {
						readBytes = append(readBytes, b)
					}
				} else if !isWhiteSpece(b) {
					readBytes = append(readBytes, b)
					enSampleEnd = true
				}
			case "P5":
				if overFF {
					if enSampleEnd {
						pixel = (pixel << 8) | int(b-'0')
						img.SetGray16(j, i,
							color.Gray16{uint16(pixel * 65536.0 / d.maxValue)},
						)
						enSampleEnd = false
						j += 1
					} else {
						pixel = int(b - '0')
						enSampleEnd = true
					}
				} else {
					pixel = int(b - '0')
					img.SetGray16(j, i,
						color.Gray16{uint16(pixel * 65536.0 / d.maxValue)},
					)
					j += 1
				}
			}
		}
	}
	return img, nil
}
