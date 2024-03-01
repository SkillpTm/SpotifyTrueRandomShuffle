// Package player is responsible for handling the main loop of the program. It executes all functionality relevant to make TrueRandomShuffle work.
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
	// general values
	userCountry string
	userID      string

	// check values
	currentlyPlayingType string
	isPlaying            bool
	isPrivateSession     bool
	repeatState          string
	shuffleState         bool
	smartShuffle         bool

	// context values
	contextHREF string
	contextType string
	contextURI  string

	contextLength int

	// shuffle playlist values
	shufflePlaylistHREF      string
	shufflePlaylistLength    int
	shufflePlaylistTrackURIs []string
	shufflePlaylistURI       string
}

// newPlayer creates and returns a player with the userID and userCountry set
func newPlayer() (*Player, error) {
	player := Player{}

	// Get the user's profile for their country
	userProfile, err := util.MakeHTTPRequest("GET", baseURL+userProfileExtension, api.UserToken.GetAccessTokenHeader(), nil, nil)
	if err != nil {
		return &player, fmt.Errorf("couldn't GET user profile; %s", err.Error())
	}

	// create player with country tag
	player.userID = userProfile["id"].(string)
	player.userCountry = userProfile["country"].(string)

	return &player, nil
}

// getShufflePlaylist gets the href and uri of the shuffle playlist. If it needs to it will also create the playlist.
func (player *Player) setShufflePlaylist() error {
	// get the shuffle playlist values from the JSON
	shufflePlaylistMap, err := util.GetJSONData(util.AppConfig.ShufflePlaylistPath)
	if err != nil {
		return fmt.Errorf("couldn't get temp playlist json; %s", err.Error())
	}

	// did we get temp a playlist from the json?
	if shufflePlaylistMap["uri"].(string) != "" {
		player.shufflePlaylistHREF = shufflePlaylistMap["href"].(string)
		player.shufflePlaylistURI = shufflePlaylistMap["uri"].(string)

		// reset the temp playlist from a previous execution to avoid missmatching with our data
		err = player.clearShufflePlaylist()
		if err != nil {
			return fmt.Errorf("couldn't clear shuffle playlist; %s", err.Error())
		}

		return nil
	}

	href, uri, err := player.createShufflePlaylist()
	if err != nil {
		return fmt.Errorf("couldn't create shuffle playlist; %s", err.Error())
	}

	player.shufflePlaylistHREF = href
	player.shufflePlaylistURI = uri

	return nil
}

// clearShufflePlaylist ensours that on spotify's end the shuffle playlist is empty
func (player *Player) clearShufflePlaylist() error {
	var trackURIs []string

	if len(player.shufflePlaylistTrackURIs) > 0 {
		trackURIs = player.shufflePlaylistTrackURIs
		// maybe reset palyer var here??
	} else {
		shufflePlaylistTrackURIs, err := player.getShufflePlaylistTrackURIs()
		if err != nil {
			return fmt.Errorf("couldn't get all shuffle playlist tracks; %s", err.Error())
		}
		trackURIs = shufflePlaylistTrackURIs
	}

	// if we don't have any uris now the shuffle playlist is already empty
	if len(trackURIs) == 0 {
		return nil
	}

	// now clear the shuffle Playlist
	headers := api.UserToken.GetAccessTokenHeader()
	headers["Content-Type"] = "application/json"

	bodyData := map[string]interface{}{"tracks": []map[string]string{}}

	// populate body
	for _, uri := range trackURIs {
		bodyData["tracks"] = append(bodyData["tracks"].([]map[string]string), map[string]string{"uri": uri})
	}

	// make call to remove all tracks of the shuffle playlist
	_, err := util.MakeHTTPRequest("DELETE", player.shufflePlaylistHREF+"/tracks", headers, nil, bodyData)
	if err != nil {
		return fmt.Errorf("couldn't DELETE request all tracks from the temp playlist; %s", err.Error())
	}

	return nil
}

// getShufflePlaylistTrackURIs makes a call to Spotify and returns a slice of the URIs currently in the shuffle playlist
func (player *Player) getShufflePlaylistTrackURIs() ([]string, error) {
	// define the storage var early
	var trackURIs []string

	// get current shuffle playlist tracks
	tracksResponse, err := util.MakeHTTPRequest("GET", player.shufflePlaylistHREF+"/tracks", api.UserToken.GetAccessTokenHeader(), nil, nil)
	if err != nil {
		return trackURIs, fmt.Errorf("couldn't GET request random track from context; %s", err.Error())
	}

	// append the URIs to trackURIs
	for _, item := range tracksResponse["items"].([]interface{}) {
		trackURIs = append(trackURIs, item.(map[string]interface{})["track"].(map[string]interface{})["uri"].(string))
	}

	return trackURIs, nil
}

