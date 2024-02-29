package main

// <---------------------------------------------------------------------------------------------------->

import (
	"fmt"
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
		log.Fatal(fmt.Errorf("couldn't setup config; %s", err.Error()))
	}
	
	api.AuthUser()
	err = player.Start()
	if err != nil {
		util.LogError(fmt.Errorf("couldn't continue main loop; %s", err.Error()))
	}
}

// TODO: update scopes