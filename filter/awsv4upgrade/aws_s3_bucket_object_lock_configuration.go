package awsv4upgrade

import (
	"github.com/minamijoyo/tfedit/tfeditor"
	"github.com/minamijoyo/tfedit/tfwrite"
)

// AWSS3BucketObjectLockConfigurationFilter is a filter implementation for
// upgrading the object_lock_configuration argument of aws_s3_bucket.
// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#object_lock_configuration-rule-argument
type AWSS3BucketObjectLockConfigurationFilter struct{}

var _ tfeditor.ResourceFilter = (*AWSS3BucketObjectLockConfigurationFilter)(nil)

// NewAWSS3BucketObjectLockConfigurationFilter creates a new instance of
// AWSS3BucketObjectLockConfigurationFilter.
func NewAWSS3BucketObjectLockConfigurationFilter() tfeditor.ResourceFilter {
	return &AWSS3BucketObjectLockConfigurationFilter{}
}

// ResourceFilter upgrades the object_lock_configuration argument of
// aws_s3_bucket.
func (f *AWSS3BucketObjectLockConfigurationFilter) ResourceFilter(inFile *tfwrite.File, resource *tfwrite.Resource) (*tfwrite.File, error) {
	oldNestedBlock := "object_lock_configuration"
	newResourceType := "aws_s3_bucket_object_lock_configuration"

	nestedBlocks := resource.FindNestedBlocksByType(oldNestedBlock)
	if len(nestedBlocks) == 0 {
		return inFile, nil
	}

	resourceName := resource.Name()
	newResource := tfwrite.NewEmptyResource(newResourceType, resourceName)
	inFile.AppendResource(newResource)
	setBucketArgument(newResource, resource)

	objectLockBlock := nestedBlocks[0]

	// Move a rule block to a new resource
	ruleBlocks := objectLockBlock.FindNestedBlocksByType("rule")
	for _, ruleBlock := range ruleBlocks {
		newResource.AppendNestedBlock(ruleBlock)
		objectLockBlock.RemoveNestedBlock(ruleBlock)
	}

	return inFile, nil
}
