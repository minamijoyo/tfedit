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
	blocks := findResourceByType(inFile.Body(), "aws_s3_bucket")
	for _, block := range blocks {
		nested := block.Body().FirstMatchingBlock("logging", []string{})
		if nested == nil {
			continue
		}

		labels := block.Labels()
		resourceName := labels[1]

		inFile.Body().AppendNewline()
		newblock := inFile.Body().AppendNewBlock("resource", []string{"aws_s3_bucket_logging", resourceName})
		newblock.Body().SetAttributeTraversal("bucket", hcl.Traversal{
			hcl.TraverseRoot{Name: "aws_s3_bucket"},
			hcl.TraverseAttr{Name: resourceName},
			hcl.TraverseAttr{Name: "id"},
		})
		newblock.Body().AppendUnstructuredTokens(nested.Body().BuildTokens(nil))
		block.Body().RemoveBlock(nested)
	}

	return inFile, nil
}
