package main

import (
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
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
			log.Println("FileUploaded")
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
		_, err = io.Copy(w, file)
		if err != nil {
			log.Println(err)
			return
		}

		headers := w.Header()
		headers.Set("Content-Disposition", "attachment; filename="+filename)
		//headers.Set("Content-Type", t)
		//w.Write(file)
		log.Println("FileDownloaded")
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
