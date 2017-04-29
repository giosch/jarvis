#!/usr/bin/python

import sys, os
from pocketsphinx.pocketsphinx import *
from sphinxbase.sphinxbase import *
import pyaudio
import wave
import audioop
from collections import deque
import subprocess
import math
import time
from executer import Executer

DEBUG = 1
def plog(s):
    if DEBUG:
        print s

myExecuter = Executer()

CHUNK = 1024  # CHUNKS of bytes to read each time from mic
FORMAT = pyaudio.paInt16
CHANNELS = 1
RATE = 16000
#finally is in DB!!!
THRESHOLD = 55
#number of frame over the THRESHOLD to define the command started
FRAME_TO_START = 10
FRAME_TO_STOP = 3
ACTIVATION_TOKEN = 'jarvis bitch'
# TOKEN_THRESHOLD = '1e-35' # for silence
TOKEN_THRESHOLD = 1e-37

modeldir = "model"
datadir = "test/data"
# Create a decoder with certain model
config = Decoder.default_config()
config.set_string('-hmm', os.path.join(modeldir, 'en-us/en-us'))
config.set_string('-dict', os.path.join(modeldir, 'en-us/cmudict-en-us.dict'))
config.set_string('-keyphrase', ACTIVATION_TOKEN)
config.set_float('-kws_threshold', TOKEN_THRESHOLD)

#opening Microphone
p = pyaudio.PyAudio()
stream = p.open(format=FORMAT, channels=CHANNELS, rate=RATE, input=True, frames_per_buffer=CHUNK)
stream.start_stream()

# Process audio chunk by chunk. On keyphrase detected perform action and restart search
decoder = Decoder(config)
decoder.start_utt()


def save_speech(data):
    filename = 'output_'+str(int(time.time()))
    # writes data to WAV file
    data = ''.join(data)
    wf = wave.open(filename + '.wav', 'wb')
    wf.setnchannels(CHANNELS)
    wf.setsampwidth(p.get_sample_size(FORMAT))
    wf.setframerate(RATE)
    wf.writeframes(data)
    wf.close()
    return filename + '.wav'


#iniziamo facendolo partire subito. Poi vediamo se aggiungere anche la gestione
#del silenzio all'inizio
def recordUntilSilence(silenceSeconds = 1,maxSeconds = 5,prevSeconds = 1.5):
    filename = ""
    audio2send = []
    cur_data = ''  # current chunk of audio data
    rel = RATE/CHUNK
    slid_win = deque(maxlen=silenceSeconds * rel)
    #Prepend audio from 0.5 seconds before noise was detected
    prev_audio = deque(maxlen=prevSeconds * rel)

    assert((silenceSeconds * rel)>= FRAME_TO_START)

    started = False
    chunkPassed = 0
    while True:
        chunkPassed += 1
        if (chunkPassed > (rel * maxSeconds)):
            plog("Timeout recording, decoding phrase")
            filename = save_speech(list(prev_audio) + audio2send)
            #filename = save_speech(audio2send)
            break
        cur_data = stream.read(CHUNK)
        rms = audioop.rms(cur_data,2)
        decibel = 20 * math.log10(rms)
        slid_win.append(decibel)

        if DEBUG:
            if len(slid_win)>=1 :
                plog("Current chunk volume : "
                     + str(slid_win[-1]) +
                     " silence" if slid_win[-1]<THRESHOLD else " someone speaking")

        if sum([x > THRESHOLD for x in slid_win]) >= FRAME_TO_START:
            if started == False:
                plog("Starting recording of phrase")
                started = True
            audio2send.append(cur_data)
        elif started and (sum([x > THRESHOLD for x in slid_win]) <= FRAME_TO_STOP):
            plog("Finished recording, decoding phrase")
            filename = save_speech(list(prev_audio) + audio2send)
            #filename = save_speech(audio2send)
            break
        else:
            prev_audio.append(cur_data)
    return filename

#need to add try catch everywhere
#we want to fail safe
def analyzeAndExecute(filename):
    r = subprocess.check_output(['python','transcribe.py', filename])
    os.remove(filename)
    plog("Command words :"+r)
    myExecuter.do(r)
    return r



while True:
    buf = stream.read(1024)
    if buf:
        decoder.process_raw(buf, False, False)
    else:
        break
    if decoder.hyp() != None:
        #print ([(seg.word, seg.prob, seg.start_frame, seg.end_frame) for seg in decoder.seg()])
        plog("Detected keyphrase")
        plog("Starting recording to send to gspeech")
        filename = recordUntilSilence(silenceSeconds = 1.5,maxSeconds = 5)
        if filename:
            analyzeAndExecute(filename) #on separate thread
        plog("Waiting for the keyphrase")
        decoder.end_utt()
        decoder.start_utt()
