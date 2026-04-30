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
	"golang.org/x/net/html"
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
		if len(s.Nodes) == 0 {
			return
		}

		text := strings.TrimSpace(extractNodeText(s.Nodes[0]))
		if text != "" {
			lyrics = append(lyrics, text)
		}
	})

	if len(lyrics) == 0 {
		return "", errors.New("lyrics container not found")
	}

	return normalizeLyrics(strings.Join(lyrics, "")), nil
}

func normalizeLyrics(raw string) string {
	raw = strings.TrimSpace(raw)

	if strings.Contains(raw, "Read More") {
		if idx := strings.Index(raw, "["); idx >= 0 {
			raw = raw[idx:]
		}
	}

	return strings.TrimSpace(raw)
}

func extractNodeText(node *html.Node) string {
	if node == nil {
		return ""
	}

	var b strings.Builder
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n == nil {
			return
		}

		if n.Type == html.ElementNode && n.Data == "br" {
			b.WriteString("\n")
			return
		}

		if n.Type == html.TextNode {
			b.WriteString(n.Data)
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}

	walk(node)
	return b.String()
}

func cleanString(s string) string {
	reg := regexp.MustCompile(`[^\w\s-]`)
	s = reg.ReplaceAllString(s, "")

	s = strings.TrimSpace(s)

	return s
}
