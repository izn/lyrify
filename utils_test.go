package main

import (
	"testing"
)

func TestCleanString(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"Hello, World!", "Hello World"},
		{"Hello-World123!", "Hello-World123"},
		{"Special@Characters", "SpecialCharacters"},
		{"   Trimmed Spaces   ", "Trimmed Spaces"},
	}

	for _, tc := range testCases {
		result := CleanString(tc.input)
		if result != tc.expected {
			t.Errorf("expected '%s', got '%s'", tc.expected, result)
		}
	}
}

func TestRemoveHTMLTags(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"<p>Hello, <b>World</b>!</p>", "Hello, World!"},
		{"<div><p>Test</p></div>", "Test"},
		{"<h1>Title</h1><p>Content</p>", "TitleContent"},
	}

	for _, tc := range testCases {
		result := RemoveHTMLTags(tc.input)
		if result != tc.expected {
			t.Errorf("expected '%s', got '%s'", tc.expected, result)
		}
	}
}
