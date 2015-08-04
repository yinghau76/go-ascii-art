package main

import (
	"image"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/yinghau76/go-ascii-art"
)

var templates = template.Must(template.ParseFiles("new.html", "ascii-art.html"))

func newHandler(rw http.ResponseWriter, req *http.Request) {
	renderTemplate(rw, "new", nil)
}

func uploadHandler(rw http.ResponseWriter, req *http.Request) {
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

	aa := asciiart.New(header.Filename, img)
	renderTemplate(rw, "ascii-art", aa)
}

func renderTemplate(rw http.ResponseWriter, name string, data interface{}) error {
	err := templates.ExecuteTemplate(rw, name+".html", data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err
	}
	return nil
}

func init() {
	http.HandleFunc("/art/", uploadHandler)
	http.HandleFunc("/", newHandler)
}

func main() {
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatal("Failed to start HTTP server:", err.Error())
	}
}
