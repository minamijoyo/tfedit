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
		Use:   "filter <FILTER_TYPE>",
		Short: "Apply a built-in filter",
		Long: `Apply a built-in filter

Arguments:
  FILTER_TYPE    A type of filter.
                 Valid values are:
                 - awsv4upgrade
                   Upgrade configurations to AWS provider v4.
                   Only aws_s3_bucket refactor is supported.
`,
		RunE: runFilterCmd,
	}

	return cmd
}

func runFilterCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 argument, but got %d arguments", len(args))
	}

	file := viper.GetString("file")
	update := viper.GetBool("update")
	filter, err := filter.NewFilterByType(args[0])
	if err != nil {
		return err
	}

	c := newDefaultClient(cmd)
	return c.Edit(file, update, filter)
}