// createShufflePlaylist makes the shuffle playlist on spotify and returns it's href and uri
func (player *Player) createShufflePlaylist() (string, string, error) {
	// declare vars for storage early
	var shufflePlaylisthref string
	var shufflePlaylisturi string

	headers := api.UserToken.GetAccessTokenHeader()
	headers["Content-Type"] = "application/json"

	bodyData := map[string]interface{}{
		"name":        "TrueRandomShuffle",
		"description": "DON'T CHANGE ANYTHING in this playlist. This playlist was automatically created by SpotifyTrueRandomShuffle. You may remove it from your library.",
		"public":      false,
	}

	createPlaylistResponse, err := util.MakeHTTPRequest("POST", fmt.Sprintf("%susers/%s/playlists", baseURL, player.userID), headers, nil, bodyData)
	if err != nil {
		return shufflePlaylisthref, shufflePlaylisturi, fmt.Errorf("couldn't POST request create temp playlist; %s", err.Error())
	}

	// set the response values
	shufflePlaylisthref = createPlaylistResponse["href"].(string)
	shufflePlaylisturi = createPlaylistResponse["uri"].(string)

	// immeaditly remove the shuffle playlist from the user's library
	_, err = util.MakeHTTPRequest("DELETE", shufflePlaylisthref+"/followers", api.UserToken.GetAccessTokenHeader(), nil, nil)
	if err != nil {
		return shufflePlaylisthref, shufflePlaylisturi, fmt.Errorf("couldn't POST request create temp playlist; %s", err.Error())
	}

	// write to the shuffle playlist JSON for use after a restart
	err = util.WriteJSONData(
		util.AppConfig.ShufflePlaylistPath,
		map[string]interface{}{
			"href": shufflePlaylisthref,
			"uri":  shufflePlaylisturi,
		},
	)
	if err != nil {
		return shufflePlaylisthref, shufflePlaylisturi, fmt.Errorf("couldn't write shuffle playlist data to JSON; %s", err.Error())
	}

	return shufflePlaylisthref, shufflePlaylisturi, nil
}

// runChecks ensures that all check and context values are up to date and valid || false is an error
func (player *Player) runChecks(playbackResponse *map[string]interface{}) (bool, error) {
	// check if there is a playback state
	if len(*playbackResponse) == 0 {
		return false, nil
	}

	// set the default check values
	player.setCheckValues(playbackResponse)

	// run the default checks
	if !player.playbackChecks() {
		return false, nil
	}

	// check if a context exists (i.e. if the user is listening to a song outside of an album/playlist)
	if (*playbackResponse)["context"] == nil {
		return false, nil
	}

	err := player.updateContext(playbackResponse)
	if err != nil {
		return false, fmt.Errorf("couldn't update context; %s", err.Error())
	}

	// if thex context only has 1 track fail
	if player.contextLength == 1 {
		return false, nil
	}

	if !player.shuffleCheck((*playbackResponse)["context"].(map[string]interface{})["uri"].(string)) {
		return false, nil
	}

	return true, nil
}

// loadPlaybackOnPlayer receives the response from the API call for playback State and loads it onto the player
func (player *Player) setCheckValues(playbackResponse *map[string]interface{}) {
	// add all relevant values to the player
	player.currentlyPlayingType = (*playbackResponse)["currently_playing_type"].(string)
	player.isPlaying = (*playbackResponse)["is_playing"].(bool)
	player.isPrivateSession = (*playbackResponse)["device"].(map[string]interface{})["is_private_session"].(bool)
	player.repeatState = (*playbackResponse)["repeat_state"].(string)
	player.shuffleState = (*playbackResponse)["shuffle_state"].(bool)
	player.smartShuffle = (*playbackResponse)["smart_shuffle"].(bool)
}

// playbackChecks makes sure all remaining factors pass || false is an error
func (player *Player) playbackChecks() bool {
	// the playback response has to pass all of these checks
	return player.isPlaying && // is the user playing something right now
		!player.isPrivateSession && // is in a private session
		player.currentlyPlayingType == "track" && // is listening to a track
		player.repeatState != "track" && // doesn't have repeat on a track turned on
		player.contextType != "show" && // the context can't be a show
		player.contextType != "artist" // the context can't be an artist
}

