package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

//Just a wrapper for messages
type ActuatorMessage struct {
	//TODO MetaData seems poorly defined as a string. Refine this struct
	MetaData string
	//Using a string to simplify stuff, []byte is serialized as base64 in json
	//since we are targeting simple systems i would like to avoid having an
	//arduino b64(ende)code a simple command
	Payload string
}

//The currently supposedly connected actuators
var actuators ActuatorPool

//A pool of actuators, use it to store actuators
type ActuatorPool struct {
	sync.RWMutex
	actuators map[string]Actuator
}

//Use this to thread safely add an actuator to the pool upon connection
func (ap *ActuatorPool) Add(a Actuator) {
	ap.Lock()
	defer ap.Unlock()
	ap.actuators[a.name] = a
}

//Use this to send a message to an actuator in the pool
func (ap *ActuatorPool) Send(target string, msg *ActuatorMessage) (*ActuatorMessage, error) {
	ap.RLock()
	defer ap.RUnlock()

	act, found := ap.actuators[target]
	if !found {
		return nil, fmt.Errorf("actuator not found: %s", target)
	}

	//propagates both error and return value
	return act.Send(msg)
}

//This is used to represent a connected actuator and communicate with it
type Actuator struct {
	name string
	enc  *json.Encoder
	dec  *json.Decoder
	conn io.ReadWriteCloser
}

func NewActuator(name string, conn io.ReadWriteCloser) *Actuator {
	return &Actuator{
		name: name,
		enc:  json.NewEncoder(conn),
		dec:  json.NewDecoder(conn),
		conn: conn,
	}
}

func (ac *Actuator) Send(msg *ActuatorMessage) (*ActuatorMessage, error) {
	//TODO handle disconnection
	err := ac.enc.Encode(*msg)
	if err != nil {
		return nil, err
	}
	err = ac.dec.Decode(msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func handleActuator(conn net.Conn) {
	//This is meant as a reminder that we are ignoring an error
	defer func() { _ = conn.Close() }()
	n, err := conn.Write([]byte(serverMessage))
	if err != nil {
		log.Println(n, err)
		return
	}
	name, err := authenticateActuator(&conn)
	if err != nil {
		return
	}
	actuators.Add(Actuator{name: name, conn: conn})
}

func authenticateActuator(conn *net.Conn) (string, error) {
	//TODO read subject from tls connection or authenticate the actuator somehow
	return "dummyname", nil
}

//TODO make ports configurable
func listenForActuators(tlsConfig *tls.Config) {
	var ln net.Listener
	var err error
	if tlsConfig != nil {
		ln, err = tls.Listen("tcp", ":1337", tlsConfig)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		ln, err = net.Listen("tcp", ":1338")
		if err != nil {
			log.Println(err)
			return
		}
	}
	defer func() { _ = ln.Close() }()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleActuator(conn)
	}
}
