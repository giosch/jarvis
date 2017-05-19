package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"os"
)

//I'm using the same certs for both actuators and https server
//Server listens on 8443 for HTTPS clients and on 1337 for actuators
//When the first actuator connects the server performs a switch to
//interactive connection

const serverMessage = `Welcome to jarvis dummy interface`

func main() {
	log.SetFlags(log.Lshortfile)
	tlsConfig := loadCerts("./trustedCerts")
	go listenForSensors(tlsConfig)
	go listenForActuators(tlsConfig)
	listenForActuators(nil)
}

func loadCerts(path string) *tls.Config {
	certsPool := x509.NewCertPool()
	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		filePath := path + string(os.PathSeparator) + f.Name()
		log.Println("Loading certs for: " + filePath)
		cert, _ := ioutil.ReadFile(filePath)
		certsPool.AppendCertsFromPEM(cert)
	}
	log.Printf("Loaded %d certs.\n", len(certsPool.Subjects()))

	serverCert, _ := tls.LoadX509KeyPair("public.pem", "private.pem")

	tlsConfig := &tls.Config{
		ClientCAs:    certsPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{serverCert},
	}
	tlsConfig.BuildNameToCertificate()
	return tlsConfig
}
