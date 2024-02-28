// Package player ...
package player

import (
	"errors"

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