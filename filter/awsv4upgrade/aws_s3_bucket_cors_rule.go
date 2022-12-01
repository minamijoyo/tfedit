package awsv4upgrade

import (
	"github.com/minamijoyo/tfedit/tfwrite"
)

// AWSS3BucketCorsRuleResourceFilter is a filter implementation for upgrading
// the cors_rule argument of aws_s3_bucket.
// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#cors_rule-argument
func AWSS3BucketCorsRuleResourceFilter(inFile *tfwrite.File, resource *tfwrite.Resource) (*tfwrite.File, error) {
	oldNestedBlock := "cors_rule"
	newResourceType := "aws_s3_bucket_cors_configuration"

	nestedBlocks := resource.FindNestedBlocksByType(oldNestedBlock)
	if len(nestedBlocks) == 0 {
		return inFile, nil
	}

	resourceName := resource.Name()
	newResource := tfwrite.NewEmptyResource(newResourceType, resourceName)
	inFile.AppendBlock(newResource)
	setParentBucket(newResource, resource)

	for _, nestedBlock := range nestedBlocks {
		newResource.AppendNestedBlock(nestedBlock)
		resource.RemoveNestedBlock(nestedBlock)
	}

	return inFile, nil
}
