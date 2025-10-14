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
	Aliases:      []string{"users", "ls", "list"},
	Long:         ``,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		////////////////// ENT
		entDB, err := usercfg.NewENT(cmd.Context(), logger.New(), dialect.SQLite, "rag.sqlite")
		if err != nil {
			return err
		}
		usersEnt, err := entDB.GetUserQuery(ctx).All(ctx)
		if err != nil {
			return err
		}
		colCnt := 0
		fmt.Println("User list ENT:")
		for _, u := range usersEnt {
			//u.Edges.Confluence
			cols := u.Edges.Collections
			fmt.Printf(" %s\n", u.Name)
			if len(cols) > 0 {
				for _, c := range cols {
					fmt.Printf("   %s (%s)\n", c.Name, c.Edges.Sources[0].Parts)
					colCnt++
				}
			}
			// b, err := json.MarshalIndent(u, "", "  ")
			// if err != nil {
			// 	fmt.Println(err)
			// }
			// fmt.Print(string(b))
		}
		fmt.Printf("Count:\n Users: %v\n Collections: %v\n", len(usersEnt), colCnt)
		//////////////////// GROM
		db, err := usercfg.Create(ctx, logger.New())
		if err != nil {
			return err
		}
		colCnt = 0
		fmt.Println("User list:")
		users, err := db.Users(ctx)
		if err != nil {
			return err
		}
		for _, u := range users {
			fmt.Printf(" %s\n", u.Name)
			if len(u.Collections) > 0 {
				for _, c := range u.Collections {
					fmt.Printf("   %s (%s)\n", c.DisplayName, c.Source.Parts)
					colCnt++
				}
			}
		}
		fmt.Printf("Count:\n Users: %v\n Collections: %v\n", len(users), colCnt)

		return nil
	},
}

var dbGormCmd = &cobra.Command{
	Use:          "migrate",
	Short:        "migrate to gorm",
	Aliases:      []string{"gorm"},
	Long:         ``,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		return usercfg.Migrate2Gorm(ctx, logger.New())
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
		db, err := usercfg.Create(ctx, logger.New())
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
