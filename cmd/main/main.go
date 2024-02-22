package main

// <---------------------------------------------------------------------------------------------------->

import (
	"github.com/zmb3/spotify/v2"

	"github.com/SkillpTm/SpotifyTrueRandomShuffle/pkg/setup"
)

// <---------------------------------------------------------------------------------------------------->

var Client *spotify.Client = nil
var User *spotify.PrivateUser = nil

func main() {
	Client, User = setup.Setup()
}