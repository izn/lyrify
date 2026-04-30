package lyrics

import (
	"errors"
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
	for _, title := range titleCandidates(track.Title) {
		url := buildGeniusURL(track.Artist, title)
		lyrics, err := fetchLyricsFromURL(url)
		if err == nil && strings.TrimSpace(lyrics) != "" {
			return lyrics, nil
		}
	}

	return "", errors.New("lyrics not found on Genius for current track")
}

func fetchLyricsFromURL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("genius returned status %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return extractLyrics(string(body))
}

func buildGeniusURL(artistRaw, titleRaw string) string {
	artist := cleanString(artistRaw)
	title := cleanString(titleRaw)

	artistSlug := strings.ReplaceAll(artist, " ", "-")
	titleSlug := strings.ReplaceAll(title, " ", "-")

	artistEncoded := url.PathEscape(artistSlug)
	titleEncoded := url.PathEscape(titleSlug)

	geniusURL := fmt.Sprintf("https://genius.com/%s-%s-lyrics", artistEncoded, titleEncoded)

	return geniusURL
}

func titleCandidates(title string) []string {
	title = strings.TrimSpace(title)
	if title == "" {
		return []string{title}
	}

	candidates := []string{title}

	suffixes := []string{
		`\s*-\s*\d{4}\s+Remaster(?:ed)?$`,
		`\s*-\s*Remaster(?:ed)?(?:\s+\d{4})?$`,
		`\s*-\s*\d{4}\s+Version$`,
		`\s*-\s*Live.*$`,
		`\s*-\s*Radio\s+Edit$`,
		`\s*-\s*Mono(?:\s+Version)?$`,
		`\s*-\s*Stereo(?:\s+Version)?$`,
	}

	clean := title
	for _, pattern := range suffixes {
		re := regexp.MustCompile(`(?i)` + pattern)
		clean = strings.TrimSpace(re.ReplaceAllString(clean, ""))
	}

	if clean != "" && clean != title {
		candidates = append(candidates, clean)
	}

	return candidates
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

	if len(lyrics) == 0 {
		return "", errors.New("lyrics container not found")
	}

	return strings.Join(lyrics, ""), nil
}

func cleanString(s string) string {
	reg := regexp.MustCompile(`[^\w\s-]`)
	s = reg.ReplaceAllString(s, "")

	s = strings.TrimSpace(s)

	return s
}
