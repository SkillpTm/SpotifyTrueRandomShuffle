// Package player ...
package player

import (
	"fmt"
	"math/rand"

	"github.com/SkillpTm/SpotifyTrueRandomShuffle/internal/api"
	"github.com/SkillpTm/SpotifyTrueRandomShuffle/internal/util"
)

// <---------------------------------------------------------------------------------------------------->

// Player holds all the data and functions relevant for the main loop
type Player struct {
	userCountry string
	userID string

	contextHREF string
	contextType string
	contextURI string
	currentlyPlayingType string
	isPlaying bool
	isPrivateSession bool
	repeatState string
	shuffleState bool
	smartShuffle bool

	contextLength int

	tempPlaylistHREF string
	tempPlaylistTrackURIs []string
	tempPlaylistURI string
}



// loadPlaybackOnPlayer receives the response from the API call for playback State and loads it onto the player
func (player *Player) loadPlaybackOnPlayer(playbackResponse map[string]interface{}) {
	// add all relevant values to the player
	player.currentlyPlayingType = playbackResponse["currently_playing_type"].(string)
	player.isPlaying = playbackResponse["is_playing"].(bool)
	player.isPrivateSession = playbackResponse["device"].(map[string]interface{})["is_private_session"].(bool)
	player.repeatState = playbackResponse["repeat_state"].(string)
	player.shuffleState = playbackResponse["shuffle_state"].(bool)
	player.smartShuffle = playbackResponse["smart_shuffle"].(bool)

	// check if our context is a new album/playlist and not the temp playlist. Only then set it as the context on the player.
	if (player.tempPlaylistURI != playbackResponse["context"].(map[string]interface{})["uri"].(string)) {
		player.contextHREF = playbackResponse["context"].(map[string]interface{})["href"].(string)
		player.contextType = playbackResponse["context"].(map[string]interface{})["type"].(string)
		player.contextURI = playbackResponse["context"].(map[string]interface{})["uri"].(string)
	}
}



// shuffleCheck makes sure shuffle is turned one
func (player *Player) shuffleCheck(playbackContextURI string) bool {
	// if we're playing the temp playlist the shuffle state doesn't matter
	if playbackContextURI == player.tempPlaylistURI {
		return false
	}

	// otherwise return the actual shuffle state
	return player.shuffleState
}



// playbackChecks makes sure all remaing factors pass
func (player *Player) playbackChecks() bool {
	// the playback response has to pass all of these checks
	return	player.isPlaying &&							// is the user playing something right now
			!player.isPrivateSession &&					// is in a private session
			player.currentlyPlayingType == "track" &&	// is listening to a track
			player.repeatState != "track" &&			// doesn't have repeat on a track turned on
			!player.smartShuffle &&						// has smart shuffle turned off
			player.contextType != "show" &&				// the context can't be a show
			player.contextType != "artist"				// the context can't be an artist
}



// validateTempPlaylist ensures that we have set the href, id and uri on the player. If needed it will also create and populate the playlist
func (player *Player) validateTempPlaylist() error {
	// do we already have a temp playlist on the player
	if (player.tempPlaylistURI != "") {
		// check if the temp plalyist has the required amount of tracks in it
		err := player.validateTempPlaylistTracks()
		if err != nil {
			return fmt.Errorf("couldn't validate temp playlist tracks; %s", err.Error())
		}

		return nil
	}

	// get the temp playlist uri from the JSON
	tempPlaylistMap, err := util.GetJSONData(util.AppConfig.TempPlaylistPath)
	if err != nil {
		return fmt.Errorf("couldn't get temp playlist json; %s", err.Error())
	}

	tempPlaylistHREF := tempPlaylistMap["href"].(string)
	tempPlaylistURI := tempPlaylistMap["uri"].(string)

	// did we get all temp playlist values from the json?
	if (tempPlaylistURI != "") {
		player.tempPlaylistHREF = tempPlaylistHREF
		player.tempPlaylistURI = tempPlaylistURI

		// reset the temp playlist from a previous exectuion to avoid missmatching with our data
		err = player.resetTempPlaylist()
		if err != nil {
			return fmt.Errorf("couldn't reset temp playlist; %s", err.Error())
		}

		return nil
	}

	// since there is no temp playlist we create one
	err = player.createTempPlaylist()
	if err != nil {
		return fmt.Errorf("couldn't create temp playlist; %s", err.Error())
	}

	// it's impossible to actually delete a playlist. We're unfollowing it, so the user won't mess with it
	_, err = util.MakeHTTPRequest("DELETE", player.tempPlaylistHREF + "/followers", api.UserToken.GetAccessTokenHeader(), nil, nil)
	if err != nil {
		return fmt.Errorf("couldn't DELETE request temp playlist; %s", err.Error())
	}

	// since the temp playlist just got created populate it
	err = player.populateTempPlaylist(util.AppConfig.TempPlaylistSize)
	if err != nil {
		return fmt.Errorf("couldn't populate temp playlist; %s", err.Error())
	}

	return nil
}


