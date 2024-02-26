// Package player ...
package player

// <---------------------------------------------------------------------------------------------------->

import (
	"time"

	"github.com/SkillpTm/SpotifyTrueRandomShuffle/internal/api"
	"github.com/SkillpTm/SpotifyTrueRandomShuffle/internal/util"
)

// <---------------------------------------------------------------------------------------------------->

const (
    baseURL = "https://api.spotify.com/v1/"
    userProfileExtension = "me"
    playbackStateExtension = "me/player"
)

var refreshTime = time.Duration(util.AppConfig.LoopRefreshTime) * time.Second

// <---------------------------------------------------------------------------------------------------->



// Start is our main loop which repeats infinitly and provides with all parts needed for TrueRandomShuffle
func Start() {

    // Get the user's profile for their country
    userProfile, err := util.MakeGETRequest(baseURL + userProfileExtension, api.UserToken.GetAccessToken())
    if err != nil {
        util.LogError(err)
    }

    // create player with country tag
    userPlayer := Player{userCountry: userProfile["country"].(string),}


    // the main program loop starts here
    for {
        time.Sleep(refreshTime)

        // check if our playlist still exists and is on the player
    }
}