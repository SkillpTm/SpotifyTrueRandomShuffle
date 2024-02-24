// Package player ...
package player

// <---------------------------------------------------------------------------------------------------->

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/SkillpTm/SpotifyTrueRandomShuffle/internal/api"
	"github.com/SkillpTm/SpotifyTrueRandomShuffle/internal/util"
)

// <---------------------------------------------------------------------------------------------------->

const (
	refreshTime = 1 * time.Second
	userProfileURL = "https://api.spotify.com/v1/me"
	playbackStateURL = "https://api.spotify.com/v1/me/player"
	queueURL = "https://api.spotify.com/v1/me/player/queue"
	addToQueueURL = "https://api.spotify.com/v1/me/player/queue"
)

// <---------------------------------------------------------------------------------------------------->



func Start() {

	userProfile, err := util.MakeGETRequest(userProfileURL, api.UserToken.GetAccessToken())
	if err != nil {
		log.Fatal(err)
	}

	contextData := Player{country: userProfile["country"].(string),}


	for {
		time.Sleep(refreshTime)

		playbackResponse, err := util.MakeGETRequest(playbackStateURL, api.UserToken.GetAccessToken())
		if err != nil {
			log.Fatal(err)
		}

		if (playbackResponse["context"] == nil) {
			continue
		}

		contextData.isPlaying = playbackResponse["is_playing"].(bool)
		contextData.isPrivateSession = playbackResponse["device"].(map[string]interface{})["is_private_session"].(bool)
		contextData.currentlyPlayingType = playbackResponse["currently_playing_type"].(string)
		contextData.repeatState = playbackResponse["repeat_state"].(string)
		contextData.shuffleState = playbackResponse["shuffle_state"].(bool)
		contextData.smartShuffle = playbackResponse["smart_shuffle"].(bool)
		contextData.contextType = playbackResponse["context"].(map[string]interface{})["type"].(string)
		contextData.contextHREF = playbackResponse["context"].(map[string]interface{})["href"].(string)


		// continue if any of the checks fail
		if (contextData.RunChecks()) {
			continue
		}

		queueResponse, err := util.MakeGETRequest(queueURL, api.UserToken.GetAccessToken())
		if err != nil {
			log.Fatal(err)
		}

		// check if the next song in the queue is the song we last added
		if (contextData.lastQueueSongURI == queueResponse["queue"].([]interface{})[0].(map[string]interface{})["uri"]) {
			continue
		}

		contextResponse, err := util.MakeGETRequest(contextData.contextHREF, api.UserToken.GetAccessToken())
		if err != nil {
			log.Fatal(err)
		}

		length := 0

		if (contextData.contextType == "album") {
			length = int(contextResponse["total_tracks"].(float64))
		} else if (contextData.contextType == "playlist") {
			length = int(contextResponse["tracks"].(map[string]interface{})["total"].(float64))
		}

		randomTrackURL := contextData.contextHREF + fmt.Sprintf("/tracks?market=%s&limit=%d&offset=%d", contextData.country, 1, rand.Intn(length))

		randomTrackResponse, err := util.MakeGETRequest(randomTrackURL, api.UserToken.GetAccessToken())
		if err != nil {
			log.Fatal(err)
		}

		randomTrackURI := ""

		if (contextData.contextType == "album") {
			randomTrackURI = randomTrackResponse["items"].([]interface{})[0].(map[string]interface{})["uri"].(string)
		} else if (contextData.contextType == "playlist") {
			randomTrackURI = randomTrackResponse["items"].([]interface{})[0].(map[string]interface{})["track"].(map[string]interface{})["uri"].(string)
		}

		_, err = util.MakePOSTRequest(fmt.Sprintf("%s?uri=%s", addToQueueURL, randomTrackURI), map[string]string{}, map[string]string{"Authorization": "Bearer " + api.UserToken.GetAccessToken(),})
		if err != nil {
			log.Fatal(err)
		}

		contextData.lastQueueSongURI = randomTrackURI
	}
}