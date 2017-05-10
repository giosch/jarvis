package main

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	//This skips certificate validation
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Post("https://localhost:8443/", "application/json", strings.NewReader(`{"key":"value"}`))
	if err != nil {
		log.Fatal(err)
	}
	io.Copy(os.Stdout, resp.Body)
}
