package main

import (
	"bytes"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/anthonynsimon/bild/imgio"
)

type imageHandler struct {
	mu    sync.Mutex // guards n
	count int
}

func (handler *imageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	imgUrl := "https://cdn.raster.app/raster/raster/k5Ov2Vn1pi?ixlib=js-3.6.0&s=bd19ef4b04ed05e37b36aef3b6067de2"
	handler.mu.Lock()

	// Unlocks at the end of the func
	defer handler.mu.Unlock()

	response, e := http.Get(imgUrl)

	if e != nil {
		log.Fatal(e)
	}

	defer response.Body.Close()

	file, err := os.Create("input.jpeg")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}

	img, err := imgio.Open("input.jpeg")

	if err != nil {
		log.Fatal(err)
		return
	}

	transformed := transform(img, r)

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, transformed, nil)

	if err != nil {
		log.Fatal(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(buf.Bytes())
}

func main() {
	http.Handle("/new", new(imageHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
