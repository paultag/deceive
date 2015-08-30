package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

func main() {
	host := flag.String("host", "localhost", "host to serve on")
	port := flag.Int("port", 1984, "port to serve on")
	cert := flag.String("cert", "", "server cert to use")
	key := flag.String("key", "", "server key to use")
	ca := flag.String("ca", "", "ca cert to use")
	root := flag.String("root", "", "root to write to")

	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		clientName := r.TLS.PeerCertificates[0].Subject.CommonName

		dir, fpath := path.Split(r.URL.Path)
		dir = path.Clean(path.Join("/", dir))
		targetDir := path.Join(*root, dir)

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
	})

	caPool := x509.NewCertPool()
	x509CaCrt, err := ioutil.ReadFile(*ca)
	if err != nil {
		panic(err)
	}
	if ok := caPool.AppendCertsFromPEM(x509CaCrt); !ok {
		panic(fmt.Errorf("Error appending CA cert from PEM!"))
	}

	log.Printf("Listening...\n")

	s := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", *host, *port),
		Handler: http.DefaultServeMux,
		TLSConfig: &tls.Config{
			ClientAuth: tls.RequireAndVerifyClientCert,
			ClientCAs:  caPool,
		},
	}
	log.Fatal(s.ListenAndServeTLS(*cert, *key))
}
