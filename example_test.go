package asciiart_test

import (
	"fmt"
	"image"

	"github.com/yinghau76/go-ascii-art"
)

func ExampleAsciiArt() {
	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{160, 160}})
	aa := asciiart.New("image", img)
	if len(aa.Art) != (80+1)*80 {
		fmt.Printf("Failed to generate art")
	} else {
		fmt.Printf("Art generated!")
	}
	// Output: Art generated!
}
