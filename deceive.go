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
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"pault.ag/go/config"
)

type Deceive struct {
	Host   string `flag:"host" description:"host to serve on behalf on"`
	Port   int    `flag:"port" description:"server port to host on"`
	Cert   string `flag:"cert" description:"server tls cert"`
	Key    string `flag:"key" description:"server tls key"`
	CaCert string `flag:"ca" description:"ca cert"`
	Root   string `flag:"root" description:"filesystem root"`
}

func GetConfig() Deceive {
	conf := Deceive{
		Host:   "localhost",
		Port:   1984,
		Cert:   "/etc/deceive/deceive.crt",
		Key:    "/etc/deceive/deceive.key",
		CaCert: "/etc/deceive/ca.crt",
		Root:   "/var/lib/deceive/",
	}
	flags, err := config.LoadFlags("deceive", &conf)
	if err != nil {
		panic(err)
	}
	flags.Parse(os.Args[1:])

	if !path.IsAbs(conf.Root) {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		conf.Root = path.Clean(path.Join(cwd, conf.Root))
	}

	return conf
}

func main() {
	conf := GetConfig()

	caPool := x509.NewCertPool()
	x509CaCrt, err := ioutil.ReadFile(conf.CaCert)
	if err != nil {
		panic(err)
	}
	if ok := caPool.AppendCertsFromPEM(x509CaCrt); !ok {
		panic(fmt.Errorf("Error appending CA cert from PEM!"))
	}

	s := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Handler: http.DefaultServeMux,
		TLSConfig: &tls.Config{
			ClientAuth: tls.RequireAndVerifyClientCert,
			ClientCAs:  caPool,
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		clientName := r.TLS.PeerCertificates[0].Subject.CommonName
		l := func(message string, args ...interface{}) {
			log.Printf("%s: %s", clientName, fmt.Sprintf(message, args...))
		}
		HandleUpload(l, conf, w, r, clientName)
	})

	log.Printf("Listening...\n")
	log.Fatal(s.ListenAndServeTLS(conf.Cert, conf.Key))
}

// vim: foldmethod=marker
