// Package player ...
package player

// <---------------------------------------------------------------------------------------------------->

import (
	"log"
	"time"

	"github.com/SkillpTm/SpotifyTrueRandomShuffle/internal/api"
	"github.com/SkillpTm/SpotifyTrueRandomShuffle/internal/util"
)

// <---------------------------------------------------------------------------------------------------->

const playbackStateURL = "https://api.spotify.com/v1/me/player"

// <---------------------------------------------------------------------------------------------------->

func Start() {

	for {
		time.Sleep(1 * time.Second)

		responseMap, err := util.MakeGETRequest(playbackStateURL, api.UserToken.GetAccessToken())
		if err != nil {
			log.Fatal(err)
		}

		if (responseMap["context"] == nil) {
			continue
		}

		contextData := Player{
			isPlaying: responseMap["is_playing"].(bool),
			currentlyPlayingType: responseMap["currently_playing_type"].(string),
			repeatState: responseMap["repeat_state"].(string),
			shuffleState: responseMap["shuffle_state"].(bool),
			smartShuffle: responseMap["smart_shuffle"].(bool),
			contextType: responseMap["context"].(map[string]interface{})["type"].(string),
			contextHREF: responseMap["context"].(map[string]interface{})["href"].(string),
		}

		// continue if any of the checks fail
		if (contextData.RunChecks()) {
			continue
		}


	}
}