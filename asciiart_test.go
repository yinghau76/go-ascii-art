package asciiart

import (
	"image"
	"testing"
)

func Test(t *testing.T) {
	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{160, 160}})
	art := generateAsciiArt(img)
	if len(art) != (80+1)*80 {
		t.Error("Does not generate art with width 80")
	}
}
