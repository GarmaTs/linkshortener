package fileserver

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

type UploadHandler struct {
	UploadDir string
	HostAddr  string
}

func canSaveFile(pathToFile string) bool {
	_, err := os.Stat(pathToFile)
	if err != nil {
		// File does not exists, then can save it
		if errors.Is(err, os.ErrNotExist) {
			return true
		}
	}

	return false
}

// makeNewFileName checks if file exists or not
// If exists, makes new name that does not exists
func (h *UploadHandler) makeNewFileName(pathToFile string) (string, error) {
	if canSaveFile(pathToFile) {
		return pathToFile, nil
	}

	i := 0
	for i < 100 {
		i++
		nameWithoutExt := strings.TrimSuffix(pathToFile, path.Ext(pathToFile))
		newName := nameWithoutExt + "_copy (" + strconv.Itoa(i) + ")" + path.Ext(pathToFile)

		if canSaveFile(newName) {
			return newName, nil
		}
	}

	return "", errors.New("Can not make new file name")
}

func (h *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}

	// Create upload directory if not exists
	_, err = os.Stat(h.UploadDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = os.Mkdir(h.UploadDir, 0777)
			if err != nil {
				log.Println(err)
				http.Error(w, "Unable to save file", http.StatusInternalServerError)
				return
			}
		} else {
			if err != nil {
				log.Println(err)
				http.Error(w, "Unable to save file", http.StatusInternalServerError)
				return
			}
		}
	}

	filePath, err := h.makeNewFileName(h.UploadDir + "/" + header.Filename)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}

	err = os.WriteFile(filePath, data, 0777)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}

	fileLink := h.HostAddr + "/" + header.Filename
	fmt.Fprintln(w, fileLink)
}
