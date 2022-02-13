package awsv4upgrade

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/minamijoyo/hcledit/editor"
	"github.com/zclconf/go-cty/cty"
)

// AWSS3BucketACLFilter is a filter implementation for upgrading the acl
// argument of aws_s3_bucket.
// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#acl-argument
type AWSS3BucketACLFilter struct{}

var _ editor.Filter = (*AWSS3BucketACLFilter)(nil)

// Filter upgrades the acl argument of aws_s3_bucket.
func (f *AWSS3BucketACLFilter) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	blocks := findResourceByType(inFile.Body(), "aws_s3_bucket")
	for _, block := range blocks {
		if block.Body().GetAttribute("acl") == nil {
			continue
		}

		resourceName := getResourceName(block)
		newblock := appendNewResourceBlock(inFile.Body(), "aws_s3_bucket_acl", resourceName)
		setBucketArgument(newblock, resourceName)
		newblock.Body().SetAttributeValue("acl", cty.StringVal("private"))
		block.Body().RemoveAttribute("acl")
	}

	return inFile, nil
}
