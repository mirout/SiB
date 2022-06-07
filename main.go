package main

import (
	"github.com/go-chi/chi/v5"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")
	if filename != "" {
		file, err := os.Create("storage/" + filename)
		if err != nil {
			log.Println(err)
			w.Write([]byte("Something went wrong"))
			return
		}
		defer file.Close()

		_, err = io.Copy(file, r.Body)

		if err != nil {
			log.Println(err)
			w.Write([]byte("Something went wrong"))
		} else {
			w.Write([]byte("Success"))
		}

	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("You should specify filename"))
	}
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")
	if filename != "" {
		file, err := ioutil.ReadFile("storage/" + filename)
		if err != nil {
			log.Println(err)
			w.Write([]byte("Something went wrong"))
			return
		}
		t := http.DetectContentType(file[:512])
		headers := w.Header()
		headers.Set("Content-Disposition", "attachment; filename="+filename)
		headers.Set("Content-Type", t)
		w.Write(file)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("You should specify filename"))
	}
}

func main() {
	r := chi.NewRouter()
	r.Post("/upload", UploadHandler)
	r.Get("/download", DownloadHandler)

	http.ListenAndServe(":3333", r)
}
