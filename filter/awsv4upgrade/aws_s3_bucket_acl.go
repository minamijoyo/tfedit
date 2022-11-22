package awsv4upgrade

import (
	"github.com/minamijoyo/tfedit/tfwrite"
)

// AWSS3BucketACLResourceFilter is a filter implementation for upgrading the
// acl argument of aws_s3_bucket.
// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#acl-argument
func AWSS3BucketACLResourceFilter(inFile *tfwrite.File, resource *tfwrite.Resource) (*tfwrite.File, error) {
	oldAttribute := "acl"
	newResourceType := "aws_s3_bucket_acl"

	attr := resource.GetAttribute(oldAttribute)
	if attr == nil {
		return inFile, nil
	}

	resourceName := resource.Name()
	newResource := tfwrite.NewEmptyResource(newResourceType, resourceName)
	inFile.AppendResource(newResource)
	setParentBucket(newResource, resource)
	newResource.AppendAttribute(attr)
	resource.RemoveAttribute(oldAttribute)

	return inFile, nil
}
