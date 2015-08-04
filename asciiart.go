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
	BLACK      = '@'
	CHARCOAL   = '#'
	DARKGRAY   = '8'
	MEDIUMGRAY = '&'
	MEDIUM     = 'o'
	GRAY       = ':'
	SLATEGRAY  = '*'
	LIGHTGRAY  = '.'
	WHITE      = ' '
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
				art[pos] = WHITE
			} else if r >= 200 {
				art[pos] = LIGHTGRAY
			} else if r >= 180 {
				art[pos] = SLATEGRAY
			} else if r >= 160 {
				art[pos] = GRAY
			} else if r >= 130 {
				art[pos] = MEDIUM
			} else if r >= 100 {
				art[pos] = MEDIUMGRAY
			} else if r >= 70 {
				art[pos] = DARKGRAY
			} else if r >= 50 {
				art[pos] = CHARCOAL
			} else {
				art[pos] = BLACK
			}
			pos++
		}
		art[pos] = '\n'
		pos++
	}

	return string(art)
}

func New(title string, img image.Image) *AsciiArt {
	return &AsciiArt{title, generateAsciiArt(img)}
}
