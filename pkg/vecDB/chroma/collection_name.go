package chroma

import "strings"

func FixCollectionName(n string) string {
	return strings.ReplaceAll(n, " ", "_")
}
