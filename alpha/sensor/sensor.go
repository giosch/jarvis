package main

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"bytes"
	"bufio"
	"encoding/json"
	"fmt"
)

type SensorMessage struct {
	Destination string
	IsAudio     bool
	Message     []byte
}

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
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		toSend := SensorMessage{"dummyname",false,[]byte(input)}
		jsonGenerated,_ := json.Marshal(toSend)
		resp, err := client.Post("https://localhost:8443/command", "application/json", bytes.NewReader(jsonGenerated))
		if err != nil {
			log.Fatal(err)
		}
		io.Copy(os.Stdout, resp.Body)
		fmt.Println("")
	}
}
