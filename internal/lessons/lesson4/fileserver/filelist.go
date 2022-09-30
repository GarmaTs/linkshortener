package fileserver

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
)

type FileListHander struct {
	Dir       string
	extension string
}

func (fh *FileListHander) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fh.extension = ".*"

	switch r.Method {
	case http.MethodGet:
		ext := r.FormValue("extension")
		if len(ext) > 0 {
			fh.extension = ext
		}
	}

	_, err := os.Stat(fh.Dir)
	if err != nil {
		log.Println(err)
		fmt.Fprintln(w, "No files")
		return
	}

	files, err := os.ReadDir(fh.Dir)
	if err != nil {
		log.Println(err)
		fmt.Fprintln(w, "No files")
		return
	}

	for _, f := range files {
		fi, err := os.Stat(fh.Dir + "/" + f.Name())
		if err != nil {
			log.Fatal(err)
		}
		if fh.extension == ".*" {
			fmt.Fprintf(w, fmt.Sprintf("%s\t%s\t%d\n", f.Name(), path.Ext(f.Name()), fi.Size()))
		} else {
			if fh.extension == path.Ext(f.Name()) {
				fmt.Fprintf(w, fmt.Sprintf("%s\t%s\t%d\n", f.Name(), path.Ext(f.Name()), fi.Size()))
			}
		}
	}
}
