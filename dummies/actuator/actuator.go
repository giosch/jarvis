package main

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"log"
	"os"
)

//Connects to the server and drops to interactive mode
func main() {
	log.SetFlags(log.Lshortfile)

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

	conn, err := tls.Dial("tcp", "127.0.0.1:1337", tlsConfig)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	go io.Copy(os.Stdout, conn)
	io.Copy(conn, os.Stdin)
}
