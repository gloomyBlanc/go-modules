package pnm

import (
	"image"
	"io"
)

func Encode(w io.Writer, img image.Image) error {
	// ヘッダ情報の登録

	switch img.(type) {
	case *image.Gray, *image.Gray16:
	case *image.NRGBA, *image.NRGBA64:
	}
	// ラスタデータの出力
	return nil
}

type pnmHeader struct {
	magicNumber   string
	width, height int
	maxValue      int
}
