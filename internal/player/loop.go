// Package player ...
package player

// <---------------------------------------------------------------------------------------------------->

import (
	"fmt"
	"time"

	"github.com/SkillpTm/SpotifyTrueRandomShuffle/internal/api"
	"github.com/SkillpTm/SpotifyTrueRandomShuffle/internal/util"
)

// <---------------------------------------------------------------------------------------------------->

const (
	baseURL = "https://api.spotify.com/v1/"
	getPlaylistExtension = "playlists/"
	playbackStateExtension = "me/player"
	startPlaybackExtension = "me/player/play"
	tooglePlaybackShuffleExtension = "me/player/shuffle"
	userProfileExtension = "me"
)

// <---------------------------------------------------------------------------------------------------->



// Start is our main loop which repeats infinitely and provides with all parts needed for TrueRandomShuffle
func Start() error {
	// make a new player
	userPlayer, err := newPlayer()
	if err != nil {
		return fmt.Errorf("couldn't create a new player; %s", err.Error())
	}

	// set the shuffle playlist (create it if necessary)
	err = userPlayer.setShufflePlaylist()
	if err != nil {
		return fmt.Errorf("couldn't get shuffle playlist; %s", err.Error())
	}

	// the main program loop starts here
	for {
		// slow down the loop so we don't get rate limited
		time.Sleep(time.Duration(int64(util.AppConfig.LoopRefreshTime * float64(time.Second))))

		// get the playback state for tests and context
        playbackResponse, err := util.MakeHTTPRequest("GET", baseURL + playbackStateExtension, api.UserToken.GetAccessTokenHeader(), nil, nil)
        if err != nil {
            return fmt.Errorf("couldn't GET request playback state; %s", err.Error())
        }

		// run general checks
		passed, err := userPlayer.runChecks(&playbackResponse)
        if err != nil {
            return fmt.Errorf("couldn't run all checks; %s", err.Error())
        }

		if !passed {
			continue
		}

		// remove tracks behind the currently playing one
		err = userPlayer.removeFinishedTracks(playbackResponse["item"].(map[string]interface{})["uri"].(string))
		if err != nil {
			return fmt.Errorf("couldn't remove finished tracks; %s", err.Error())
		}

		// fill shuffle playlist up until it has shufflePlaylistLength tracks
		err = userPlayer.fillShufflePlaylist()
		if err != nil {
			return fmt.Errorf("couldn't fill shuffle playlist; %s", err.Error())
		}

		// check if we're listening to the temp playlist right now
		if (playbackResponse["context"].(map[string]interface{})["uri"].(string) == userPlayer.shufflePlaylistURI) {
			continue
		}

		// start playing our palylist
		err = userPlayer.startShufflePlaylist()
		if err != nil {
			return fmt.Errorf("couldn't start playing temp playlist; %s", err.Error())
		}
	}
}