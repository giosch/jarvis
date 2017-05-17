package main

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
)

type SensorMessage struct {
	Destination string
	IsAudio     bool
	Message     []byte
}

func listenForSensors(tlsConfig *tls.Config) {
	http.HandleFunc("/command", handleSensors)

	server := &http.Server{
		Addr:      ":8443",
		TLSConfig: tlsConfig,
	}

	err := server.ListenAndServeTLS("public.pem", "private.pem")
	panic(err)
}

func handleSensors(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var msg SensorMessage
	err := decoder.Decode(&msg)
	if err != nil {
		log.Println(err)
	}
	defer func() { _ = req.Body.Close() }()
	var text string
	if msg.IsAudio {
		text = speechToText(msg.Message)
	} else {
		text = string(msg.Message)
	}
	log.Println("Testo ricevuto", text)
	//TODO define metadata type and function
	var metadata string
	_, err = w.Write(handleText(msg.Destination, text, metadata))
	if err != nil {
		log.Println(err)
	}
}

//TODO define what exactly is metadata, a string seems pretty dull...
func handleText(target, text, metadata string) []byte {
	//TODO invoke vikyscript
	//TODO take vikyscript command output

	//TODO check if the response requires to repeat the job and keep looping
	//for keeplooping || err == nil{
	response, err := actuators.Send(target, &ActuatorMessage{metadata, text, ""})
	if err != nil {
		return []byte("TODO this is a temporary error message: " + err.Error())
	}
	//}

	return []byte(response.Payload)
}

func speechToText(buf []byte) string {
	//TODO
	return ""
}
