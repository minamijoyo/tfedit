package awsv4upgrade

import (
	"github.com/minamijoyo/tfedit/tfeditor"
	"github.com/minamijoyo/tfedit/tfwrite"
)

// AWSS3BucketLoggingFilter is a filter implementation for upgrading the
// logging argument of aws_s3_bucket.
// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#logging-argument
type AWSS3BucketLoggingFilter struct{}

var _ tfeditor.ResourceFilter = (*AWSS3BucketLoggingFilter)(nil)

// NewAWSS3BucketLoggingFilter creates a new instance of AWSS3BucketLoggingFilter.
func NewAWSS3BucketLoggingFilter() tfeditor.ResourceFilter {
	return &AWSS3BucketLoggingFilter{}
}

// ResourceFilter upgrades the logging argument of aws_s3_bucket.
func (f *AWSS3BucketLoggingFilter) ResourceFilter(inFile *tfwrite.File, resource *tfwrite.Resource) (*tfwrite.File, error) {
	oldNestedBlock := "logging"
	newResourceType := "aws_s3_bucket_logging"

	nestedBlocks := resource.FindNestedBlocksByType(oldNestedBlock)
	if len(nestedBlocks) == 0 {
		return inFile, nil
	}

	resourceName := resource.Name()
	newResource := tfwrite.NewEmptyResource(newResourceType, resourceName)
	inFile.AppendResource(newResource)
	setParentBucket(newResource, resource)
	newResource.AppendUnwrappedNestedBlockBody(nestedBlocks[0])
	resource.RemoveNestedBlock(nestedBlocks[0])

	return inFile, nil
}
