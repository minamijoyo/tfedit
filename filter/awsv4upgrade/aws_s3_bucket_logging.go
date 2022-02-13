package awsv4upgrade

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/minamijoyo/hcledit/editor"
)

// AWSS3BucketLoggingFilter is a filter implementation for upgrading the
// logging argument of aws_s3_bucket.
// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#logging-argument
type AWSS3BucketLoggingFilter struct {
}

var _ editor.Filter = (*AWSS3BucketLoggingFilter)(nil)

// NewAWSS3BucketLoggingFilter creates a new instance of AWSS3BucketLoggingFilter.
func NewAWSS3BucketLoggingFilter() editor.Filter {
	return &AWSS3BucketLoggingFilter{}
}

// Filter upgrades the logging argument of aws_s3_bucket.
func (f *AWSS3BucketLoggingFilter) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	blocks := inFile.Body().Blocks()
	for _, b := range blocks {
		if b.Type() != "resource" {
			continue
		}

		labels := b.Labels()
		if len(labels) == 2 && labels[0] != "aws_s3_bucket" {
			continue
		}
		bName := labels[1]

		if nested := b.Body().FirstMatchingBlock("logging", []string{}); nested != nil {
			inFile.Body().AppendNewline()
			newblock := inFile.Body().AppendNewBlock("resource", []string{"aws_s3_bucket_logging", bName})
			newblock.Body().SetAttributeTraversal("bucket", hcl.Traversal{
				hcl.TraverseRoot{Name: "aws_s3_bucket"},
				hcl.TraverseAttr{Name: bName},
				hcl.TraverseAttr{Name: "id"},
			})
			newblock.Body().AppendUnstructuredTokens(nested.Body().BuildTokens(nil))
			b.Body().RemoveBlock(nested)
		}
	}

	return inFile, nil
}
