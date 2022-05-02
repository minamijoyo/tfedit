package cmd

import (
	"os"

	"github.com/minamijoyo/hcledit/editor"
	"github.com/spf13/cobra"
)

// RootCmd is a top level command instance
var RootCmd = &cobra.Command{
	Use:           "tfedit",
	Short:         "A refactoring tool for Terraform",
	SilenceErrors: true,
	SilenceUsage:  true,
}

func init() {
	setDefaultStream(RootCmd)
}

func setDefaultStream(cmd *cobra.Command) {
	cmd.SetIn(os.Stdin)
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
}

func newDefaultClient(cmd *cobra.Command) editor.Client {
	o := &editor.Option{
		InStream:  cmd.InOrStdin(),
		OutStream: cmd.OutOrStdout(),
		ErrStream: cmd.ErrOrStderr(),
	}
	return editor.NewClient(o)
}
