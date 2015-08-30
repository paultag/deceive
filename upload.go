package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

func HandleUpload(config Config, w http.ResponseWriter, r *http.Request, clientName string) {
	dir, fpath := path.Split(r.URL.Path)
	dir = path.Clean(path.Join("/", dir))
	targetDir := path.Join(config.Root, dir)

	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		log.Printf("Unknown directory\n")
		return
	}

	targetFile := path.Join(targetDir, fpath)

	fd, err := os.Create(targetFile)
	if err != nil {
		log.Printf("Error: %s\n", err)
		return
	}
	defer fd.Close()

	written, err := io.Copy(fd, r.Body)
	if err != nil {
		log.Printf("Error: %s\n", err)
		return
	}
	log.Printf("%s written %d bytes, %s\n", clientName, written, targetFile)
}
