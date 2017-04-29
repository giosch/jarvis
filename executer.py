import spotify.spotify as spotLib

class Executer:
    def __init__(self):
        spotify = spotLib.Spotify()
        self.commands = [{'words':['porno','immagini'],'action':self.magic},
                        #{'words':[],'action':},
                        {'words':['prossima','canzone'],'action':spotify.next},
                        {'words':['salta','canzone'],'action':spotify.next},
                        {'words':['cambia','canzone'],'action':spotify.next},
                        {'words':['precedente','canzone'],'action':spotify.previous},
                        {'words':['prima','canzone'],'action':spotify.previous},
                        {'words':['metti','quella','prima'],'action':spotify.previous},
                        {'words':['rimetti','quella','prima'],'action':spotify.previous},
                        {'words':['ferma','musica'],'action':spotify.stop},
                        {'words':['ferma','canzone'],'action':spotify.stop},
                        {'words':['stop','musica'],'action':spotify.stop},
                        {'words':['stop','canzone'],'action':spotify.stop},
                        {'words':['metti','musica'],'action':spotify.play},
                        {'words':['partire','musica'],'action':spotify.play}]
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
