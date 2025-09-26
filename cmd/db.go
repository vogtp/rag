package cmd

import (
	"fmt"

	"entgo.io/ent/dialect"
	"github.com/spf13/cobra"
	"github.com/vogtp/rag/pkg/logger"
	"github.com/vogtp/rag/pkg/usercfg"
)

func addDB() {
	rootCmd.AddCommand(dbCmd)
	dbCmd.AddCommand(dbAddCmd)
}

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "manage the  DB",

	Aliases:      []string{},
	Long:         ``,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		db, err := usercfg.New(cmd.Context(), logger.New(), dialect.SQLite, "rag.sqlite")
		if err != nil {
			return err
		}
		users, err := db.User.Query().All(ctx)
		if err != nil {
			return err
		}
		fmt.Println("User list:")
		for _, u := range users {

			fmt.Printf(" %v\n", u)
		}
		fmt.Printf("Count: %v\n", len(users))
		return nil
	},
}

var dbAddCmd = &cobra.Command{
	Use:   "add <username>",
	Short: "add users to the DB",

	Aliases:      []string{},
	Long:         ``,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		db, err := usercfg.New(cmd.Context(), logger.New(), usercfg.Dialect, usercfg.DBFileName)
		if err != nil {
			return err
		}
		if len(args) < 1 {
			return cmd.Usage()
		}
		u, err := db.User.Create().SetName(args[0]).Save(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("Created user: %v\n", u)

		return nil
	},
}
