package main

import (
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

func (aa *AsciiArt) Generate(img image.Image) {
	rect := img.Bounds()

	gray := image.NewGray(rect)
	draw.Draw(gray, rect, img, rect.Min, draw.Src)

	size := rect.Size()
	art := make([]byte, (size.X+1)*size.Y)
	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			r, _, _, _ := gray.At(x, y).RGBA()
			if r >= 128 {
				art[y*(size.X+1)+x] = ' '
			} else if r >= 64 {
				art[y*(size.X+1)+x] = '.'
			} else {
				art[y*(size.X+1)+x] = '*'
			}
			art[(y+1)*(size.X+1)-1] = '\n'
		}
	}

	// Convert image to ascii art
	aa.Art = string(art)
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

	image, _, err := image.Decode(file)
	if err != nil {
		log.Println("Failed to decoder uploaded image:", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	aa := &AsciiArt{Title: header.Filename}
	aa.Generate(image)

	err = templates.ExecuteTemplate(rw, "ascii-art.html", aa)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/new/", NewHandler)
	http.HandleFunc("/upload/", UploadHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Failed to start HTTP server:", err.Error())
	}
}
