import subprocess
import re
##TODO
#   -adding check if spotify is running, if not open it, every time

#the check does not work i think
class VolumeManagerLubuntu:
    def __init__(self):
        self.muted = False
        self.prevVol = self.get()
        #maybe get current vol
        return


        """
        Simple mixer control 'Master',0
  Capabilities: pvolume pswitch pswitch-joined
  Playback channels: Front Left - Front Right
  Limits: Playback 0 - 65536
  Mono:
  Front Left: Playback 52432 [80%] [on]
  Front Right: Playback 52432 [80%] [on]


NEED TO PARSE THIS TO GET CURRENT
  """
    def up(self,percentage=15):
        res = subprocess.check_output(['amixer','-D','pulse','sset','Master',str(percentage)+'%+'])
        m = re.search('\[(.+)%\]', res)
        return int(m.group(0)[1:-2])

    def down(self,percentage=15):
        res = subprocess.check_output(['amixer','-D','pulse','sset','Master',str(percentage)+'%-'])
        m = re.search('\[(.+)%\]', res)
        return int(m.group(0)[1:-2])

    def mute(self):
        #need to save on self.volume the current to unmute
        self.muted = True
        self.prevVol = self.get()
        return self.down(100)

    #
    def unmute(self):
        #need to save on self.volume the current to unmute

        if (self.get()==0) and self.muted:
            self.muted = False
            self.up(self.prevVol)
        return self.prevVol

    def get(self):
        res = subprocess.check_output(['amixer','-D','pulse','sset','Master','0%+'])
        m = re.search('\[(.+)%\]', res)
        return int(m.group(0)[1:-2])
