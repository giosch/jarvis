package main

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {

	// Load client cert
	cert, err := tls.LoadX509KeyPair("public.pem", "private.pem")
	if err != nil {
		log.Fatal(err)
	}

	// Load CA cert
	caCert, err := ioutil.ReadFile("../core/public.pem")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
		//FIXME this needs to go away, just here for testing ease
		InsecureSkipVerify: true,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	resp, err := client.Post("https://localhost:8443/command", "application/json", strings.NewReader(`{"key":"value"}`))
	if err != nil {
		log.Fatal(err)
	}
	io.Copy(os.Stdout, resp.Body)
}
