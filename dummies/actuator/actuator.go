package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"log"
)

//Just a wrapper for messages
type ActuatorMessage struct {
	//TODO MetaData seems poorly defined as a string. Refine this struct
	MetaData string
	//Using a string to simplify stuff, []byte is serialized as base64 in json
	//since we are targeting simple systems i would like to avoid having an
	//arduino b64(ende)code a simple command
	Payload string
	//Used to communicate problems
	ErrorMessage string
}

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

	//gestire il messaggio iniziale del server

	in := json.NewDecoder(conn)
	out := json.NewEncoder(conn)
	for {
		var msg *ActuatorMessage
		var err error
		err = in.Decode(msg)
		if err != nil {
			log.Println(err)
			return
		}
		err = out.Encode(&ActuatorMessage{"", "", "Functionality not implemented yet"})
		if err != nil {
			log.Println(err)
			return
		}
	}
}
