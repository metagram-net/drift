package main

import (
	"github.com/metagram-net/drift"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func setupCmd(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "setup",
		Aliases: []string{"init"},
		Short:   "Set up the migrations directory",
		Args:    cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			path, err := drift.Setup(viper.GetString("migrations-dir"))
			if err != nil {
				cli.Exitf(1, "set up migrations: %s", err)
			}

			cli.Infof("Created the first migration file: %s", path)
			cli.Infof("Run the migrate command to apply it.")
		},
	}
	return cmd
}
