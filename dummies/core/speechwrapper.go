package main

import (
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

func getResult(resp *speechpb.RecognizeResponse, err error) string {
	if err != nil {
		return err.Error()
	}
	res := resp.Results[0]
	var curConf float32
	var transcript string
	for _, alt := range res.GetAlternatives() {
		if alt.GetConfidence() > curConf {
			curConf = alt.GetConfidence()
			transcript = alt.GetTranscript()
		}
	}
	return transcript
}
