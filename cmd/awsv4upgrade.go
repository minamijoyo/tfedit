package cmd

import (
	"fmt"

	"github.com/minamijoyo/tfedit/filter/awsv4upgrade"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(newAWSV4UpgradeCmd())
}

func newAWSV4UpgradeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aws_v4_upgrade",
		Short: "Upgrade resources to AWS provider v4",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(
		newAWSV4UpgradeAWSS3BucketCmd(),
	)

	return cmd
}

func newAWSV4UpgradeAWSS3BucketCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aws_s3_bucket",
		Short: "Upgrade aws_s3_bucket to AWS provider v4",
		Long:  "Upgrade aws_s3_bucket to AWS provider v4",
		RunE:  runAWSV4UpgradeAWSS3BucketCmd,
	}

	return cmd
}

func runAWSV4UpgradeAWSS3BucketCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("expected no argument, but got %d arguments", len(args))
	}

	file := viper.GetString("file")
	update := viper.GetBool("update")

	filter := awsv4upgrade.NewAWSS3BucketFilter()
	c := newDefaultClient(cmd)
	return c.Edit(file, update, filter)
}
