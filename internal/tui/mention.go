package tui

import (
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

// maxSuggestions caps how many names the @mention popup offers at once.
const maxSuggestions = 5

// MentionPrefix returns the partial username of an unfinished @mention at the
// end of text, and whether one is in progress. The "@" must start a word, so
// "matt@rdbrck" reads as an address rather than a mention — the same rule
// mentionRegex applies to received messages.
func MentionPrefix(text string) (string, bool) {
	at := strings.LastIndex(text, "@")
	if at < 0 {
		return "", false
	}
	if prev, _ := utf8.DecodeLastRuneInString(text[:at]); isWordRune(prev) {
		return "", false
	}
	prefix := text[at+1:]
	if strings.IndexFunc(prefix, unicode.IsSpace) >= 0 {
		return "", false
	}
	return prefix, true
}

func isWordRune(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

// MatchUsernames returns the candidates starting with prefix, compared
// case-insensitively, sorted alphabetically and capped at maxSuggestions. An
// empty prefix matches everything, so typing a bare "@" lists everyone online.
func MatchUsernames(prefix string, candidates []string) []string {
	lower := strings.ToLower(prefix)
	var matches []string
	for _, c := range candidates {
		if strings.HasPrefix(strings.ToLower(c), lower) {
			matches = append(matches, c)
		}
	}
	sort.Slice(matches, func(i, j int) bool {
		return strings.ToLower(matches[i]) < strings.ToLower(matches[j])
	})
	if len(matches) > maxSuggestions {
		matches = matches[:maxSuggestions]
	}
	return matches
}
