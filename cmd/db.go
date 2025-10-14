package cmd

import (
	"fmt"

	"entgo.io/ent/dialect"
	"github.com/spf13/cobra"
	"github.com/vogtp/rag/pkg/logger"
	"github.com/vogtp/rag/pkg/usercfg"
	"github.com/vogtp/rag/pkg/usercfg/db"
	"github.com/vogtp/rag/pkg/usercfg/db/ent/user"
)

func addDB() {
	rootCmd.AddCommand(dbCmd)
	dbCmd.AddCommand(dbGormCmd)
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
	Aliases:      []string{},
	Long:         ``,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		db, err := usercfg.New(cmd.Context(), logger.New(), dialect.SQLite, "rag.sqlite")
		if err != nil {
			return err
		}
		users, err := db.GetUserQuery(ctx).All(ctx)
		if err != nil {
			return err
		}
		colCnt := 0
		fmt.Println("User list:")
		for _, u := range users {
			//u.Edges.Confluence
			cols := u.Edges.Collections
			fmt.Printf(" %s\n", u.Name)
			if len(cols) > 0 {
				for _, c := range cols {
					fmt.Printf("   %s (%s)\n", c.Name, c.Edges.Sources[0].Parts)
					colCnt++
				}
				fmt.Println()
			}
			// b, err := json.MarshalIndent(u, "", "  ")
			// if err != nil {
			// 	fmt.Println(err)
			// }
			// fmt.Print(string(b))
		}
		fmt.Printf("Count:\n Users: %v\n Collections: %v\n", len(users), colCnt)
		return nil
	},
}

var dbGormCmd = &cobra.Command{
	Use:          "gorm",
	Short:        "list users",
	Aliases:      []string{},
	Long:         ``,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		dbEnt, err := usercfg.New(cmd.Context(), logger.New(), usercfg.Dialect, usercfg.DBFileName)
		if err != nil {
			return err
		}
		usrs, err := dbEnt.GetUserQuery(ctx).All(ctx)
		if err != nil {
			return err
		}
		for _, u := range usrs {
			fmt.Println(u.Name)
		}
		return db.InitGorm()
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
		db, err := usercfg.New(cmd.Context(), logger.New(), usercfg.Dialect, usercfg.DBFileName)
		if err != nil {
			return err
		}
		if len(args) > 0 {
			usr, err := db.GetUserQuery(ctx).Where(user.Name(args[0])).First(ctx)
			if err != nil {
				return err
			}
			return db.CleanupUserCollections(ctx, usr)
		}
		usrs, err := db.GetUserQuery(ctx).All(ctx)
		if err != nil {
			return err
		}
		for _, usr := range usrs {
			if err := db.CleanupUserCollections(ctx, usr); err != nil {
				return err
			}
		}
		return nil
	},
}
