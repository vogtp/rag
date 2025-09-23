package confluence

import (
	"log/slog"
	"net/url"
	"regexp"
	"strings"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
)

var (
	markdownRegex = regexp.MustCompile(`\[[^][]+]\((https?://[^()]+)\)`)
)

func parsePage(slog *slog.Logger, p string) (string, []string) {
	// retrun html2text.HTML2Text(p)
	markdown, err := htmltomarkdown.ConvertString(p)
	if err != nil {
		slog.Error("cannot encode html to markdown", "err", err)
		return p, []string{}
	}
	results := markdownRegex.FindAllStringSubmatch(markdown, -1)
	pdfLnks := make([]string, 0, len(results))
	for v := range results {
		lnk := results[v][1]
		if !strings.HasSuffix(strings.ToLower(lnk), ".pdf") {
			continue
		}
		if _, err =url.Parse(lnk); err != nil{
			slog.Warn("Incorrect confluence link", "link",lnk, "err",err)
			continue
		}
		pdfLnks = append(pdfLnks, lnk)
		slog.Debug("Markdown links", "href", lnk)
	}
	return markdown, pdfLnks
}
