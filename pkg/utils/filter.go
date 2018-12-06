package utils

import (
	"regexp"
	"strings"
)

var r = regexp.MustCompile(`\W+`)

func DetectBadWords(txt string, badWords []string) bool {
	for _, w := range r.Split(txt, -1) {
		w = strings.ToLower(w)
		for _, bad := range badWords {
			if bad == w {
				return true
			}
		}
	}
	return false
}
