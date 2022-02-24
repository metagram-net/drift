package main

import (
	"database/sql"

	_ "github.com/jackc/pgx/v4/stdlib" // database/sql driver: pgx
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/metagram-net/drift"
)

func migrateCmd(cli *CLI) *cobra.Command {
	// Set the default ID out of range to distinguish explicit zero.
	uptoID := drift.MigrationID(-1)

	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run migrations",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, _ []string) {
			ctx := cmd.Context()
			dir := viper.GetString("migrations-dir")

			db, err := sql.Open("pgx", viper.GetString("database-url"))
			if err != nil {
				cli.Exitf(1, "open database connection: %s", err)
			}
			defer db.Close()

			var upto *drift.MigrationID
			if uptoID >= 0 {
				upto = &uptoID
			}

			err = drift.Migrate(ctx, cli, db, dir, upto)
			if err != nil {
				cli.Exitf(1, "run migrations: %s", err)
			}
		},
	}

	flags := cmd.Flags()
	flags.Var(&uptoID, "upto", "Maximum migration ID to run (default: run all migrations)")
	return cmd
}
