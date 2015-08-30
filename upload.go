package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	clientName := r.TLS.PeerCertificates[0].Subject.CommonName

	dir, fpath := path.Split(r.URL.Path)
	dir = path.Clean(path.Join("/", dir))
	targetDir := path.Join("/tmp/fnord/", dir)

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
