package lyrics

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/izn/lyrify/spotify"
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
	artist := cleanString(track.Artist)
	title := cleanString(track.Title)

	artistSlug := strings.ReplaceAll(artist, " ", "-")
	titleSlug := strings.ReplaceAll(title, " ", "-")

	artistEncoded := url.PathEscape(artistSlug)
	titleEncoded := url.PathEscape(titleSlug)

	geniusURL := fmt.Sprintf("https://genius.com/%s-%s-lyrics", artistEncoded, titleEncoded)

	return geniusURL
}

func extractLyrics(rawHTML string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(rawHTML))
	if err != nil {
		return "", err
	}

	var lyrics []string

	doc.Find("div[data-lyrics-container='true']").Each(func(i int, s *goquery.Selection) {
		s.Contents().Each(func(j int, node *goquery.Selection) {
			if goquery.NodeName(node) == "#text" {
				text := strings.TrimSpace(node.Text())
				if text != "" {
					lyrics = append(lyrics, text)
				}
				return
			}

			if goquery.NodeName(node) == "br" {
				lyrics = append(lyrics, "\n")
				return
			}

			if goquery.NodeName(node) == "a" {
				text := strings.TrimSpace(node.Text())
				if text != "" {
					lyrics = append(lyrics, text)
				}
			}
		})
	})

	return strings.Join(lyrics, ""), nil
}

func cleanString(s string) string {
	reg := regexp.MustCompile(`[^\w\s-]`)
	s = reg.ReplaceAllString(s, "")

	s = strings.TrimSpace(s)

	return s
}
