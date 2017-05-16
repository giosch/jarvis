package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

const (
	timeoutDuration      = time.Second * 10
	secureActuatorPort   = "1337"
	insecureActuatorPort = "1338"
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

func (am *ActuatorMessage) Error() string {
	return am.ErrorMessage
}

//The currently supposedly connected actuators
var actuators ActuatorPool

//A pool of actuators, use it to store actuators
type ActuatorPool struct {
	sync.RWMutex
	curid     int64
	actuators map[string]map[int64]Actuator
}

//Use this to thread safely add an actuator to the pool upon connection
func (ap *ActuatorPool) Add(a Actuator) {
	ap.Lock()
	defer ap.Unlock()

	a.id = ap.curid
	if ap.actuators == nil {
		ap.actuators = make(map[string]map[int64]Actuator)
	}
	if ap.actuators[a.name] == nil {
		ap.actuators[a.name] = make(map[int64]Actuator)
	}
	ap.actuators[a.name][ap.curid] = a
	ap.curid++
}

//Thread safely remove an actuator, upgrades an RLock to a WLock
func (ap *ActuatorPool) Remove(a *Actuator) {
	ap.Lock()
	defer ap.Unlock()
	delete(ap.actuators[a.name], a.id)
}

//Use this to send a message to an actuator in the pool
func (ap *ActuatorPool) Send(target string, msg *ActuatorMessage) (*ActuatorMessage, error) {
	//The reason we get the full lock and not just the read lock is that we can
	//modify the list of actuators if some disconnects in the meantime
	ap.Lock()
	defer ap.Unlock()

	acts, found := ap.actuators[target]
	if !found || len(acts) == 0 /*delete also the parent in this case?*/ {
		return nil, fmt.Errorf("actuator not found: %s", target)
	}

	results := make(chan *ActuatorMessage, len(acts))
	var wg sync.WaitGroup
	wg.Add(len(acts))
	for _, act := range acts {
		//Parameter is passed to avoid races
		go func(act Actuator,msg *ActuatorMessage) {
			res, e := act.Send(msg)
			if e != nil {
				res = &ActuatorMessage{ErrorMessage: e.Error()}
				delete(ap.actuators[act.name], act.id)
			}
			results <- res
			wg.Done()
		}(act,msg)
	}
	//Wait for all comms to complete
	wg.Wait()
	//This is used to make the following loop conclude gracefully
	close(results)
	//If there was at least one success we return that message,
	//otherwise we return the last error (which is less likely to be a disconnection)
	var res *ActuatorMessage
	for res = range results {
		if res.Error() == "" {
			return res, nil
		}
	}
	//Res implements the builtin error interface, so it can be returned as error value
	return nil, res
}

//This is used to represent a connected actuator and communicate with it
type Actuator struct {
	name string
	id   int64
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
	//If it can handle timeouts, let's use them
	setDeadline(ac.conn)
	var err error
	err = ac.enc.Encode(*msg)
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

	//If it can handle timeouts, let's use them
	setDeadline(conn)
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

func listenForActuators(tlsConfig *tls.Config) {
	var ln net.Listener
	var err error
	if tlsConfig != nil {
		ln, err = tls.Listen("tcp", ":"+secureActuatorPort, tlsConfig)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		ln, err = net.Listen("tcp", ":"+insecureActuatorPort)
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

func setDeadline(conn io.ReadWriteCloser) {
	nconn, ok := conn.(net.Conn)
	if ok && timeoutDuration != 0 {
		_ = nconn.SetDeadline(time.Now().Add(timeoutDuration))
	}
}
