package main

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/anthonynsimon/bild/imgio"
)

type imageHandler struct {
	mu sync.Mutex // guards n
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func (handler *imageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	imgUrl := "https://cdn.raster.app/raster/raster/eMueOFL3Hw?ixlib=js-3.6.0&s=cc3998a09bd1b0d388497ae6148c38db"

	handler.mu.Lock()

	// Unlocks at the end of the func
	defer handler.mu.Unlock()

	img := fetch(imgUrl)

	transformed := transform(img, r)

	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, transformed, nil)

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
	print("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func fetch(imgUrl string) image.Image {
	hashed := fmt.Sprint(hash(imgUrl))
	path := "./cache/" + hashed + ".jpeg"

	img, err := imgio.Open(path)

	if err != nil {
		response, e := http.Get(imgUrl)

		print("Fetching...")

		if e != nil {
			log.Fatal(e)
		}

		defer response.Body.Close()

		os.Mkdir("cache", 0777)
		file, err := os.Create(path)

		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()

		_, err = io.Copy(file, response.Body)
		if err != nil {
			log.Fatal(err)
		}

		img, err = imgio.Open(path)

		if err != nil {
			log.Fatal(err)
			return nil
		}
	}

	return img
}
