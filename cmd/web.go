package cmd

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/vogtp/rag/pkg/usercfg"
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
	Aliases: []string{"w", "rag", "r", "http"},
}

var webStartCmd = &cobra.Command{
	Use:          "start",
	Short:        "Start RAG web server",
	Aliases:      []string{"s", "serve"},
	SilenceUsage: true,
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
	userCfg, err := usercfg.Create(ctx, slog, usercfg.DBFileName)
	if err != nil {
		return err
	}
	api, err := web.New(ctx, slog, userCfg)
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
