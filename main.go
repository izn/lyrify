package main

import (
	"fmt"
)

func main() {
	track, err := CurrentSpotifyTrack()
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
