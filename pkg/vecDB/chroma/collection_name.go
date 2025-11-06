package chroma

import (
	"slices"
	"strings"
	"unicode"
)

var (
	allowed []rune = []rune{'_', '-'}
)

func FixCollectionName(n string) string {
	var ret strings.Builder
	for _, c := range n {
		if !isValidRune(c) {
			c = '_'
		}
		ret.WriteRune(c)
	}
	r := ret.String()
	if !unicode.IsLetter(rune(r[len(r)-1])) {
		r = r[:len(r)-1]
	}
	return r
}

func isValidRune(r rune) bool {
	if unicode.IsLetter(r) || unicode.IsDigit(r) {
		return true
	}
	if slices.Contains(allowed, r) {
		return true
	}
	return false
}
