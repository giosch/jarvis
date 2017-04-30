import plugins.spotify as spotLib
import plugins.volumeManager as volumeManager

class Executer:
    def __init__(self):
        self.spotify = spotLib.Spotify()
        self.volManager = volumeManager.VolumeManagerLubuntu()
        self.commands = [{'words':['porno','immagini'],'action':self.magic},
                        #{'words':[],'action':},
                        {'words':['prossima','canzone'],'action':self.spotify.next},
                        {'words':['salta','canzone'],'action':self.spotify.next},
                        {'words':['cambia','canzone'],'action':self.spotify.next},
                        {'words':['precedente','canzone'],'action':self.spotify.previous},
                        {'words':['prima','canzone'],'action':self.spotify.previous},
                        {'words':['metti','quella','prima'],'action':self.spotify.previous},
                        {'words':['rimetti','quella','prima'],'action':self.spotify.previous},
                        {'words':['ferma','musica'],'action':self.spotify.stop},
                        {'words':['ferma','canzone'],'action':self.spotify.stop},
                        {'words':['stop','musica'],'action':self.spotify.stop},
                        {'words':['stop','spotify'],'action':self.spotify.stop},
                        {'words':['stop','canzone'],'action':self.spotify.stop},
                        {'words':['metti','pausa','musica'],'action':self.spotify.pause},
                        {'words':['metti','pausa','spotify'],'action':self.spotify.pause},
                        {'words':['metti','musica'],'action':self.spotify.play},
                        {'words':['partire','musica'],'action':self.spotify.play},

                        {'words':['abbassa','volume'],'action':self.volManager.down},
                        {'words':['alza','volume'],'action':self.volManager.up},
                        {'words':['muta','musica'],'action':self.volManager.mute},
                        {'words':['muta','volume'],'action':self.volManager.mute},
                        {'words':['metti','muto'],'action':self.volManager.mute},
                        {'words':['rialza','musica'],'action':self.volManager.unmute},
                        {'words':['rialza','volume'],'action':self.volManager.unmute},
                        {'words':['togli','muto'],'action':self.volManager.unmute}]
    def do(self,s):
        #s = s.split(' ') #if active need to eliminate the \n at the end
        matches = []
        for command in self.commands:
            words_not_matched = sum([1 for x in command['words'] if x not in s])
            if words_not_matched == 0:
                matches.append(command)
        actions = set([x['action'] for x in matches])
        #policy : execute in order
        for i in actions:
            i()
        return True if len(actions) != 0 else False

    def magic(self):
        import webbrowser
        url = 'http://json-porn.com/demo/search/'
        chrome_path = '/usr/bin/google-chrome %s'
        webbrowser.get(chrome_path).open(url)
