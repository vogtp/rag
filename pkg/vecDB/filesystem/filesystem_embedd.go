package filesystem

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/vogtp/rag/pkg/types"
	vecdb "github.com/vogtp/rag/pkg/vecDB"
	"github.com/vogtp/rag/pkg/vecDB/pdf"
)

func Generate(ctx context.Context, path string) chan *types.EmbeddDocument {
	out := make(chan *types.EmbeddDocument, 1)
	go walkPath(ctx, out, path)
	return out
}

func walkPath(ctx context.Context, out chan *types.EmbeddDocument, path string) {
	defer close(out)

	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if d == nil {
			slog.Warn("Directory not present", "path", path, "dirEntry", d)
			return fmt.Errorf("directory does not exist: %s", path)
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
		doc := types.EmbeddDocument{
			IDMetaKey:   vecdb.MetaPath,
			IDMetaValue: path,
			Modified:    i.ModTime(),
		}

		if strings.HasSuffix(strings.ToLower(path), ".pdf") {
			docs, err := pdf.SplitFromFile(ctx, path)
			if err != nil {
				slog.Warn("Cannot split PDF", "path", path, "err", err)
				return nil
			}
			for _, d := range docs {
				out <- &d
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
		out <- &doc
		return nil
	})
	if err != nil {
		slog.Error("Error walking path", "err", err, "path", path)
	}
}
