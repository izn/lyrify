package spotify

import (
	"errors"
	"runtime"
)

var clients = make(map[string]SpotifyClient)

func registerClient(os string, client SpotifyClient) {
	clients[os] = client
}

func NewSpotifyClient() (SpotifyClient, error) {
	client, ok := clients[runtime.GOOS]
	if !ok {
		return nil, errors.New("unsupported operating system")
	}
	return client, nil
}
