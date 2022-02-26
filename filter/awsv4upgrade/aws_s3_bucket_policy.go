package awsv4upgrade

import (
	"github.com/minamijoyo/tfedit/tfeditor"
	"github.com/minamijoyo/tfedit/tfwrite"
)

// AWSS3BucketPolicyFilter is a filter implementation for upgrading the policy
// argument of aws_s3_bucket.
// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#policy-argument
type AWSS3BucketPolicyFilter struct{}

var _ tfeditor.ResourceFilter = (*AWSS3BucketPolicyFilter)(nil)

// NewAWSS3BucketPolicyFilter creates a new instance of AWSS3BucketPolicyFilter.
func NewAWSS3BucketPolicyFilter() tfeditor.ResourceFilter {
	return &AWSS3BucketPolicyFilter{}
}

// ResourceFilter upgrades the policy argument of aws_s3_bucket.
func (f *AWSS3BucketPolicyFilter) ResourceFilter(inFile *tfwrite.File, resource *tfwrite.Resource) (*tfwrite.File, error) {
	oldAttribute := "policy"
	newResourceType := "aws_s3_bucket_policy"

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
