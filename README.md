## TrueRandomShuffle:
Spotify's shuffle is not random, they use an algorithm based on your listening behaviour. TrueRandomShuffle inserts <b>random</b> songs into your queue from the album/playlist you're listening to, to solve this problem.

## Usage:
Turing on shuffle (not smart shuffle) also turns on TrueRandomShuffle, which means the bot will add a random song from the album/playlist you're listening to, to your queue. Due to the limitations of the Spotify Web API the latest song added by the queue remains in the queue even after turning off shuffle or when changing albums/playlists.

TrueRandomShuffle is turned off during private sessions. It also doesn't effect you repeating a track.

<sup>note: The currently playing track and the next track in queue *can* be the same.</sup>

## How to setup TrueRandomShuffle:

### Make an app on `Spotify for Developers`:
Create an app here: https://developer.spotify.com/dashboard

### Clone this repository:
```bash
git clone https://github.com/SkillpTm/SpotifyTrueRandomShuffle
```

### Make a .env:
Make a .env at ./.env with this format:
```sh
SPOTIFY_ID = "[INSERT ID HERE]"
SPOTIFY_SECRET = "[INSERT SECRET HERE]"
SPOTIFY_REDIRECT_DOMAIN = "[INSERT DOMAIN HERE]" # example: "http://localhost"
```

### Compile TrueRandomShuffle:
```bash
go build -o ./SpotifyTrueRandomShuffle.exe ./cmd/main/
```

### Launch TrueRandomShuffle:
windows:
```bash
start ./SpotifyTrueRandomShuffle.exe
```
linux:
```bash
./SpotifyTrueRandomShuffle.exe
```

### Open the link:
The program will print a link to your CLI that you'll need to open and agree to. Afterwards TrueRandomShuffle is turned on.