package filesystem

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/textsplitter"
	vecdb "github.com/vogtp/rag/pkg/vecDB"
)

func Generate(ctx context.Context, path string) chan vecdb.EmbeddDocument {
	out := make(chan vecdb.EmbeddDocument, 3)
	go walkPath(ctx, out, path)
	return out
}

func walkPath(ctx context.Context, out chan vecdb.EmbeddDocument, path string) {
	defer close(out)

	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		slog.Debug("Processing walkdir", "path", path, "dirEntry", d, "err", err)
		if d.IsDir() {
			return nil
		}

		slog := slog.With("path", path)
		i, err := d.Info()
		if err != nil {
			slog.Warn("cannot get info", "err", err)
			return err
		}
		doc := vecdb.EmbeddDocument{
			IDMetaKey:   vecdb.MetaPath,
			IDMetaValue: path,
			Modified:    i.ModTime(),
		}

		if strings.HasSuffix(strings.ToLower(path), ".pdf") {
			f, err := os.Open(path)
			if err != nil {
				slog.Warn("Cannot read pdf file", "err", err)
				return nil
			}
			pdfLoader := documentloaders.NewPDF(f, i.Size())
			docs, err := pdfLoader.LoadAndSplit(ctx, textsplitter.NewTokenSplitter())
			if err != nil {
				slog.Warn("Cannot split docs", "err", err)
				return nil
			}
			for _, d := range docs {
				doc := vecdb.EmbeddDocument{
					IDMetaKey:   vecdb.MetaPath,
					IDMetaValue: path,
					MetaData:    make(map[string]any),
					Modified:    i.ModTime(),
					Document:    d.PageContent,
				}
				doc.MetaData[vecdb.MetaURL] = fmt.Sprintf("file://%s", path)
				out <- doc
			}
			slog.Info("Added PDF")
			return nil
		}

		doc.MetaData = make(map[string]any)
		doc.MetaData[vecdb.MetaPath] = path
		doc.MetaData[vecdb.MetaUpdated] = i.ModTime().String()

		slog.Debug("adding document to chroma", "path", path)
		txt, err := os.ReadFile(path)
		if err != nil {
			slog.Warn("cannot read document", "err", err)
			return err
		}
		doc.Document = string(txt)
		out <- doc
		return nil
	})
	if err != nil {
		slog.Error("Error walking path", "err", err, "path", path)
	}
}
