package confluence

import (
	"bytes"
	"io"
	"log/slog"
	"net/url"
	"os"
	"regexp"
	"strings"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
)

var (
	markdownRegex = regexp.MustCompile(`\[[^][]+]\((https?://[^()]+)\)`)
)

func parsePage(slog *slog.Logger, p string) (string, []string) {

	// capture stdout since one of the libs writes to it
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	markdown, err := htmltomarkdown.ConvertString(p)

	// read output
	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r) //nolint
		outC <- buf.String()
	}()
	w.Close()
	os.Stdout = old
	out := <-outC
	if len(out) > 0 {
		slog.Warn("HTML to MD convert output", "output", out)
	}

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
		if _, err = url.Parse(lnk); err != nil {
			slog.Warn("Incorrect confluence link", "link", lnk, "err", err)
			continue
		}
		pdfLnks = append(pdfLnks, lnk)
		slog.Debug("Markdown links", "href", lnk)
	}
	return markdown, pdfLnks
}
