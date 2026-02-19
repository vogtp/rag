package pdf

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/vogtp/langchaingo/documentloaders"
	"github.com/vogtp/langchaingo/textsplitter"
	"github.com/vogtp/rag/pkg/cfg"
	"github.com/vogtp/rag/pkg/types"
	vecdb "github.com/vogtp/rag/pkg/vecDB"
)

type HttpStatusError struct {
	StatusCode int
	Status     string
}

func (hse HttpStatusError) Error() string {
	return fmt.Sprintf("status not OK: %v (%s)", hse.StatusCode, hse.Status)
}

func SplitFromLink(ctx context.Context, link string) ([]types.EmbeddDocument, error) {
	// slog := slog.With("pdf.url", link)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, link, nil)
	if err != nil {
		// slog.Warn("Cannot create http request to get PDF", "err", err)
		return nil, fmt.Errorf("create http get request: %w", err)
	}
	req.Header.Set("User-Agent", cfg.UserAgent())
	resp, err := http.DefaultClient.Do(req)
	// resp, err := http.Get(link)
	if err != nil {
		// slog.Warn("Cannot http get PDF", "err", err)
		return nil, fmt.Errorf("do http request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		// slog.Info("http get returned an error for PDF", "code", resp.StatusCode, "status", resp.Status)
		return nil, HttpStatusError{Status: resp.Status, StatusCode: resp.StatusCode}
	}
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		// slog.Warn("Cannot read http PDF response body", "err", err)
		return nil, fmt.Errorf("read PDF response body: %w", err)
	}
	ra := bytes.NewReader(buf)
	return SplitFromReaderAt(ctx, ra, ra.Size(), vecdb.MetaURL, link)
}

func SplitFromFile(ctx context.Context, path string) ([]types.EmbeddDocument, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open pdf %q: %v", path, err)
	}
	fi, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat pdf %q: %v", path, err)
	}

	docs, err := SplitFromReaderAt(context.Background(), f, fi.Size(), vecdb.MetaPath, path)
	if err != nil {
		return nil, fmt.Errorf("spliting PDF %q: %w", path, err)
	}
	for i, d := range docs {
		if !fi.ModTime().IsZero() {
			d.MetaData[vecdb.MetaUpdated] = fi.ModTime().String()
		}
		d.MetaData[vecdb.MetaURL] = fmt.Sprintf("file://%s", path)
		docs[i] = d
	}
	return docs, nil
}

func SplitFromReaderAt(ctx context.Context, ra io.ReaderAt, size int64, idKey string, idValue string) ([]types.EmbeddDocument, error) {

	loader := documentloaders.NewPDF(ra, size)
	split := textsplitter.NewRecursiveCharacter()
	split.Separators = []string{"\n\n", "\n", ". ", " § "}
	// split.ChunkSize = 400
	// split.ChunkOverlap = 50
	docs, err := loader.LoadAndSplit(ctx, split)
	if err != nil {
		// slog.Warn("Cannot split PDF", "err", err, "size", size)
		return nil, fmt.Errorf("split PDF: %w", err)
	}
	ret := make([]types.EmbeddDocument, 0, len(docs))
	for _, d := range docs {
		if len(d.PageContent) < 10 {
			continue
		}
		doc := types.EmbeddDocument{
			Title:       filepath.Base(idValue),
			IDMetaKey:   idKey,
			IDMetaValue: idValue,
			MetaData:    make(map[string]any),
			Modified:    time.Now(),
			Document:    d.PageContent,
		}
		ret = append(ret, doc)
	}
	return ret, nil
}
