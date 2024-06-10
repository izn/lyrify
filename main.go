package main

import (
	"fmt"
	"github.com/lyrify/v2/spotify"
)

func main() {
	spotifyClient, err := spotify.NewSpotifyClient()
	if err != nil {
		fmt.Println(err)
		return
	}

	track, err := spotifyClient.CurrentTrack()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(track.Artist, "-", track.Title)

	lyrics, err := FetchLyrics(track)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(lyrics)
}
