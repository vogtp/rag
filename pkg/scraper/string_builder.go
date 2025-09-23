package scraper

import (
	"strings"
)

type stringBuilder struct {
	strings.Builder
}

func (b *stringBuilder) WriteString(s string)  {
	//fmt.Print(s)
	b.Builder.WriteString(s)
}
