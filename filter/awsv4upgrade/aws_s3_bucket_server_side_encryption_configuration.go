package awsv4upgrade

import (
	"github.com/minamijoyo/tfedit/tfwrite"
)

// AWSS3BucketServerSideEncryptionConfigurationResourceFilter is a filter
// implementation for upgrading the server_side_encryption_configuration
// argument of aws_s3_bucket.
// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#server_side_encryption_configuration-argument
func AWSS3BucketServerSideEncryptionConfigurationResourceFilter(inFile *tfwrite.File, resource *tfwrite.Resource) (*tfwrite.File, error) {
	oldNestedBlock := "server_side_encryption_configuration"
	newResourceType := "aws_s3_bucket_server_side_encryption_configuration"

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
