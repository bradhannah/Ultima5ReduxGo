package text

import (
	"bytes"
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/text/language"
	"log"
)

var (
	//go:embed assets/ultima-v-warriors-of-destiny/ultima-v-warriors-of-destiny.ttf
	rawUltimaTTF []byte
)

type UltimaFont struct {
	rawUltimaTTF     *[]byte
	ultimaFaceSource *text.GoTextFaceSource
	textFace         *text.GoTextFace
}

func NewUltimaFont() *UltimaFont {
	ultimaFont := UltimaFont{}
	ultimaFont.rawUltimaTTF = &rawUltimaTTF

	var err error
	ultimaFont.ultimaFaceSource, err = text.NewGoTextFaceSource(bytes.NewReader(rawUltimaTTF))
	if err != nil {
		log.Fatal(err)
	}

	ultimaFont.textFace = &text.GoTextFace{
		Source:    ultimaFont.ultimaFaceSource,
		Direction: text.DirectionLeftToRight,
		Size:      24,
		Language:  language.English,
	}

	return &ultimaFont
}
