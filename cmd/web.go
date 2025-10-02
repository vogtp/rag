package cmd

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/vogtp/rag/pkg/cfg"
	"github.com/vogtp/rag/pkg/rag"
	"github.com/vogtp/rag/pkg/vecDB/chroma"
	"github.com/vogtp/rag/pkg/web"
)

func addWeb() {
	rootCmd.AddCommand(webCmd)
	webCmd.AddCommand(webStartCmd)
}

var webCmd = &cobra.Command{
	Use:     "web",
	Short:   "Manage RAG web server",
	Aliases: []string{"w", "rag", "r"},
}

var webStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start RAG web server",
	//Aliases: []string{"w", "rag", "r"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return startWeb(cmd.Context())
	},
}

func startWeb(ctx context.Context) error {
	slog := slog.Default()
	_, err := startChroma(ctx, slog)
	if err != nil {
		return fmt.Errorf("start chroma: %w", err)
	}
	ragCfgs, err := cfg.GetRagConfig()
	if err != nil {
		return fmt.Errorf("read RAG config: %w", err)
	}
	rags := make([]rag.Manager, len(ragCfgs))
	for i, ragCfg := range ragCfgs {
		rag, err := rag.New(ctx, slog, ragCfg)
		if err != nil {
			return fmt.Errorf("start rag %q backend: %w", ragCfg.Name, err)
		}
		rags[i] = *rag
	}
	api, err := web.New(ctx, slog, rags)
	if err != nil {
		return fmt.Errorf("start http server: %w", err)
	}
	return api.Run(ctx)
}

func startChroma(ctx context.Context, slog *slog.Logger) (func(ctx context.Context) error, error) {
	c, err := chroma.NewContainer(slog)
	if err != nil {
		return nil, fmt.Errorf("create chroma container: %w", err)
	}
	return c.EnsureStarted(ctx)
}