// updateContext ensures the context on the player is always up to date, by potentially reseting and updating the context and the shuffle playlist
func (player *Player) updateContext(playbackResponse *map[string]interface{}) error {
	err := player.clearContext((*playbackResponse)["context"].(map[string]interface{})["uri"].(string))
	if err != nil {
		return fmt.Errorf("couldn't clear context on player and reset shuffle playlist; %s", err.Error())
	}
	// if the context isn't empty it's up to date
	if player.contextURI != "" {
		return nil
	}

	err = player.setContext(playbackResponse)
	if err != nil {
		return fmt.Errorf("couldn't set context; %s", err.Error())
	}

	return nil
}

// clearContext returns a bool on if it cleared the context and emptied the shuffle playlist
func (player *Player) clearContext(newContextURI string) error {
	// if there is no context, don't clear it
	if player.contextURI == "" {
		return nil
	}

	// check that the user is playing a context that isn't the shuffle playlist or the original context
	if newContextURI == player.shufflePlaylistURI ||
		newContextURI == player.contextURI {
		return nil
	}

	err := player.clearShufflePlaylist()
	if err != nil {
		return fmt.Errorf("couldn't clear shuffle playlist; %s", err.Error())
	}

	// clear the context values from the player
	player.contextHREF = ""
	player.contextType = ""
	player.contextURI = ""
	player.contextLength = 0
	player.shufflePlaylistTrackURIs = nil

	return nil
}

// setContext only sets a new context if the user is listening to anything but the shuffle playlist
func (player *Player) setContext(playbackResponse *map[string]interface{}) error {
	// check if our context is a new album/playlist and not the shuffle playlist
	if player.shufflePlaylistURI == (*playbackResponse)["context"].(map[string]interface{})["uri"].(string) {
		return nil
	}

	// set the new context
	player.contextHREF = (*playbackResponse)["context"].(map[string]interface{})["href"].(string)
	player.contextType = (*playbackResponse)["context"].(map[string]interface{})["type"].(string)
	player.contextURI = (*playbackResponse)["context"].(map[string]interface{})["uri"].(string)

	length, err := player.getContextLength()
	if err != nil {
		return fmt.Errorf("couldn't get context length; %s", err.Error())
	}
	player.contextLength = length

	// if the context length is bigger than our playlist size limit, limit the playlist length to the size
	if player.contextLength >= util.AppConfig.ShufflePlaylistSize {
		player.shufflePlaylistLength = util.AppConfig.ShufflePlaylistSize
	} else {
		player.shufflePlaylistLength = player.contextLength
	}

	return nil
}

func (player *Player) getContextLength() (int, error) {
	var length int

	contextResponse, err := util.MakeHTTPRequest("GET", fmt.Sprintf("%s?market=%s", player.contextHREF, player.userCountry), api.UserToken.GetAccessTokenHeader(), nil, nil)
	if err != nil {
		return length, fmt.Errorf("couldn't GET request context; %s", err.Error())
	}

	// depending on the context type the length is in another part of the JSON
	if player.contextType == "album" {
		length = int(contextResponse["total_tracks"].(float64))
	} else if player.contextType == "playlist" {
		length = int(contextResponse["tracks"].(map[string]interface{})["total"].(float64))
	}

	return length, nil
}

// shuffleCheck makes sure shuffle is turned one || false is an error
func (player *Player) shuffleCheck(playbackContextURI string) bool {
	// if we're playing the temp playlist the shuffle state doesn't matter
	if playbackContextURI == player.shufflePlaylistURI {
		return true
	}

	// otherwise return the actual shuffle state
	return player.shuffleState && !player.smartShuffle
}

// removeFinishedTracks removes all tracks in the playlist before the current one
func (player *Player) removeFinishedTracks(currentTrackURI string) error {
	// check if a song has been finished/skipped, while ignoring manually queued songs
	for index, shuffleTrackURI := range player.shufflePlaylistTrackURIs {
		// check where the currently playing track is in the shuffle playlist
		if currentTrackURI != shuffleTrackURI {
			continue
		}

		// if the track is the first track we don't need to remove anything
		if index == 0 {
			break
		}

		headers := api.UserToken.GetAccessTokenHeader()
		headers["Content-Type"] = "application/json"

		bodyData := map[string]interface{}{"tracks": []map[string]string{}}

		// populate body
		for _, uri := range player.shufflePlaylistTrackURIs {
			if uri == shuffleTrackURI {
				break
			}

			bodyData["tracks"] = append(bodyData["tracks"].([]map[string]string), map[string]string{"uri": uri})
		}

		_, err := util.MakeHTTPRequest("DELETE", player.shufflePlaylistHREF+"/tracks", headers, nil, bodyData)
		if err != nil {
			return fmt.Errorf("couldn't DELETE request remove playlist items; %s", err.Error())
		}

		// simply shorten shufflePlaylistTrackURIs by index to remove multiple URIs if needed
		player.shufflePlaylistTrackURIs = player.shufflePlaylistTrackURIs[index:]

		break
	}

	return nil
}

