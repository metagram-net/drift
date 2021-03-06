package main

import (
	"os"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/metagram-net/drift"
)

func newCmd(cli *CLI) *cobra.Command {
	var (
		// Set the default ID out of range to distinguish explicit zero.
		id   drift.MigrationID = -1
		slug string
	)

	cmd := &cobra.Command{
		Use:   "new",
		Short: "Create a new migration file",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			dir := viper.GetString("migrations-dir")
			templateFile := viper.GetString("template-file")

			tmpl, err := migrationTemplate(templateFile)
			if err != nil {
				cli.Exitf(1, "apply migration template: %s", err)
			}

			path, err := drift.NewFile(cli, dir, id, slug, tmpl)
			if err != nil {
				cli.Exitf(1, "write migration file: %s", err)
			}

			cli.Infof("Created new migration file: %s", path)
			cli.Printf(path)
		},
	}
	flags := cmd.Flags()
	flags.Var(&id, "id", "Migration ID override (default: Unix timestamp in seconds)")
	flags.StringVar(&slug, "slug", "", "Short text used to name the migration")
	cmd.MarkFlagRequired("slug")
	flags.String("template", "", "Template file for the migration")
	viper.BindPFlag("template-file", flags.Lookup("template"))
	return cmd
}

func migrationTemplate(path string) (*template.Template, error) {
	if path == "" {
		// Drift uses a sensible default template in case of nil.
		return nil, nil
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return template.New("migration").Parse(string(b))
}
