package main 

import (
	"regexp"
	"strings"
)

func cleanString(s string) string {
	reg := regexp.MustCompile(`[^\w\s-]`)
	s = reg.ReplaceAllString(s, "")

	s = strings.TrimSpace(s)

	return s
}

func removeHTMLTags(htmlString string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(htmlString, "")
}
