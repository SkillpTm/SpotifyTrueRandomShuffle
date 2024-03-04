package main

// <---------------------------------------------------------------------------------------------------->

import (
	"fmt"
	"log"
	"strings"

	"github.com/SkillpTm/SpotifyTrueRandomShuffle/internal/auth"
	"github.com/SkillpTm/SpotifyTrueRandomShuffle/internal/player"
	"github.com/SkillpTm/SpotifyTrueRandomShuffle/internal/util"
)

// <---------------------------------------------------------------------------------------------------->

// main is the entry into the program, it sets up our config, envs, authorizes with the user and starts the main loop
func main() {
	err := util.Setup()
	if err != nil {
		log.Fatal(fmt.Errorf("couldn't setup config; %s", err.Error()))
	}

	auth.User()

	for {
		err = player.Start()
		// we can't return nil, so we don't error check
		util.LogError(fmt.Errorf("couldn't continue main loop; %s", err.Error()), false)

		// check if Spotify terminated our connection
		if strings.Contains(err.Error(), "connection reset by peer") {
			// forcefully refresh our Token
			forceErr := auth.UserToken.ForceRefreshToken()
			if forceErr != nil {
				util.LogError(forceErr, true)
			}

			continue
		}

		// if an unhandled error occured end the program
		return
	}
}
