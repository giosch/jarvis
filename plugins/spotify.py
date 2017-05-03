import dbus
import subprocess
import time
import signal
import os
##TODO
#   -adding check if spotify is running, if not open it, every time

#the check does not work i think
class SpotifyLinux:
    def __init__(self):
        self.checkAndLoop()

    def checkAndLoop(self):
        if(self.checkIfRunningAndInit()):
            return
        while(not self.checkIfRunningAndInit()):
            if not self.pidExitst():
                self.startSpotify()
        #iv'e yet to find a better way
        time.sleep(2)
        return

    def initInterface(self):
        self.bus = dbus.SessionBus()
        self.proxy = self.bus.get_object('org.mpris.MediaPlayer2.spotify', '/org/mpris/MediaPlayer2')
        self.interface = dbus.Interface(self.proxy, dbus_interface='org.mpris.MediaPlayer2.Player')
        self.properties = dbus.Interface(self.proxy, dbus_interface='org.freedesktop.DBus.Properties')

    def checkIfRunningAndInit(self):
        try:
            self.initInterface()
            self.getProperties()
            return True
        except dbus.exceptions.DBusException as e:
            return False

    #Fails with a defunct spotify process active, need to fix
    def pidExitst(self):
        try:
            output = subprocess.check_output(["pidof","spotify"])
            if len(output)>0:
                return True
            else:
                return False
        except:
            return False

    def startSpotify(self):
        def preexec_function():
            # Ignore the SIGINT signal by setting the handler to the standard
            # signal handler SIG_IGN.
            signal.signal(signal.SIGINT, signal.SIG_IGN)
        FNULL = open(os.devnull, 'w')
        subprocess.Popen(['spotify'],preexec_fn = os.setpgrp, stdout=FNULL, stderr=FNULL)

    def getProperties(self):
        return self.properties.GetAll('org.mpris.MediaPlayer2.Player').items()

    def getMetadata(self):
        self.checkAndLoop()
        prop = self.getProperties()
        for i in prop:
            if i[0] == 'Metadata':
                return i[1]
        return None

    def next(self):
        self.checkAndLoop()
        self.interface.Next()

    def stop(self):
        self.checkAndLoop()
        self.interface.Stop()

    def play(self):
        self.checkAndLoop()
        self.interface.Play()

    def pause(self):
        self.checkAndLoop()
        self.interface.Pause()

    def playPause(self):
        self.checkAndLoop()
        self.interface.PlayPause()

    def restart(self):
        self.checkAndLoop()
        current = str(self.getMetadata()['mpris:trackid'])
        self.openUri(current)

    def previous(self):
        self.checkAndLoop()
        current = str(self.getMetadata()['mpris:trackid'])
        self.interface.Previous()
        new = str(self.getMetadata()['mpris:trackid'])
        if (current == new):
            self.interface.Previous()

    def openUri(self,uri):
        self.checkAndLoop()
        if uri:
            self.interface.OpenUri(uri)