// resetTempPlaylist clears the temp playlist, this is only required when restarting the application
func (player *Player) resetTempPlaylist() error {
	// get temp playlist tracks
	tempPlaylistTracks, err := util.MakeHTTPRequest("GET", player.tempPlaylistHREF + "/tracks",  api.UserToken.GetAccessTokenHeader(), nil, nil)
	if err != nil {
		return fmt.Errorf("couldn't GET request random track from context; %s", err.Error())
	}

	// return early if the temp playlist is already empty
	if (len(tempPlaylistTracks["items"].([]interface{})) == 0) {
		// re-populate the playlist
		err = player.populateTempPlaylist(util.AppConfig.TempPlaylistSize)
		if err != nil {
			return fmt.Errorf("couldn't populate temp playlist; %s", err.Error())
		}

		return nil
	}

	resetTempPlaylistHeaders := api.UserToken.GetAccessTokenHeader()
	resetTempPlaylistHeaders["Content-Type"] = "application/json"

	var resetTempPlaylistURIs []map[string]string

	for _, item := range tempPlaylistTracks["items"].([]interface{}) {
		resetTempPlaylistURIs = append(resetTempPlaylistURIs, map[string]string{
			"uri" : item.(map[string]interface{})["track"].(map[string]interface{})["uri"].(string),
		})
	}

	resetTempPlaylistBody := map[string]interface{}{
		"tracks" : resetTempPlaylistURIs,
	}

	// delete temp playlist tracks
	_, err = util.MakeHTTPRequest("DELETE", player.tempPlaylistHREF + "/tracks", resetTempPlaylistHeaders, nil, resetTempPlaylistBody)
	if err != nil {
		return fmt.Errorf("couldn't DELETE request all tracks from the temp playlist; %s", err.Error())
	}

	// re-populate the playlist
	err = player.populateTempPlaylist(util.AppConfig.TempPlaylistSize)
	if err != nil {
		return fmt.Errorf("couldn't populate temp playlist; %s", err.Error())
	}

	return nil
}



// validateTempPlaylistTracks checks if the temp playlist has the config specified amount of tracks in it. If that's not the case it populates it.
func (player *Player) validateTempPlaylistTracks() error {
	playlistTracksResponse, err := util.MakeHTTPRequest("GET", player.tempPlaylistHREF + "/tracks", api.UserToken.GetAccessTokenHeader(), nil, nil)
	if err != nil {
		return fmt.Errorf("couldn't GET request temp playlist tracks; %s", err.Error())
	}

	playlistTracksLength := int(playlistTracksResponse["total"].(float64))

	// add the compared amount of songs to the temp playlist (this value can be zero)
	err = player.populateTempPlaylist(util.AppConfig.TempPlaylistSize - playlistTracksLength)
	if err != nil {
		return fmt.Errorf("couldn't populate temp playlist; %s", err.Error())
	}

	return nil
}



// createTempPlaylist creates the temp playlist needed for the main loop and sets temp playlist values to the player.
func (player *Player) createTempPlaylist() error {
	createPlaylistHeaders := api.UserToken.GetAccessTokenHeader()
	createPlaylistHeaders["Content-Type"] = "application/json"

	bodyData := map[string]interface{}{
		"name": "TrueRandomShuffle",
		"desciption": "DO NOT add ANYTHING to this playlist. This playlist was automatically create by SpotifyTrueRandomShuffle.",
		"public": false,
	}

	createPlaylistResponse, err := util.MakeHTTPRequest("POST", baseURL + createPlaylistExtension, createPlaylistHeaders, nil, bodyData)
	if err != nil {
		return fmt.Errorf("couldn't POST request create temp playlist; %s", err.Error())
	}

	err = util.WriteJSONData(
		util.AppConfig.TempPlaylistPath, 
		map[string]interface{}{
			"href" : createPlaylistResponse["href"].(string),
			"uri" : createPlaylistResponse["uri"].(string),
		},
	)
	if err != nil {
		return fmt.Errorf("couldn't write temp playlist data to JSON; %s", err.Error())
	}

	player.tempPlaylistHREF = createPlaylistResponse["href"].(string)
	player.tempPlaylistURI = createPlaylistResponse["uri"].(string)
	return nil
}



