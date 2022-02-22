package awsv4upgrade

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/minamijoyo/hcledit/editor"
	"github.com/minamijoyo/tfedit/tfeditor"
	"github.com/minamijoyo/tfedit/tfwrite"
)

// AWSS3BucketACLFilter is a filter implementation for upgrading the acl
// argument of aws_s3_bucket.
// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#acl-argument
type AWSS3BucketACLFilter struct{}

var _ editor.Filter = (*AWSS3BucketACLFilter)(nil)
var _ tfeditor.ResourceFilter = (*AWSS3BucketACLFilter)(nil)

// NewAWSS3BucketACLFilter creates a new instance of AWSS3BucketACLFilter.
func NewAWSS3BucketACLFilter() editor.Filter {
	return &AWSS3BucketACLFilter{}
}

// Filter upgrades the acl argument of aws_s3_bucket.
func (f *AWSS3BucketACLFilter) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	m := tfeditor.NewResourcesByTypeFilter("aws_s3_bucket", f)
	return m.Filter(inFile)
}

// ResourceFilter upgrades the acl argument of aws_s3_bucket.
func (f *AWSS3BucketACLFilter) ResourceFilter(inFile *tfwrite.File, resource *tfwrite.Resource) (*tfwrite.File, error) {
	oldAttribute := "acl"
	oldResourceRefAttribute := "id"
	newResourceType := "aws_s3_bucket_acl"
	newResourceRefAttribute := "bucket"

	attr := resource.GetAttribute(oldAttribute)
	if attr == nil {
		return inFile, nil
	}

	resourceName := resource.Name()
	newResource := tfwrite.NewEmptyResource(newResourceType, resourceName)
	inFile.AppendResource(newResource)
	newResource.SetAttributeByReference(newResourceRefAttribute, resource, oldResourceRefAttribute)
	newResource.AppendAttribute(attr)
	resource.RemoveAttribute(oldAttribute)

	return inFile, nil
}
