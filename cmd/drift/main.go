package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const defaultMigrationsDir = "migrations"

func init() {
	viper.SetConfigName("drift")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("DRIFT")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("migrations-dir", defaultMigrationsDir)
	viper.SetDefault("verbosity", 1)
	viper.SetDefault("template-file", "")
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	go func() {
		<-ctx.Done()
		stop()
		log.Print("Interrupt received, cleaning up before quitting. Interrupt again to force-quit.")
	}()

	err := rootCmd().ExecuteContext(ctx)
	if err != nil {
		os.Exit(1)
	}
}

func rootCmd() *cobra.Command {
	cli := &CLI{
		stdout:    os.Stdout,
		stderr:    os.Stderr,
		verbosity: InfoLevel,
	}

	cmd := &cobra.Command{
		Use:     "drift",
		Short:   "Manage SQL migrations",
		Version: "0.1.1",
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			err := viper.ReadInConfig()
			var notFound viper.ConfigFileNotFoundError
			if errors.As(err, &notFound) {
				// The config file is optional, so use the defaults.
			} else if err != nil {
				return err
			}

			cli.SetVerbosity(Verbosity(viper.GetInt("verbosity")))
			return nil
		},
	}

	flags := cmd.PersistentFlags()
	flags.String("migrations-dir", defaultMigrationsDir, "Directory containing migration files")
	flags.CountP("verbosity", "v", "Log verbosity")
	viper.BindPFlags(flags)

	cmd.AddCommand(
		migrateCmd(cli),
		newCmd(cli),
		setupCmd(cli),
		renumberCmd(cli),
		migrationTemplateCmd(cli),
	)
	return cmd
}
