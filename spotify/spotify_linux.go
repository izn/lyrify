// +build linux

package spotify

import (
	"errors"
)

type LinuxSpotifyClient struct{}

func init() {
	registerClient("linux", &LinuxSpotifyClient{})
}

func (c *LinuxSpotifyClient) CurrentTrack() (Track, error) {
	return Track{}, errors.New("Not implemented for Linux yet")
}
