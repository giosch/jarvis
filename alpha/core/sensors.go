package main

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"

	vikyscript "github.com/empijei/VikyScript/dummy"
)

func init() {
	//FIXME this are just dummies
	_, _, _ = vikyscript.Parse("ping:ping (?P<actuator>[a-z]+)")
	_, _, _ = vikyscript.Parse("conncheck:controlla se è connesso l'attuatore (?P<actuator>[a-z]+)")
}

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
	cmd, params, err := vikyscript.Match(text)
	//TODO hm...... e ora? 🤔
	if err != nil {
		log.Println(err)
		return []byte(err.Error())
	}
	log.Printf("Command recognized: %s with parameters %#v\n", cmd, params)

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
	text := googleSpeechToText(buf)
	log.Println(text)
	return text
}
