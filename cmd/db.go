package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vogtp/rag/pkg/logger"
	"github.com/vogtp/rag/pkg/usercfg"
)

func addDB() {
	rootCmd.AddCommand(dbCmd)
	dbCmd.AddCommand(dbUserCmd)
	dbCmd.AddCommand(dbCleanupCmd)
}

var dbCmd = &cobra.Command{
	Use:     "db",
	Short:   "manage the DB",
	Aliases: []string{},
	Long:    ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	},
}

var dbUserCmd = &cobra.Command{
	Use:          "user",
	Short:        "list users",
	Aliases:      []string{"users", "ls", "list"},
	Long:         ``,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		db, err := usercfg.Create(ctx, logger.New(), usercfg.DBFileName)
		if err != nil {
			return err
		}
		colCnt := 0
		fmt.Println("User list:")
		users, err := db.Users(ctx)
		if err != nil {
			return err
		}
		for _, u := range users {
			fmt.Printf(" * %s\n", u.Name)
			if len(u.Collections) > 0 {
				for _, c := range u.Collections {
					fmt.Printf("   %s (%s)\n", c.Displayname, c.Source.Parts)
					colCnt++
				}
			}
		}
		fmt.Printf("Count:\n Users: %v\n Collections: %v\n", len(users), colCnt)

		return nil
	},
}

var dbCleanupCmd = &cobra.Command{
	Use:   "cleanup [username]",
	Short: "cleanup old collections from DB",

	Aliases:      []string{},
	Long:         ``,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		db, err := usercfg.Create(ctx, logger.New(), usercfg.DBFileName)
		if err != nil {
			return err
		}
		if len(args) > 0 {
			usr, err := db.User(ctx, args[0])
			if err != nil {
				return err
			}
			return db.CleanupUserCollections(ctx, usr)
		}
		usrs, err := db.Users(ctx)
		if err != nil {
			return err
		}
		for _, usr := range usrs {
			if err := db.CleanupUserCollections(ctx, &usr); err != nil {
				return err
			}
		}
		return nil
	},
}
