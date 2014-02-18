package main

import (
	"github.com/nfnt/resize"
	"html/template"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
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
	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			r, _, _, _ := gray.At(x, y).RGBA()
			r /= 256
			if r >= 230 {
				art[y*(size.X+1)+x] = WHITE
			} else if r >= 200 {
				art[y*(size.X+1)+x] = LIGHTGRAY
			} else if r >= 180 {
				art[y*(size.X+1)+x] = SLATEGRAY
			} else if r >= 160 {
				art[y*(size.X+1)+x] = GRAY
			} else if r >= 130 {
				art[y*(size.X+1)+x] = MEDIUM
			} else if r >= 100 {
				art[y*(size.X+1)+x] = MEDIUMGRAY
			} else if r >= 70 {
				art[y*(size.X+1)+x] = DARKGRAY
			} else if r >= 50 {
				art[y*(size.X+1)+x] = CHARCOAL
			} else {
				art[y*(size.X+1)+x] = BLACK
			}
		}
		art[(y+1)*(size.X+1)-1] = '\n'
	}

	return string(art)
}

func NewAsciiArt(title string, img image.Image) *AsciiArt {
	return &AsciiArt{title, generateAsciiArt(img)}
}

var templates = template.Must(template.ParseFiles("new.html", "ascii-art.html"))

func NewHandler(rw http.ResponseWriter, req *http.Request) {
	err := templates.ExecuteTemplate(rw, "new.html", nil)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func UploadHandler(rw http.ResponseWriter, req *http.Request) {
	file, header, err := req.FormFile("image")
	if err != nil {
		log.Println("Not image was found:", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	img, _, err := image.Decode(file)
	if err != nil {
		log.Println("Failed to decoder uploaded image:", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	aa := NewAsciiArt(header.Filename, img)

	err = templates.ExecuteTemplate(rw, "ascii-art.html", aa)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/art/", UploadHandler)
	http.HandleFunc("/", NewHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Failed to start HTTP server:", err.Error())
	}
}
