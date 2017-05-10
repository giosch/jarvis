package main

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

//I'm using the same certs for both actuators and https server
//Server listens on 8443 for HTTPS clients and on 1337 for actuators
//When the first actuator connects the server performs a switch to
//interactive connection

const serverMessage = `Welcome to jarvis dummy interface`

func main() {
	log.SetFlags(log.Lshortfile)
	certPools := x509.NewCertPool()
	appendCertFromFile("../actuator/public.pem", certPools)
	appendCertFromFile("../sensor/public.pem", certPools)

	serverCert, _ := tls.LoadX509KeyPair("public.pem", "private.pem")

	tlsConfig := &tls.Config{
		ClientCAs:    certPools,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{serverCert},
	}
	tlsConfig.BuildNameToCertificate()

	go listenForSensors(tlsConfig)
	listenForActuators(tlsConfig)
}

func appendCertFromFile(path string, caCertPool *x509.CertPool) {
	caCert, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool.AppendCertsFromPEM(caCert)

}

func listenForSensors(tlsConfig *tls.Config) {
	HelloSensor := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, serverMessage)
		data, _ := ioutil.ReadAll(req.Body)
		println(string(data))
	}

	http.HandleFunc("/", HelloSensor)

	server := &http.Server{
		Addr:      ":8443",
		TLSConfig: tlsConfig,
	}

	server.ListenAndServeTLS("public.pem", "private.pem") //private cert
}

func listenForActuators(tlsConfig *tls.Config) {
	handleConnection := func(conn net.Conn) {
		defer conn.Close()
		n, err := conn.Write([]byte(serverMessage))
		if err != nil {
			log.Println(n, err)
			return
		}
		go io.Copy(conn, os.Stdin)
		io.Copy(os.Stdout, conn)
	}
	ln, err := tls.Listen("tcp", ":1337", tlsConfig)
	if err != nil {
		log.Println(err)
		return
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}
