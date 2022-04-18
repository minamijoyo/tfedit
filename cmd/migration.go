package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/minamijoyo/tfedit/migration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	migrationCmd := newMigrationCmd()
	flags := migrationCmd.PersistentFlags()
	flags.StringP("file", "f", "-", "A path to input Terraform plan file in JSON format")
	flags.StringP("out", "o", "-", "Write a migration file to a given path")
	_ = viper.BindPFlag("file", flags.Lookup("file"))
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
		Use:   "awsv4upgrade",
		Short: "Generate a migration file for awsv4upgrade",
		Long: `Generate a migration file for awsv4upgrade

Read a Terraform plan file in JSON format and
generate a migration file in tfmigrate HCL format.
Only aws_s3_bucket refactor is supported.
`,
		RunE: runMigrationAwsv4upgradeCmd,
	}

	return cmd
}

func runMigrationAwsv4upgradeCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("expected 0 argument, but got %d arguments", len(args))
	}

	planFile := viper.GetString("file")
	migrationFile := viper.GetString("out")

	var planJSON []byte
	var err error
	if planFile == "-" {
		planJSON, err = io.ReadAll(cmd.InOrStdin())
		if err != nil {
			return fmt.Errorf("failed to read input from stdin: %s", err)
		}
	} else {
		planJSON, err = os.ReadFile(planFile)
		if err != nil {
			return fmt.Errorf("failed to read file: %s", err)
		}
	}

	output, err := migration.Generate(planJSON)
	if err != nil {
		return err
	}

	if migrationFile == "-" {
		fmt.Fprintf(cmd.OutOrStdout(), string(output))
	} else {
		if err := os.WriteFile(migrationFile, output, os.ModePerm); err != nil {
			return fmt.Errorf("failed to write file: %s", err)
		}
	}

	return nil
}
