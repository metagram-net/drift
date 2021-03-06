package main

import (
	_ "github.com/jackc/pgx/v4/stdlib" // database/sql driver: pgx
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/metagram-net/drift"
)

const renumberLong string = `Renumber migrations to fix filesystem sorting.

This command renames migration files so that string sorting matches the numeric
sorting of the IDs. This happens by adding or removing prefix zeroes on the IDs
to make all the IDs the shortest width that fits them all.

Other commands ignore zero prefixes when interpreting IDs as integers. This
renumbering is never necessary for correctness.`

func renumberCmd(cli *CLI) *cobra.Command {
	var write bool

	cmd := &cobra.Command{
		Use:   "renumber",
		Short: "Renumber migrations to fix filesystem sorting",
		Long:  renumberLong,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, _ []string) {
			dir := viper.GetString("migrations-dir")
			err := drift.Renumber(cli, dir, write)
			if err != nil {
				cli.Exitf(1, "renumber: %s", err)
			}
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&write, "write", "w", false, "Execute renames instead of just printing them")
	return cmd
}
