package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vogtp/rag/pkg/cfg"
	"github.com/vogtp/rag/pkg/logger"
)

// New creates a cobra
func New() *cobra.Command {

	cfg.Parse()
	logger.New()

	addVecDB()
	addOllama()
	addWeb()
	addDB()
	//experiments.Add(rootCmd)

	return rootCmd
}

var rootCmd = &cobra.Command{
	Use:   "ragctl",
	Short: "RAG commandline",
	// SilenceUsage:  true,
	// SilenceErrors: true,
	// CompletionOptions: cobra.CompletionOptions{
	// 	DisableDefaultCmd: true,
	// },
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Print(cmd.UsageString())
		cmd.Printf("Version %v\n", cfg.Version)
	},
}
