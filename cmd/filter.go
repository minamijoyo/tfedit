package cmd

import (
	"fmt"

	"github.com/minamijoyo/tfedit/filter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(newFilterCmd())
}

func newFilterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "filter",
		Short: "Apply a built-in filter",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(
		newFilterAwsv4upgradeCmd(),
	)

	return cmd
}

func newFilterAwsv4upgradeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "awsv4upgrade",
		Short: "Apply a built-in filter for awsv4upgrade",
		Long: `Apply a built-in filter for awsv4upgrade

Upgrade configurations to AWS provider v4.
Only aws_s3_bucket refactor is supported.
`,
		RunE: runFilterAwsv4upgradeCmd,
	}

	return cmd
}

func runFilterAwsv4upgradeCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("expected 0 argument, but got %d arguments", len(args))
	}

	file := viper.GetString("file")
	update := viper.GetBool("update")
	filter, err := filter.NewFilterByType("awsv4upgrade")
	if err != nil {
		return err
	}

	c := newDefaultClient(cmd)
	return c.Edit(file, update, filter)
}
