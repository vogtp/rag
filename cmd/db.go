package cmd

import (
	"encoding/json"
	"fmt"

	"entgo.io/ent/dialect"
	"github.com/spf13/cobra"
	"github.com/vogtp/rag/pkg/logger"
	"github.com/vogtp/rag/pkg/usercfg"
	"github.com/vogtp/rag/pkg/usercfg/db/ent"
	"github.com/vogtp/rag/pkg/usercfg/db/ent/sourcesystem"
	"github.com/vogtp/rag/pkg/usercfg/db/ent/user"
)

func addDB() {
	rootCmd.AddCommand(dbCmd)
	dbCmd.AddCommand(dbAddCmd)
	dbCmd.AddCommand(dbTestCmd)
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
		users, err := db.All(ctx)
		if err != nil {
			return err
		}
		colCnt := 0
		fmt.Println("User list:")
		for _, u := range users {
			//u.Edges.Confluence
			cols, err := u.Collections(ctx)
			if err != nil {
				fmt.Println(err)
			}
			u.QueryCollections()
			fmt.Printf(" %s\n", u.Name)
			if len(cols) > 0 {
				fmt.Print("   ")
				for _, c := range cols {
					fmt.Printf("%s ", c.Name)
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

		srcConflPub, err := db.SourceSystem.Create().SetName("Intranet").SetURL("https://intranet-dev.unibas.ch/").SetParts("HR,KB").SetKey("xxxx").SetType(sourcesystem.TypeConfluence).Save(ctx)
		if err != nil {
			return err
		}
		srcHttplPub, err := db.SourceSystem.Create().SetName("Website").SetURL("https://its.unibas.ch/").SetParts("HR,KB").SetType(sourcesystem.TypeHTTP).Save(ctx)
		if err != nil {
			return err
		}
		srcConflPriv, err := db.SourceSystem.Create().SetName("Intranet").SetURL("https://intranet-dev.unibas.ch/").SetParts("mySpace,KB").SetKey("xxxx").SetType(sourcesystem.TypeConfluence).Save(ctx)
		if err != nil {
			return err
		}
		colPub, err := db.Collection.Create().SetName("Public collection").AddSources(srcConflPub, srcHttplPub).SetAPIKey("myPublicOpenaiAPIkey").Save(ctx)
		if err != nil {
			return err
		}
		colPriv, err := db.Collection.Create().AddSources(srcConflPriv).SetName("Private collection").SetAPIKey("mySecretOpenaiAPIkey").Save(ctx)
		if err != nil {
			return err
		}
		u, err := db.User.Create().SetName(args[0]).AddCollections(colPub, colPriv).OnConflict().UpdateNewValues().ID(ctx)
		if err != nil {
			return err
		}

		fmt.Printf("Created user: %v\n", u)

		return nil
	},
}
var tmp = `{
  "id": 8589934593,
  "Name": "vogtpa",
  "OpenaiAPIkey": "MyCoolKey",
  "edges": {
    "Confluence": [
      {
        "id": 1,
        "Name": "Intra",
        "URL": "https://intranet-dev.unibas.ch/",
        "APIKey": "yyyy",
        "edges": {
          "Spaces": [
            {
              "id": 4294967297,
              "Name": "Our space",
              "SpaceKey": "MsS"
            }
          ]
        }
      }
    ]
  }
}`

var dbTestCmd = &cobra.Command{
	Use: "test",

	Aliases:      []string{},
	Long:         ``,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		db, err := usercfg.New(cmd.Context(), logger.New(), usercfg.Dialect, usercfg.DBFileName)
		if err != nil {
			return err
		}
		var usr ent.User
		if err := json.Unmarshal([]byte(tmp), &usr); err != nil {
			return err
		}

		u, err := db.User.Query().Where(user.ID(usr.ID)).First(ctx)
		if err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(tmp), &u); err != nil {
			return err
		}
		fmt.Printf("Marsh user: %v\n", u)
		// u, err = db.User.UpdateOne(u).SetUser(u).AddConfluence(u.Confluence()...).Save(ctx)
		// if err != nil {
		// 	return err
		// }
		fmt.Printf("Updated user: %v\n", u)

		return nil
	},
}
