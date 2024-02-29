## TrueRandomShuffle:
Spotify's shuffle is not random, they use an algorithm based on your listening behaviour. TrueRandomShuffle creates an hidden playlist, which it uses to randomoize your queue. It's completely frictionless, just press shuffle play on your playlist and let TrueRandomShuffle do the rest.

## Usage:
Turing on shuffle (not smart shuffle) also turns on TrueRandomShuffle, which means you'll be redirected to a hidden playlist. This playlist shouldn't be modififed, you can skip (even multiple songs) as usual.

TrueRandomShuffle is turned off during private sessions. It also doesn't effect you repeating a track or when you use smart shuffle.

### Customization:

You may edited the following values in ./configs/config.json:
- loopRefreshTime: This changes how often the main loop repeats itself (setting this too low may cause rate limiting from Spotify, which will stop TrueRandomShuffle).
- requestAuthEveryTime: This changes if you have to click "accept" in the browser for every restart.
- shufflePlaylistSize: This changes how big the hidden playlist for TrueRandomShuffle will be. Making it too big may occur rate limiting and more delay between the loops. Making it too small may mean TrueRandomShuffle can't refill the hidden playlist fast enough, if you spam skip.

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
The program will print a link to your CLI, which you'll need to open and you'll need to agree to the authorization notice. Afterwards TrueRandomShuffle is turned on.