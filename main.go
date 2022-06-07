package main

import (
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
	"os"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")
	if filename != "" {
		file, err := Create("storage/" + filename)
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
			log.Println("File uploaded")
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
		file, err := Open("storage/" + filename)
		if err != nil {
			log.Println(err)
			w.Write([]byte("Something went wrong"))
			return
		}

		headers := w.Header()
		headers.Set("Content-Disposition", "attachment; filename="+filename)

		_, err = io.Copy(w, file)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("File downloaded")
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("You should specify filename"))
	}
}

func main() {
	r := chi.NewRouter()
	r.Post("/upload", UploadHandler)
	r.Get("/download", DownloadHandler)

	if err := os.MkdirAll("storage", os.ModePerm); err != nil {
		log.Println(err)
	}

	http.ListenAndServe(":3333", r)
}
