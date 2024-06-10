// +build darwin

package spotify

import (
	"errors"
	"os/exec"
	"strings"
)

type DarwinSpotifyClient struct{}

func init() {
	registerClient("darwin", &DarwinSpotifyClient{})
}

func (c *DarwinSpotifyClient) CurrentTrack() (Track, error) {
	script := `
		tell application "Spotify"
			if it is running then
				set trackName to name of current track
				set artistName to artist of current track
				return artistName & "\t" & trackName
			else
				return "Spotify is not running"
			end if
		end tell
	`

	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return Track{}, err
	}

	parts := strings.Split(string(output), "\t")

	if len(parts) < 2 {
		return Track{}, errors.New("Is Spotify running?")
	}

	return Track{
		Artist: parts[0],
		Title:  parts[1],
	}, nil
}