// populateTempPlaylist adds the amount fo missing songs to the temp playlist
func (player *Player) populateTempPlaylist(missingSongs int) error {
	if (missingSongs == 0) {
		return nil
	}

	var newTracks []string

	// loop to add song URIs to newTracks
	for range missingSongs {
		randomTrackURL := fmt.Sprintf("%s/tracks?market=%s&limit=%d&offset=%d", player.contextHREF, player.userCountry, 1, rand.Intn(player.contextLength))

		randomTrackResponse, err := util.MakeHTTPRequest("GET", randomTrackURL,  api.UserToken.GetAccessTokenHeader(), nil, nil)
		if err != nil {
			return fmt.Errorf("couldn't GET request random track from context; %s", err.Error())
		}

		randomTrackURI := ""

		// depending on if the context is an album or a playlist the URI is in a different position in the JSON
		if (player.contextType == "album") {
			randomTrackURI = randomTrackResponse["items"].([]interface{})[0].(map[string]interface{})["uri"].(string)
		} else if (player.contextType == "playlist") {
			randomTrackURI = randomTrackResponse["items"].([]interface{})[0].(map[string]interface{})["track"].(map[string]interface{})["uri"].(string)
		}

		newTracks = append(newTracks, randomTrackURI)
	}

	bodyData := map[string]interface{}{
		"uris": newTracks,
	}

	// add songs to temp playlist
	_, err := util.MakeHTTPRequest( "POST", player.tempPlaylistHREF + "/tracks", api.UserToken.GetAccessTokenHeader(), nil, bodyData)
	if err != nil {
		return fmt.Errorf("couldn't POST request add in new items to temp playlist; %s", err.Error())
	}

	player.tempPlaylistTrackURIs = append(player.tempPlaylistTrackURIs, newTracks...)

	return nil
}



// resetContext resets all context values on the player, as weel as tempPlaylistTrackURIs and it delete all tracks from the temp playlist
func (player *Player) resetContext() error {
	resetContextHeaders := api.UserToken.GetAccessTokenHeader()
	resetContextHeaders["Content-Type"] = "application/json"

	var resetConextURIs []map[string]string

	for _, tempPlaylistTrackURI := range player.tempPlaylistTrackURIs {
		resetConextURIs = append(resetConextURIs, map[string]string{
			"uri" : tempPlaylistTrackURI,
		})
	}

	resetContextBody := map[string]interface{}{
		"tracks" : resetConextURIs,
	}

	_, err := util.MakeHTTPRequest("DELETE", player.tempPlaylistHREF + "/tracks", resetContextHeaders, nil, resetContextBody)
	if err != nil {
		return fmt.Errorf("couldn't DELETE request all tracks from the temp playlist; %s", err.Error())
	}

	// clear all context values and tempPlaylistTrackURIs so they can be set for the new context
	player.contextHREF = ""
	player.contextType = ""
	player.contextURI = ""
	player.contextLength = 0
	player.tempPlaylistTrackURIs = nil

	return nil
}



// setContextLength sets the contextLength on the player
func (player *Player) setContextLength() error {
	contextResponse, err := util.MakeHTTPRequest("GET", fmt.Sprintf("%s?market=%s", player.contextHREF, player.userCountry), api.UserToken.GetAccessTokenHeader(), nil, nil)
	if err != nil {
		return fmt.Errorf("couldn't GET request context type; %s", err.Error())
	}

	// depending on the context type the length is in another part of the JSON
	if (player.contextType == "album") {
		player.contextLength = int(contextResponse["total_tracks"].(float64))
	} else if (player.contextType == "playlist") {
		player.contextLength = int(contextResponse["tracks"].(map[string]interface{})["total"].(float64))
	}

	return nil
}



// playTempPlaylist starts playing our temp playlista nd turns shuffle on it off
func (player *Player) playTempPlaylist() error {
	startPlaybackHeaders := api.UserToken.GetAccessTokenHeader()
	startPlaybackHeaders["Content-Type"] = "application/json"

	startPlaybackData := map[string]interface{}{"context_uri" : player.tempPlaylistURI}

	_, err := util.MakeHTTPRequest("PUT", baseURL + startPlaybackExtension, startPlaybackHeaders, nil, startPlaybackData)
	if err != nil {
		return fmt.Errorf("couldn't PUT request start playback; %s", err.Error())
	}

	// turn of shuffle for the temp playlist
	_, err = util.MakeHTTPRequest("PUT", baseURL + tooglePlaybackShuffleExtension + "?state=false", api.UserToken.GetAccessTokenHeader(), nil, nil)
	if err != nil {
		return fmt.Errorf("couldn't PUT request shuffle playback; %s", err.Error())
	}

	return nil
}