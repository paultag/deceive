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
}

func GetConfig() Config {
	config := Config{
		Host:   "localhost",
		Port:   1984,
		Cert:   "",
		Key:    "",
		CaCert: "",
	}

	host := flag.String("host", config.Host, "host to serve on")
	port := flag.Int("port", config.Port, "port to serve on")
	cert := flag.String("cert", config.Cert, "server cert to use")
	key := flag.String("key", config.Key, "server key to use")
	ca := flag.String("ca", config.CaCert, "ca cert to use")

	flag.Parse()

	config.Host = *host
	config.Port = *port
	config.Cert = *cert
	config.Key = *key
	config.CaCert = *ca

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

	http.HandleFunc("/", HandleUpload)
	log.Printf("Listening...\n")
	log.Fatal(s.ListenAndServeTLS(config.Cert, config.Key))
}
