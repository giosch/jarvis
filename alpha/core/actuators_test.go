package main

import (
	"bytes"
	"encoding/json"
	"log"
	"testing"
)

type MockConn struct {
	//Reading from the connection Reads from what is stored in here
	Input *bytes.Buffer
	//Writing to the MockConn writes to this variable
	Output *bytes.Buffer

	Enc *json.Encoder
	Dec *json.Decoder

	Closed bool
}

func NewMockConn() *MockConn {
	mc := &MockConn{
		Input:  bytes.NewBuffer(nil),
		Output: bytes.NewBuffer(nil),
	}
	mc.Enc = json.NewEncoder(mc.Input)
	mc.Dec = json.NewDecoder(mc.Output)
	return mc
}

func (m *MockConn) Read(p []byte) (n int, err error) {
	log.Println("Reading from mockconn")
	_ = "breakpoint"
	return m.Input.Read(p)
}

func (m *MockConn) Write(p []byte) (n int, err error) {
	log.Println("Writing to mockconn")
	_ = "breakpoint"
	return m.Output.Write(p)
}

func (m *MockConn) Close() error {
	m.Closed = true
	return nil
}

func (m *MockConn) WhatWasWritten() (am *ActuatorMessage) {
	//decodes from conn Input buffer
	_ = m.Dec.Decode(am)
	return
}

func (m *MockConn) WhatWillBeRead(am *ActuatorMessage) {
	//encodes to conn Input buffer
	_ = m.Enc.Encode(am)
}

func TestSend(t *testing.T) {
	//Backup actuators
	savedActs := actuators
	//Restore actuators
	defer func() {
		actuators = savedActs
	}()
	mcc := NewMockConn()
	mcreceive := &ActuatorMessage{
		Payload: "This is a Test",
	}
	log.Println("Creating mock incoming message")
	mcc.WhatWillBeRead(mcreceive)
	log.Println("Creating mock actuator")
	mca := NewActuator("mock", mcc)
	log.Println("Adding actuator to pool")
	actuators.Add(*mca)
	mcsend := &ActuatorMessage{
		Payload: "Payload",
	}
	log.Println("Sending test message")
	res, err := actuators.Send("mock", mcsend)
	log.Println("Sent")
	if res == nil || err != nil {
		if res != nil && res.Error() != "" {
			t.Error("Unexpected error occurred " + res.Error())
		} else {
			t.Error("Unexpected error occurred " + err.Error())
		}
		return
	}
	if res.Payload != mcreceive.Payload || err != nil {
		t.Errorf("Was expecting %s but got %s", mcreceive.Payload, res.Payload)
		return
	}
}
