package tui_test

import (
	"reflect"
	"testing"

	"github.com/esteban/rbchat/internal/tui"
)

func TestMentionPrefix(t *testing.T) {
	cases := []struct {
		name   string
		text   string
		prefix string
		ok     bool
	}{
		{"bare at", "@", "", true},
		{"partial name", "hey @ma", "ma", true},
		{"at start of input", "@ma", "ma", true},
		{"after punctuation", "(@ma", "ma", true},
		{"no at sign", "hello there", "", false},
		{"email address", "mail me at matt@rdbrck", "", false},
		{"completed mention", "@matt hello", "", false},
		{"at followed by space", "hey @ ", "", false},
		{"earlier mention already sent", "@matt hi @es", "es", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			prefix, ok := tui.MentionPrefix(tc.text)
			if ok != tc.ok || prefix != tc.prefix {
				t.Fatalf("MentionPrefix(%q) = (%q, %v), want (%q, %v)",
					tc.text, prefix, ok, tc.prefix, tc.ok)
			}
		})
	}
}

func TestMatchUsernames(t *testing.T) {
	candidates := []string{"esteban", "matt", "Marcus", "sofia"}

	cases := []struct {
		name   string
		prefix string
		want   []string
	}{
		{"empty prefix matches all", "", []string{"esteban", "Marcus", "matt", "sofia"}},
		{"case insensitive", "m", []string{"Marcus", "matt"}},
		{"uppercase query", "MA", []string{"Marcus", "matt"}},
		{"narrows to one", "es", []string{"esteban"}},
		{"no match", "zz", nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tui.MatchUsernames(tc.prefix, candidates)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("MatchUsernames(%q) = %v, want %v", tc.prefix, got, tc.want)
			}
		})
	}
}

func TestMatchUsernamesCapsResults(t *testing.T) {
	candidates := []string{"a1", "a2", "a3", "a4", "a5", "a6", "a7"}
	if got := tui.MatchUsernames("a", candidates); len(got) != 5 {
		t.Fatalf("expected results capped at 5, got %d", len(got))
	}
}
