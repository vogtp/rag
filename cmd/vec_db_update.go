package cmd

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/spf13/cobra"
	"github.com/vogtp/rag/pkg/usercfg"
)

var vecDbUpdate = &cobra.Command{
	Use:   "update",
	Short: "Update collections",

	Aliases: []string{"u", "up"},
	Long:    ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		slog := slog.Default()
		_, err := startChroma(ctx, slog)
		if err != nil {
			return fmt.Errorf("start chroma: %w", err)
		}
		userCfg, err := usercfg.Create(ctx, slog, usercfg.DBFileName)
		if err != nil {
			return err
		}
		cols, err := userCfg.CollectionsToUpdate(ctx, time.Now())
		if err != nil {
			return err
		}
		// slog.Info("Found collections to update, doing one at a time", "count", len(cols))
		for _, c := range cols {
			slog.Info("Embedding collection", "collection", c, "next update", c.NextDBUpdate.Format(time.DateTime))
			// returning here is correct since we only want to update one collection at a time
			return c.Embbed(ctx, slog)
		}
		return nil
	},
}
