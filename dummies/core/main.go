package main

import (
	"crypto/tls"
	"io"
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
	go listenForSensors()
	listenForActuators()
}

func listenForSensors() {
	HelloSensor := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, serverMessage)
	}
	http.HandleFunc("/", HelloSensor)
	err := http.ListenAndServeTLS(":8443", "server.pem", "server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func listenForActuators() {
	cer, err := tls.LoadX509KeyPair("server.pem", "server.key")
	if err != nil {
		log.Println(err)
		return
	}

	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	ln, err := tls.Listen("tcp", ":1337", config)
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

func handleConnection(conn net.Conn) {
	defer conn.Close()
	n, err := conn.Write([]byte(serverMessage))
	if err != nil {
		log.Println(n, err)
		return
	}
	go io.Copy(conn, os.Stdin)
	io.Copy(os.Stdout, conn)
}
