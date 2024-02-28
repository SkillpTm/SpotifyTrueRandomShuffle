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

	contextHREF string
	contextLength int
	contextType string

	tempPlaylistHREF string
	tempPlaylistID string
	tempPlaylistTrackURIs []string
	tempPlaylistURI string

	userID string
}



// validateTempPlaylist ensures that we have set the href, id and uri on the player. If needed it will also create and populate the playlist
func (player *Player) validateTempPlaylist() error {
	// do we already have a temp playlist on the player
	if (player.tempPlaylistID != "" &&
		player.tempPlaylistHREF != "" &&
		player.tempPlaylistURI != "") {
		return nil
	}

	// get the temp playlist uri from the JSON
	tempPlaylistMap, err := util.GetJSONData(util.AppConfig.TempPlaylistPath)
	if err != nil {
		return errors.New("couldn't get temp playlist json: " + err.Error())
	}

	tempPlaylistHREF := tempPlaylistMap["href"].(string)
	tempPlaylistID := tempPlaylistMap["id"].(string)
	tempPlaylistURI := tempPlaylistMap["uri"].(string)

	// did we get all temp playlist values from the json?
	if (player.tempPlaylistID != "" &&
		player.tempPlaylistHREF != "" &&
		player.tempPlaylistURI != "") {
		player.tempPlaylistHREF = tempPlaylistHREF
		player.tempPlaylistID = tempPlaylistID
		player.tempPlaylistURI = tempPlaylistURI
		return nil
	}

	// since there is no temp playlist we create one
	err = player.createTempPlaylist()
	if err != nil {
		return errors.New("couldn't create temp playlist: " + err.Error())
	}

	// since the temp playlist just got created populate it
	err = player.populateTempPlaylist(util.AppConfig.TempPlaylistSize)
	if err != nil {
		return errors.New("couldn't populate temp playlist: " + err.Error())
	}

	return nil
}



// createTempPlaylist creates the temp playlist needed for the main loop and sets temp playlist values to the player 
func (player *Player) createTempPlaylist() error {
	headers := map[string]string{
		"Authorization": "Bearer " + api.UserToken.GetAccessToken(),
	}

	bodyData := map[string]interface{}{
		"name": "TrueRandomShuffle",
		"desciption": "This playlist was automatically create by SpotifyTrueRandomShuffle. Please do not delete it unless you stopped using the program. You may move this playlist in any folder.",
		"public": false,
	}

	createPlaylistResponse, err := util.MakePOSTRequest(baseURL + createPlaylistExtension, headers, map[string]string{}, bodyData)
	if err != nil {
		return errors.New("couldn't POST request create temp playlist: " + err.Error())
	}

	player.tempPlaylistHREF = createPlaylistResponse["href"].(string)
	player.tempPlaylistID = createPlaylistResponse["id"].(string)
	player.tempPlaylistURI = createPlaylistResponse["uri"].(string)
	return nil
}



// populateTempPlaylist adds the amount fo missing songs to the temp playlist
func (player *Player) populateTempPlaylist(missingSongs int) error {
	var newTracks []string

	// loop to add song URIs to newTracks
	for _ = range missingSongs {
		randomTrackURL := fmt.Sprintf("%s/tracks?market=%s&limit=%d&offset=%d", player.contextHREF, player.userCountry, 1, rand.Intn(player.contextLength))

		randomTrackResponse, err := util.MakeGETRequest(randomTrackURL, api.UserToken.GetAccessToken())
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

	headers := map[string]string{
		"Authorization": "Bearer " + api.UserToken.GetAccessToken(),
	}

	bodyData := map[string]interface{}{
		"uris": newTracks,
	}

	// add songs to temp playlist
	_, err := util.MakePOSTRequest(player.tempPlaylistHREF + "/tracks", headers, map[string]string{}, bodyData)
	if err != nil {
		return errors.New("couldn't POST request add in new items to temp playlist: " + err.Error())
	}

	player.tempPlaylistTrackURIs = append(player.tempPlaylistTrackURIs, newTracks...)

	return nil
}