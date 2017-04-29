import dbus
import subprocess
import time
##TODO
#   -adding check if spotify is running, if not open it, every time

#the check does not work i think
class Spotify:
    def __init__(self):
        try:
            self.bus = dbus.SessionBus()
            self.proxy = self.bus.get_object('org.mpris.MediaPlayer2.spotify', '/org/mpris/MediaPlayer2')
            self.interface = dbus.Interface(self.proxy, dbus_interface='org.mpris.MediaPlayer2.Player')
            self.properties = dbus.Interface(self.proxy, dbus_interface='org.freedesktop.DBus.Properties')
        except dbus.exceptions.DBusException as e:
            if "The name org.mpris.MediaPlayer2.spotify was not provided by any .service files" in e:
                #spotify is close
                subprocess.Popen(["spotify"])
                time.sleep(5)#to let it open
                self.__init__()

    def getProperties(self):
        return self.properties.GetAll('org.mpris.MediaPlayer2.Player').items()
    def getMetadata(self):
        prop = self.getProperties()
        for i in prop:
            if i[0] == 'Metadata':
                return i[1]
        return None

    def next(self):
        self.interface.Next()

    def stop(self):
        self.interface.Stop()

    def play(self):
        self.interface.Play()

    def pause(self):
        self.interface.Pause()

    def playPause(self):
        self.interface.PlayPause()

    ##not working
    def restart(self):
        current = str(self.getMetadata()['mpris:trackid'])
        self.openUri(current)

    def previous(self):
        current = str(self.getMetadata()['mpris:trackid'])
        self.interface.Previous()
        new = str(self.getMetadata()['mpris:trackid'])
        if (current == new):
            self.interface.Previous()

    def openUri(self,uri):
        if uri:
            self.interface.OpenUri(uri)
