package main

// <---------------------------------------------------------------------------------------------------->

import (
	"github.com/SkillpTm/SpotifyTrueRandomShuffle/internal/api"
	"github.com/SkillpTm/SpotifyTrueRandomShuffle/internal/player"
)

// <---------------------------------------------------------------------------------------------------->

func main() {
    api.AuthUser()
    player.Start()
}