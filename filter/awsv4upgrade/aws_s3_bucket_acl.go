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
	file := tfwrite.NewFile(inFile)
	oldResourceType := "aws_s3_bucket"
	oldAttribute := "acl"
	oldResourceRefAttribute := "id"
	newResourceType := "aws_s3_bucket_acl"
	newResourceRefAttribute := "bucket"

	targets := file.FindResourcesByType(oldResourceType)
	for _, oldResource := range targets {
		attr := oldResource.GetAttribute(oldAttribute)
		if attr == nil {
			continue
		}

		resourceName := oldResource.Name()
		newResource := tfwrite.NewEmptyResource(newResourceType, resourceName)
		file.AppendResource(newResource)
		newResource.SetAttributeByReference(newResourceRefAttribute, oldResource, oldResourceRefAttribute)
		newResource.AppendAttribute(attr)
		oldResource.RemoveAttribute(oldAttribute)
	}

	return inFile, nil
}
