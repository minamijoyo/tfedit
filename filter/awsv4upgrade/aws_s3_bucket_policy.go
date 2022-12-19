package awsv4upgrade

import (
	"github.com/minamijoyo/tfedit/tfwrite"
)

// AWSS3BucketPolicyResourceFilter is a filter implementation for upgrading the
// policy argument of aws_s3_bucket.
// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#policy-argument
func AWSS3BucketPolicyResourceFilter(inFile *tfwrite.File, resource *tfwrite.Resource) (*tfwrite.File, error) {
	if resource.SchemaType() != "aws_s3_bucket" {
		return inFile, nil
	}

	oldAttribute := "policy"
	newResourceType := "aws_s3_bucket_policy"

	attr := resource.GetAttribute(oldAttribute)
	if attr == nil {
		return inFile, nil
	}

	resourceName := resource.Name()
	newResource := tfwrite.NewEmptyResource(newResourceType, resourceName)
	inFile.AppendBlock(newResource)
	setParentBucket(newResource, resource)
	newResource.AppendAttribute(attr)
	resource.RemoveAttribute(oldAttribute)

	return inFile, nil
}
