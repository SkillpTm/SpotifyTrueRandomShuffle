package main

// <---------------------------------------------------------------------------------------------------->

import (
	"fmt"

	"github.com/SkillpTm/SpotifyTrueRandomShuffle/pkg/setup"
)

// <---------------------------------------------------------------------------------------------------->

func main() {

	id, secret, redirectURL := setup.GetEnvs()

	fmt.Println(id, secret, redirectURL)
}