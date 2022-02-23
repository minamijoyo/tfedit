package awsv4upgrade

import (
	"github.com/minamijoyo/tfedit/tfeditor"
	"github.com/minamijoyo/tfedit/tfwrite"
)

// AWSS3BucketACLFilter is a filter implementation for upgrading the acl
// argument of aws_s3_bucket.
// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#acl-argument
type AWSS3BucketACLFilter struct{}

var _ tfeditor.ResourceFilter = (*AWSS3BucketACLFilter)(nil)

// NewAWSS3BucketACLFilter creates a new instance of AWSS3BucketACLFilter.
func NewAWSS3BucketACLFilter() tfeditor.ResourceFilter {
	return &AWSS3BucketACLFilter{}
}

// ResourceFilter upgrades the acl argument of aws_s3_bucket.
func (f *AWSS3BucketACLFilter) ResourceFilter(inFile *tfwrite.File, resource *tfwrite.Resource) (*tfwrite.File, error) {
	oldAttribute := "acl"
	newResourceType := "aws_s3_bucket_acl"

	attr := resource.GetAttribute(oldAttribute)
	if attr == nil {
		return inFile, nil
	}

	resourceName := resource.Name()
	newResource := tfwrite.NewEmptyResource(newResourceType, resourceName)
	inFile.AppendResource(newResource)
	setBucketArgument(newResource, resource)
	newResource.AppendAttribute(attr)
	resource.RemoveAttribute(oldAttribute)

	return inFile, nil
}
