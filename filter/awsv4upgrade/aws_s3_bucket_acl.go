package awsv4upgrade

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/minamijoyo/hcledit/editor"
	"github.com/minamijoyo/tfedit/tfwrite"
)

// AWSS3BucketACLFilter is a filter implementation for upgrading the acl
// argument of aws_s3_bucket.
// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#acl-argument
type AWSS3BucketACLFilter struct{}

var _ editor.Filter = (*AWSS3BucketACLFilter)(nil)

// Filter upgrades the acl argument of aws_s3_bucket.
func (f *AWSS3BucketACLFilter) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	targets := tfwrite.FindResourcesByType(inFile.Body(), "aws_s3_bucket")
	for _, oldResource := range targets {
		attr := oldResource.Body().GetAttribute("acl")
		if attr == nil {
			continue
		}

		resourceName := tfwrite.GetResourceName(oldResource)
		newResource := tfwrite.AppendNewResource(inFile.Body(), "aws_s3_bucket_acl", resourceName)
		setBucketArgument(newResource, resourceName)
		newResource.Body().AppendUnstructuredTokens(attr.BuildTokens(nil))
		oldResource.Body().RemoveAttribute("acl")
	}

	return inFile, nil
}
