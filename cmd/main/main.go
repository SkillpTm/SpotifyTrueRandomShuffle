package main

// <---------------------------------------------------------------------------------------------------->

import (
	"log"

	"github.com/SkillpTm/SpotifyTrueRandomShuffle/internal/api"
	"github.com/SkillpTm/SpotifyTrueRandomShuffle/internal/player"
	"github.com/SkillpTm/SpotifyTrueRandomShuffle/internal/util"
)

// <---------------------------------------------------------------------------------------------------->

// main is the entry into the program, it sets up our config, envs, authorizes with the user and starts the main loop
func main() {
	err := util.Setup()
	if err != nil {
		log.Fatal(err)
	}
	
	api.AuthUser()
	player.Start()
}

// TODO: update scopes