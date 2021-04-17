package pnm

import (
	"errors"
	"image"
	"image/color"
	"strconv"
)

var (
	errBadPPMSample = errors.New("pnm:PPM画像のサンプル値が不正です")
)

func (d *pnmDecoder) ppmReadRaster() (image.Image, error) {
	var (
		i, j, k     int
		b           byte
		pixel       [3]int
		readBytes   []byte
		err         error
		overFF      bool
		enSampleEnd bool
	)
	overFF = (d.maxValue < 256)
	img := image.NewRGBA64(image.Rect(0, 0, d.width, d.height))

	enSampleEnd = false
	for i = 0; i < d.height; i++ {
		for j = 0; j < d.width; j++ {
			for k = 0; k < 3; {
				b, err = d.reader.ReadByte()
				if err != nil {
					return nil, errBadPGMSample
				}
				switch d.magicNumber {
				case "P3":
					if enSampleEnd {
						if isWhiteSpece(b) {
							pixel[k], err = strconv.Atoi(string(readBytes))
							if err != nil {
								return nil, errBadPGMSample
							}
							readBytes = []byte{}
							enSampleEnd = false
							k += 1
						} else {
							readBytes = append(readBytes, b)
						}
					} else if !isWhiteSpece(b) {
						readBytes = append(readBytes, b)
						enSampleEnd = true
					}
				case "P6":
					if overFF {
						if enSampleEnd {
							pixel[k] = (pixel[k] << 8) | int(b-'0')
							enSampleEnd = false
							k += 1
						} else {
							pixel[k] = int(b - '0')
							enSampleEnd = true
						}
					} else {
						pixel[k] = int(b - '0')
						k += 1
					}
				}
			}
			// pixel値の代入
			img.SetRGBA64(j, i,
				color.RGBA64{
					uint16(pixel[0]),
					uint16(pixel[1]),
					uint16(pixel[2]),
					0xFFFF,
				},
			)
			pixel = [3]int{}
		}
	}
	return img, nil
}
