package spotify

type SpotifyClient interface {
	CurrentTrack() (Track, error)
}
