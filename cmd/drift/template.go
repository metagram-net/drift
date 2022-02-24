package main

import (
	_ "github.com/jackc/pgx/v4/stdlib" // database/sql driver: pgx
	"github.com/spf13/cobra"

	"github.com/metagram-net/drift"
)

func migrationTemplateCmd(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migration-template",
		Short: "Print the embedded default migration template",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, _ []string) {
			cli.Printf(drift.DefaultTemplate())
		},
	}
	return cmd
}
