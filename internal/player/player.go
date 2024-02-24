// Package player ...
package player

// <---------------------------------------------------------------------------------------------------->



type Player struct {
	country string

	isPlaying bool
	isPrivateSession bool
	currentlyPlayingType string
	repeatState string
	shuffleState bool
	smartShuffle bool
	contextType string
	contextHREF string

	lastQueueSongURI string
}


func (player *Player) RunChecks() bool {
	if (!player.isPlaying || // check if the user is playing something right now
		player.isPrivateSession || // check if the user is in a private session
		player.currentlyPlayingType != "track" || // check if the user is listening to a track
		player.repeatState == "track" || // check if the user has repeat turned on for this track
		!player.shuffleState || // check if the player has turned on shuffle
		player.smartShuffle || // check if the user hasn't turned on smart shuffle
		player.contextType == "show" || // check if the user's context isn't a show
		player.contextType == "artist") { // check if the user's context isn't an artist
		return true
	}
	return false
}