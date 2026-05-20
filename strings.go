package main

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Capitalize uppercases the first rune of s and leaves the rest unchanged.
// Operates on runes rather than bytes so multi-byte UTF-8 input is preserved.
func capitalize(s string) string {
	if s == "" {
		return s
	}
	r, size := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError {
		return s
	}
	return string(unicode.ToTitle(r)) + s[size:]
}

// Return what's between two strings, "a" and "b", in another string
func between(orig, a, b string) string {
	if strings.Contains(orig, a) && strings.Contains(orig, b) {
		posa := strings.Index(orig, a) + len(a)
		posb := strings.LastIndex(orig, b)
		if posa > posb {
			return ""
		}
		return orig[posa:posb]
	}
	return ""
}

// Return what's between two strings, "a" and "b", in another string,
// but inclusively, so that "a" and "b" are also included in the return value.
func betweenInclusive(orig, a, b string) string {
	if strings.Contains(orig, a) && strings.Contains(orig, b) {
		posa := strings.Index(orig, a) + len(a)
		posb := strings.LastIndex(orig, b)
		if posa > posb {
			return ""
		}
		return a + orig[posa:posb] + b
	}
	return ""
}

// Return the contents between double or single quotes (or an empty string)
func betweenQuotes(orig string) string {
	var s string
	for _, quote := range []string{"\"", "'"} {
		s = between(orig, quote, quote)
		if s != "" {
			return s
		}
	}
	return ""
}

// Return the string between the quotes or after the "=", if possible
// or return the original string
func betweenQuotesOrAfterEquals(orig string) string {
	s := betweenQuotes(orig)
	// Check for exactly one "="
	if s == "" && strings.Count(orig, "=") == 1 {
		s = strings.TrimSpace(strings.Split(orig, "=")[1])
	}
	return s
}

// Does a keyword exist in the string?
// Disregards several common special characters (like -, _ and .)
func has(s, kw string) bool {
	// Convert to lowercase and remove the most common special characters
	words := strings.Trim(strings.ToLower(s), "-_.,!?()[]{}\\/:;+@")
	wordSlice := strings.Split(words, " ")
	for _, word := range wordSlice {
		if word == kw {
			return true
		}
	}
	return false
}
