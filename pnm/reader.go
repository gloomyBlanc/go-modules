package pnm

import (
	"bufio"
	"errors"
	"image"
	"image/color"
	"io"
	"strconv"
)

var (
	errBadHeader   = errors.New("pnm: ヘッダ情報が不正です")
	errBadMagicNum = errors.New("pnm: ファイル形式を正常に読み込めませんでした")
)

func init() {
	image.RegisterFormat("ppm", "P?", Decode, DecodeConfig)
	image.RegisterFormat("pgm", "P?", Decode, DecodeConfig)
	image.RegisterFormat("pbm", "P?", Decode, DecodeConfig)
	image.RegisterFormat("pnm", "P?", Decode, DecodeConfig)
}

func Decode(r io.Reader) (image.Image, error) {
	var d pnmDecoder
	img, err := d.decode(r, false)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func DecodeConfig(r io.Reader) (image.Config, error) {
	var d pnmDecoder
	_, err := d.decode(r, true)
	if err != nil {
		return image.Config{}, err
	}

	switch sortPNM(d.magicNumber) {
	case PBM, PGM:
		return image.Config{
			ColorModel: color.GrayModel,
			Width:      d.width,
			Height:     d.height,
		}, nil
	case PPM:
		return image.Config{
			ColorModel: color.RGBA64Model,
			Width:      d.width,
			Height:     d.height,
		}, nil
	}
	return image.Config{}, errBadMagicNum
}

type pnmDecoder struct {
	reader *bufio.Reader
	// ヘッダ情報
	magicNumber   string
	width, height int
	maxValue      int
}

// メソッド
func (d *pnmDecoder) decode(r io.Reader, isConfig bool) (image.Image, error) {
	d.reader = bufio.NewReader(r)
	err := d.decodeHeader()
	if err != nil {
		return nil, err
	}
	if !isConfig {
		switch sortPNM(d.magicNumber) {
		case PBM:
			return d.pbmReadRaster()
		case PGM:
		case PPM:
		}
	}
	return nil, nil
}

func (d *pnmDecoder) decodeHeader() error {
	var (
		i         int
		err       error
		b         byte
		isComment bool
		readBytes [4][]byte
	)
	// ヘッダ情報の読み込み
	isComment = false
	for i = 0; i < 4; {
		b, err = d.reader.ReadByte()
		if err != nil {
			return errBadHeader
		}
		if isComment {
			if b == '\n' {
				isComment = false
			}
		} else if isWhiteSpece(b) {
			i++
			if i == 3 && sortPNM(string(readBytes[0])) == PBM {
				i++
			}
		} else {
			readBytes[i] = append(readBytes[i], b)
		}
	}
	// メンバ変数に代入
	d.magicNumber = string(readBytes[0])
	d.width, err = strconv.Atoi(string(readBytes[1]))
	if err != nil {
		return errBadHeader
	}
	d.width, err = strconv.Atoi(string(readBytes[2]))
	if err != nil {
		return errBadHeader
	}
	return nil
}

// WhiteSpace指定文字か判別(blanks, TABs, CRs, LFs)
func isWhiteSpece(b byte) bool {
	var ws_list = []byte{
		' ', '\t', '\n', '\r', '\v',
	}
	for _, ws := range ws_list {
		if b == ws {
			return true
		}
	}
	return false
}

// 渡されたマジックナンバーによってPBMかPGMかPPMかを判別
type pnmType int

const (
	PBM pnmType = iota + 1
	PGM
	PPM
	ERR
)

func sortPNM(magic string) pnmType {
	var (
		i        int
		pbmMagic = [2]string{"P1", "P4"}
		pgmMagic = [2]string{"P2", "P5"}
		ppmMagic = [2]string{"P3", "P6"}
	)
	for i = 0; i < 2; i++ {
		switch magic {
		case pbmMagic[i]:
			return PBM
		case pgmMagic[i]:
			return PGM
		case ppmMagic[i]:
			return PPM
		}
	}
	return ERR
}
