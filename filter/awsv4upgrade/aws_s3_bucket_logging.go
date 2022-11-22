package awsv4upgrade

import (
	"github.com/minamijoyo/tfedit/tfwrite"
)

// AWSS3BucketLoggingResourceFilter is a filter implementation for upgrading
// the logging argument of aws_s3_bucket.
// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#logging-argument
func AWSS3BucketLoggingResourceFilter(inFile *tfwrite.File, resource *tfwrite.Resource) (*tfwrite.File, error) {
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
