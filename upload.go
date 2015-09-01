/* {{{ Copyright (c) Paul R. Tagliamonte <paultag@gmail.com>, 2015
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE. }}} */

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

func writeError(w http.ResponseWriter, message string, code int) error {
	return writeJSON(w, map[string]string{
		"message": "failure",
		"error":   message,
	}, code)
}

func writeSuccess(w http.ResponseWriter, data interface{}, code int) error {
	return writeJSON(w, data, 200)
}

func writeJSON(w http.ResponseWriter, data interface{}, code int) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return err
	}
	return nil
}

func HandlePUT(config Deceive, w http.ResponseWriter, r *http.Request, clientName string) {
	l := func(message string, args ...interface{}) {
		log.Printf("%s: %s", clientName, fmt.Sprintf(message, args...))
	}
	defer l("End request")

	l("Incoming request to push to %s", r.URL.Path)
	dir, fpath := path.Split(r.URL.Path)
	dir = path.Clean(path.Join("/", dir))
	targetDir := path.Join(config.Root, dir)

	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		l("Attempting to write to an unknown directory")
		writeError(w, "unknown archive", 400)
		return
	}

	targetFile := path.Clean(path.Join(targetDir, fpath))
	if !strings.HasPrefix(targetFile, config.Root) {
		l(
			"Caught an attempt to write outside the root! Whoah! %s",
			targetFile, config.Root,
		)
		writeError(w, "unknown archive", 400) // Don't let the client know..
		return
	}

	fd, err := os.Create(targetFile)
	if err != nil {
		l("Error creating target: %s: %s", targetFile, err)
		writeError(w, "error creating target!", 500)
		return
	}
	defer fd.Close()

	l("Starting write to target filename")
	written, err := io.Copy(fd, r.Body)
	if err != nil {
		l("Error writing to target: %s: %s", targetFile, err)
		writeError(w, "error writing to target!", 500)
		return
	}
	l("Wrote %d bytes.", written)
	writeSuccess(w, map[string]string{
		"message": fmt.Sprintf("Wrote %d bytes", written),
	}, 200)
}

func HandleGET(config Deceive, w http.ResponseWriter, r *http.Request, clientName string) {
}

func HandleUpload(config Deceive, w http.ResponseWriter, r *http.Request, clientName string) {
	switch r.Method {
	case "PUT":
		HandlePUT(config, w, r, clientName)
	case "GET":
		HandleGET(config, w, r, clientName)
	default:
		log.Printf("Unknown method\n")
	}
}

// vim: foldmethod=marker
