package awsv4upgrade

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/minamijoyo/hcledit/editor"
)

// AWSS3BucketFilter is a filter implementation for upgrading arguments of
// aws_s3_bucket to AWS provider v4.
// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#s3-bucket-refactor
type AWSS3BucketFilter struct {
}

var _ editor.Filter = (*AWSS3BucketFilter)(nil)

// NewAWSS3BucketFilter creates a new instance of AWSS3BucketFilter.
func NewAWSS3BucketFilter() editor.Filter {
	return &AWSS3BucketFilter{}
}

// Filter upgrades arguments of aws_s3_bucket to AWS provider v4.
func (f *AWSS3BucketFilter) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	m := editor.NewMultiFilter([]editor.Filter{
		NewAWSS3BucketACLFilter(),
		NewAWSS3BucketLoggingFilter(),
	})
	return m.Filter(inFile)
}

// findResourceByType returns all matching blocks from the body that have the
// given resourceType or returns an empty list if there is no matching block.
// This method is useful when you want to ignore the resource name.
func findResourceByType(body *hclwrite.Body, resourceType string) []*hclwrite.Block {
	var matched []*hclwrite.Block

	for _, block := range body.Blocks() {
		if block.Type() != "resource" {
			continue
		}

		labels := block.Labels()
		if len(labels) == 2 && labels[0] != resourceType {
			continue
		}

		matched = append(matched, block)
	}

	return matched
}

// getResourceName is a helper method for getting a resource name of the given block.
func getResourceName(block *hclwrite.Block) string {
	labels := block.Labels()
	return labels[1]
}

// appendNewResourceBlock is a helper method for appending a new resource block
// to the given body and returns a new block.
func appendNewResourceBlock(body *hclwrite.Body, resourceType string, resourceName string) *hclwrite.Block {
	body.AppendNewline()
	return body.AppendNewBlock("resource", []string{resourceType, resourceName})
}
