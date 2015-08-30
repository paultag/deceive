package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Config struct {
	Host   string
	Port   int
	Cert   string
	Key    string
	CaCert string
	Root   string
}

func GetConfig() Config {
	config := Config{
		Host:   "localhost",
		Port:   1984,
		Cert:   "/etc/deceive/deceive.crt",
		Key:    "/etc/deceive/deceive.key",
		CaCert: "/etc/deceive/ca.crt",
		Root:   "/var/lib/deceive/",
	}

	host := flag.String("host", config.Host, "host to serve on")
	port := flag.Int("port", config.Port, "port to serve on")
	cert := flag.String("cert", config.Cert, "server cert to use")
	key := flag.String("key", config.Key, "server key to use")
	ca := flag.String("ca", config.CaCert, "ca cert to use")
	root := flag.String("root", config.Root, "filesystem root")

	flag.Parse()

	config.Host = *host
	config.Port = *port
	config.Cert = *cert
	config.Key = *key
	config.CaCert = *ca
	config.Root = *root

	return config
}

func main() {
	config := GetConfig()

	caPool := x509.NewCertPool()
	x509CaCrt, err := ioutil.ReadFile(config.CaCert)
	if err != nil {
		panic(err)
	}
	if ok := caPool.AppendCertsFromPEM(x509CaCrt); !ok {
		panic(fmt.Errorf("Error appending CA cert from PEM!"))
	}

	s := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler: http.DefaultServeMux,
		TLSConfig: &tls.Config{
			ClientAuth: tls.RequireAndVerifyClientCert,
			ClientCAs:  caPool,
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		clientName := r.TLS.PeerCertificates[0].Subject.CommonName
		HandleUpload(config, w, r, clientName)
	})

	log.Printf("Listening...\n")
	log.Fatal(s.ListenAndServeTLS(config.Cert, config.Key))
}
