package cmd

import (
	"fmt"
	"os"

	"github.com/minamijoyo/tfedit/migration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	migrationCmd := newMigrationCmd()
	flags := migrationCmd.PersistentFlags()
	flags.StringP("out", "o", "-", "Write a migration file to a given path")
	_ = viper.BindPFlag("out", flags.Lookup("out"))

	RootCmd.AddCommand(migrationCmd)
}

func newMigrationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migration",
		Short: "Generate a migration file for a built-in filter",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(
		newMigrationAwsv4upgradeCmd(),
	)

	return cmd
}

func newMigrationAwsv4upgradeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "awsv4upgrade <PATH>",
		Short: "Generate a migration file for awsv4upgrade",
		Long: `Generate a migration file for awsv4upgrade

Generate a migration file as a tfmigrate format.
Only aws_s3_bucket refactor is supported.


Arguments:
  PATH               A path of Terraform plan file.
                     Only JSON format is supported.
`,
		RunE: runMigrationAwsv4upgradeCmd,
	}

	return cmd
}

func runMigrationAwsv4upgradeCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 argument, but got %d arguments", len(args))
	}

	planFile := args[0]
	migrationFile := viper.GetString("out")
	output, err := migration.Generate(planFile)
	if err != nil {
		return err
	}

	if migrationFile == "-" {
		// Write to stdout
		fmt.Fprintf(cmd.OutOrStdout(), string(output))
	} else {
		if err := os.WriteFile(migrationFile, output, os.ModePerm); err != nil {
			return fmt.Errorf("failed to write file: %s", err)
		}
	}

	return nil
}
