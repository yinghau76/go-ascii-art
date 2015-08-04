// Package asciiart can generate ASCII art from provided image.
package asciiart

import (
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/nfnt/resize"
)

type AsciiArt struct {
	Title string
	Art   string
}

const (
	sBLACK      = '@'
	sCHARCOAL   = '#'
	sDARKGRAY   = '8'
	sMEDIUMGRAY = '&'
	sMEDIUM     = 'o'
	sGRAY       = ':'
	sSLATEGRAY  = '*'
	sLIGHTGRAY  = '.'
	sWHITE      = ' '
)

// Reference: http://www.codeproject.com/Articles/20435/Using-C-To-Generate-ASCII-Art-From-An-Image
func generateAsciiArt(img image.Image) string {
	resized := resize.Resize(80, 0, img, resize.Bilinear)
	rect := resized.Bounds()

	gray := image.NewGray(rect)
	draw.Draw(gray, rect, resized, rect.Min, draw.Src)

	size := rect.Size()
	art := make([]byte, (size.X+1)*size.Y)
	pos := 0
	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			r, _, _, _ := gray.At(x, y).RGBA()
			r /= 256
			if r >= 230 {
				art[pos] = sWHITE
			} else if r >= 200 {
				art[pos] = sLIGHTGRAY
			} else if r >= 180 {
				art[pos] = sSLATEGRAY
			} else if r >= 160 {
				art[pos] = sGRAY
			} else if r >= 130 {
				art[pos] = sMEDIUM
			} else if r >= 100 {
				art[pos] = sMEDIUMGRAY
			} else if r >= 70 {
				art[pos] = sDARKGRAY
			} else if r >= 50 {
				art[pos] = sCHARCOAL
			} else {
				art[pos] = sBLACK
			}
			pos++
		}
		art[pos] = '\n'
		pos++
	}

	return string(art)
}

// New creates an ASCII art from an image
func New(title string, img image.Image) *AsciiArt {
	return &AsciiArt{title, generateAsciiArt(img)}
}
