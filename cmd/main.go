package main

import (
	"log"
	"net/http"
	"time"

	"github.com/GarmaTs/linkshortener/internal/lesson4/fileserver"
)

func main() {
	fileDir := "upload"

	uploadHandler := &fileserver.UploadHandler{
		UploadDir: fileDir,
	}
	upSrv := &http.Server{
		Addr:         ":80",
		Handler:      uploadHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Fatal(upSrv.ListenAndServe())
	}()

	fileListHander := &fileserver.FileListHander{
		Dir: fileDir,
	}
	flSrv := &http.Server{
		Addr:         ":8080",
		Handler:      fileListHander,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Fatal(flSrv.ListenAndServe())
}
