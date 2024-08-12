package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "0.2.1"
)

func init() {
	RootCmd.AddCommand(newVersionCmd())
}

func newVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version",
		RunE:  runVersionCmd,
	}

	return cmd
}

func runVersionCmd(cmd *cobra.Command, _ []string) error {
	_, err := fmt.Fprintln(cmd.OutOrStdout(), version)
	return err
}
