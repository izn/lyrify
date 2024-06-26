package main

import (
	"fmt"
	"github.com/izn/lyrify/spotify"
	"html"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func FetchLyrics(track spotify.Track) (string, error) {
	url := buildGeniusURL(track)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	lyrics, err := extractLyrics(string(body))
	if err != nil {
		return "", err
	}

	return lyrics, nil
}

func buildGeniusURL(track spotify.Track) string {
	artist := CleanString(track.Artist)
	title := CleanString(track.Title)

	artistSlug := strings.ReplaceAll(artist, " ", "-")
	titleSlug := strings.ReplaceAll(title, " ", "-")

	artistEncoded := url.PathEscape(artistSlug)
	titleEncoded := url.PathEscape(titleSlug)

	geniusURL := fmt.Sprintf("https://genius.com/%s-%s-lyrics", artistEncoded, titleEncoded)

	return geniusURL
}

func extractLyrics(rawHtml string) (string, error) {
	re := regexp.MustCompile(`<div data-lyrics-container="true"[^>]*>(.*?)<\/div>`)
	matches := re.FindStringSubmatch(rawHtml)

	if len(matches) < 2 {
		return "", fmt.Errorf("Lyrics not found")
	}

	lyrics := matches[1]
	lyrics = html.UnescapeString(lyrics)
	lyrics = strings.ReplaceAll(lyrics, "<br/>", "\n")
	lyrics = RemoveHTMLTags(lyrics)

	return lyrics, nil
}