// fillShufflePlaylist fills the shuffle playlist up to the size of shufflePlaylistLength
func (player *Player) fillShufflePlaylist() error {
	// validate periodically that there is no missmatch between shufflePlaylistTrackURIs and the actualy tracks
	// this can only happen if the user manually removes a track from the playlist
	// on default settings this should happen about every ~2min
	if rand.Intn(player.shufflePlaylistLength*4) == 0 {
		// get the current tracks in the playlist
		currentShuffleTracks, err := player.getShufflePlaylistTrackURIs()
		if err != nil {
			return fmt.Errorf("couldn't get current shuffle tracks; %s", err.Error())
		}

		// check if we have a miss match between our var and spotify
		if len(currentShuffleTracks) != len(player.shufflePlaylistTrackURIs) {
			// set our var to spotify's values
			player.shufflePlaylistTrackURIs = currentShuffleTracks
		}
	}

	// if our shuffle playlist is already full just return
	if player.shufflePlaylistLength-len(player.shufflePlaylistTrackURIs) == 0 {
		return nil
	}

	// we use a map for easy access to figure out if a song is already in the palylist
	var currentTracks = map[string]bool{}
	// the slice stores only the URIs to be added to the shuffle playlist
	var toBeAddedTracks []string

	if len(player.shufflePlaylistTrackURIs) > 0 {
		for _, uri := range player.shufflePlaylistTrackURIs {
			currentTracks[uri] = true // the value we insert her doesn't matter, we just need a placeholder for the check later
		}
	}

	// loop to add song URIs to toBeAddedTracks
	for {
		randomTrackURL := fmt.Sprintf("%s/tracks?market=%s&limit=%d&offset=%d", player.contextHREF, player.userCountry, 20, rand.Intn(player.contextLength))

		randomTrackResponse, err := util.MakeHTTPRequest("GET", randomTrackURL, api.UserToken.GetAccessTokenHeader(), nil, nil)
		if err != nil {
			return fmt.Errorf("couldn't GET request random track from context; %s", err.Error())
		}

		randomTrackURI := ""
		var currentTrack string

		// loop over the response tracks
		for _, item := range randomTrackResponse["items"].([]interface{}) {
			// depending on if the context is an album or a playlist the URI is in a different position in the JSON
			if player.contextType == "album" {
				currentTrack = item.(map[string]interface{})["uri"].(string)
			} else if player.contextType == "playlist" {
				currentTrack = item.(map[string]interface{})["track"].(map[string]interface{})["uri"].(string)
			}

			if currentTracks[currentTrack] {
				continue
			}

			randomTrackURI = currentTrack
			break
		}

		if randomTrackURI == "" {
			continue
		}

		currentTracks[randomTrackURI] = true
		toBeAddedTracks = append(toBeAddedTracks, randomTrackURI)

		if len(currentTracks) == player.shufflePlaylistLength {
			break
		}
	}

	bodyData := map[string]interface{}{
		"uris": toBeAddedTracks,
	}

	// finally add all the URIs to the shuffle playlist
	_, err := util.MakeHTTPRequest("POST", player.shufflePlaylistHREF+"/tracks", api.UserToken.GetAccessTokenHeader(), nil, bodyData)
	if err != nil {
		return fmt.Errorf("couldn't POST request add in new items to temp playlist; %s", err.Error())
	}

	player.shufflePlaylistTrackURIs = append(player.shufflePlaylistTrackURIs, toBeAddedTracks...)

	return nil
}

func (player *Player) startShufflePlaylist() error {
	headers := api.UserToken.GetAccessTokenHeader()
	headers["Content-Type"] = "application/json"

	bodyData := map[string]interface{}{
		"context_uri": player.shufflePlaylistURI,
		"offset": map[string]int{
			"position": 0,
		},
	}

	_, err := util.MakeHTTPRequest("PUT", baseURL+startPlaybackExtension, headers, nil, bodyData)
	if err != nil {
		return fmt.Errorf("couldn't PUT request start playback; %s", err.Error())
	}

	// turn of shuffle for the temp playlist
	_, err = util.MakeHTTPRequest("PUT", baseURL+tooglePlaybackShuffleExtension+"?state=false", api.UserToken.GetAccessTokenHeader(), nil, nil)
	if err != nil {
		return fmt.Errorf("couldn't PUT request shuffle playback; %s", err.Error())
	}

	return nil
}
