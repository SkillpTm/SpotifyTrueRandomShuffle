// Package player ...
package player

import (
	"errors"
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

	contextLength int // make something for this

	tempPlaylistHREF string
	tempPlaylistTrackURIs []string
	tempPlaylistURI string
}



func (player *Player) loadPlaybackOnPlayer(playbackResponse map[string]interface{}) {
	contextHREF := playbackResponse["context"].(map[string]interface{})["href"].(string)
	contextType := playbackResponse["context"].(map[string]interface{})["type"].(string)
	contextURI := playbackResponse["context"].(map[string]interface{})["uri"].(string)

	// check if our context is a new album/playlist and not the temp playlist. Only then set it as the context on the player.
	if (contextHREF != player.tempPlaylistHREF && contextURI != player.tempPlaylistURI) {
		player.contextHREF = contextHREF
		player.contextType = contextType
		player.contextURI = contextURI
	}

	player.currentlyPlayingType = playbackResponse["currently_playing_type"].(string)
	player.isPlaying = playbackResponse["is_playing"].(bool)
	player.isPrivateSession = playbackResponse["device"].(map[string]interface{})["is_private_session"].(bool)
	player.repeatState = playbackResponse["repeat_state"].(string)
	player.shuffleState = playbackResponse["shuffle_state"].(bool)
	player.smartShuffle = playbackResponse["smart_shuffle"].(bool)
}



// playbackChecks returns true if any of it's checks failed
func (player *Player) playbackChecks() bool {
	return	!player.isPlaying ||						// isn't playing something right now
			player.isPrivateSession ||					// is in a private session
			player.currentlyPlayingType != "track" ||	// isn't listening to a track
			player.repeatState == "track" ||			// does have repeat turned on for this track
			!player.shuffleState ||						// hasn't turned on shuffle
			player.smartShuffle ||						// has turned on smart shuffle
			player.contextType == "show" ||				// context is a show
			player.contextType == "artist"				// context is an artist
}



// validateTempPlaylist ensures that we have set the href, id and uri on the player. If needed it will also create and populate the playlist
func (player *Player) validateTempPlaylist() error {
	// do we already have a temp playlist on the player
	if (player.tempPlaylistHREF == "" ||
		player.tempPlaylistURI == "") {
		// check if the temp plalyist has teh required amount of tracks in it
		err := player.validateTempPlaylistTracks()
		if err != nil {
			return errors.New("couldn't validate temp playlist tracks: " + err.Error())
		}

		return nil
	}

	// get the temp playlist uri from the JSON
	tempPlaylistMap, err := util.GetJSONData(util.AppConfig.TempPlaylistPath)
	if err != nil {
		return errors.New("couldn't get temp playlist json: " + err.Error())
	}

	tempPlaylistHREF := tempPlaylistMap["href"].(string)
	tempPlaylistURI := tempPlaylistMap["uri"].(string)

	// did we get all temp playlist values from the json?
	if (player.tempPlaylistHREF == "" ||
		player.tempPlaylistURI == "") {
		player.tempPlaylistHREF = tempPlaylistHREF
		player.tempPlaylistURI = tempPlaylistURI
		return nil
	}

	// since there is no temp playlist we create one
	err = player.createTempPlaylist()
	if err != nil {
		return errors.New("couldn't create temp playlist: " + err.Error())
	}

	// it's impossible to actually delete a playlist. We're unfollowing it, so the user won't mess with it
	_, err = util.MakeHTTPRequest("DELETE", player.tempPlaylistHREF + "/followers", api.UserToken.GetAccessTokenHeader(), nil, nil)
	if err != nil {
		return errors.New("couldn't DELETE request temp playlist: " + err.Error())
	}

	// since the temp playlist just got created populate it
	err = player.populateTempPlaylist(util.AppConfig.TempPlaylistSize)
	if err != nil {
		return errors.New("couldn't populate temp playlist: " + err.Error())
	}

	return nil
}



// validateTempPlaylistTracks checks if the temp playlist has the config specified amount of tracks in it. If that's not the case it populates it.
func (player *Player) validateTempPlaylistTracks() error {
	playlistTracksResponse, err := util.MakeHTTPRequest("GET", player.tempPlaylistHREF + "/tracks", api.UserToken.GetAccessTokenHeader(), nil, nil)
	if err != nil {
		return errors.New("couldn't GET request temp playlist tracks: " + err.Error())
	}

	playlistTracksLength := int(playlistTracksResponse["total"].(float64))

	// add the compared amount of songs to the temp playlist (this value can be zero)
	err = player.populateTempPlaylist(util.AppConfig.TempPlaylistSize - playlistTracksLength)
	if err != nil {
		return errors.New("couldn't populate temp playlist: " + err.Error())
	}

	return nil
}



// createTempPlaylist creates the temp playlist needed for the main loop and sets temp playlist values to the player.
func (player *Player) createTempPlaylist() error {

	bodyData := map[string]interface{}{
		"name": "TrueRandomShuffle",
		"desciption": "DO NOT add ANYTHING to this playlist. This playlist was automatically create by SpotifyTrueRandomShuffle.",
		"public": false,
	}

	createPlaylistResponse, err := util.MakeHTTPRequest("POST", baseURL + createPlaylistExtension, api.UserToken.GetAccessTokenHeader(), nil, bodyData)
	if err != nil {
		return errors.New("couldn't POST request create temp playlist: " + err.Error())
	}

	err = util.WriteJSONData(
		util.AppConfig.TempPlaylistPath, 
		map[string]interface{}{
			"href" : createPlaylistResponse["href"].(string),
			"uri" : createPlaylistResponse["uri"].(string),
		},
	)
	if err != nil {
		return errors.New("couldn't write temp playlist data to JSOON: " + err.Error())
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
			return errors.New("couldn't GET request random track from context: " + err.Error())
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
		return errors.New("couldn't POST request add in new items to temp playlist: " + err.Error())
	}

	player.tempPlaylistTrackURIs = append(player.tempPlaylistTrackURIs, newTracks...)

	return nil
}


func (player *Player) setContextLength() error {
	contextResponse, err := util.MakeHTTPRequest("GET", fmt.Sprintf("%s?market=%s", player.contextHREF, player.userCountry), api.UserToken.GetAccessTokenHeader(), nil, nil)
	if err != nil {
		return errors.New("couldn't GET request context type: " + err.Error())
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
	startPlaybackHeaders["Content-Type:"] = "application/json"

	startPlaybackData := map[string]interface{}{"context_uri" : player.tempPlaylistURI}

	_, err := util.MakeHTTPRequest("PUT", baseURL + startPlaybackExtension, startPlaybackHeaders, nil, startPlaybackData)
	if (err != nil) {
		return errors.New("couldn't PUT request start playback: " + err.Error())
	}

	// turn of shuffle for the temp playlist
	_, err = util.MakeHTTPRequest("PUT", baseURL + tooglePlaybackShuffleExtension + "?state=false", api.UserToken.GetAccessTokenHeader(), nil, nil)
	if (err != nil) {
		return errors.New("couldn't PUT request start playback: " + err.Error())
	}

	return nil
}