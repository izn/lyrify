package main

import (
	"regexp"
	"strings"
)

func CleanString(s string) string {
	reg := regexp.MustCompile(`[^\w\s-]`)
	s = reg.ReplaceAllString(s, "")

	s = strings.TrimSpace(s)

	return s
}

func RemoveHTMLTags(htmlString string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(htmlString, "")
}
