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
	RootCmd.AddCommand(newMigrationCmd())
}

func newMigrationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migration",
		Short: "Generate a migration file for state operations",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(
		newMigrationFromplanCmd(),
	)

	return cmd
}

func newMigrationFromplanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fromplan",
		Short: "Generate a migration file from Terraform JSON plan file",
		Long: `Generate a migration file from Terraform JSON plan file

Read a Terraform plan file in JSON format and
generate a migration file in tfmigrate HCL format.
Currently, only import actions required by awsv4upgrade are supported.
`,
		RunE: runMigrationFromplanCmd,
	}

	flags := cmd.Flags()
	flags.StringP("file", "f", "-", "A path to input Terraform JSON plan file")
	flags.StringP("out", "o", "-", "Write a migration file to a given path")
	flags.StringP("dir", "d", "", "Set a dir attribute in a migration file")
	_ = viper.BindPFlag("migration.fromplan.file", flags.Lookup("file"))
	_ = viper.BindPFlag("migration.fromplan.out", flags.Lookup("out"))
	_ = viper.BindPFlag("migration.fromplan.dir", flags.Lookup("dir"))

	return cmd
}

func runMigrationFromplanCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("expected 0 argument, but got %d arguments", len(args))
	}

	planFile := viper.GetString("migration.fromplan.file")
	migrationFile := viper.GetString("migration.fromplan.out")
	migrationDir := viper.GetString("migration.fromplan.dir")

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

	output, err := migration.GenerateFromPlan(planJSON, migrationDir)
	if err != nil {
		return err
	}

	// Suppress creating a migration file when no action.
	// It is not only redundant, but also causes an error as an invalid migration
	// file when loaded by tfmigrate.
	if len(output) == 0 {
		// Intentionally does not return errors so that we can simply ignore
		// irrelevant directories when processing multiple directories.
		return nil
	}

	if migrationFile == "-" {
		fmt.Fprint(cmd.OutOrStdout(), string(output))
	} else {
		// nolint: gosec
		// G306: Expect WriteFile permissions to be 0600 or less
		// In general, a migration file is expected to commit to git and it does
		// not contain any credentials, so there is no problem.
		if err := os.WriteFile(migrationFile, output, 0644); err != nil {
			return fmt.Errorf("failed to write file: %s", err)
		}
	}

	return nil
}
