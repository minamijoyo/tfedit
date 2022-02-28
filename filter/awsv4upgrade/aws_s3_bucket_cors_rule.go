package awsv4upgrade

import (
	"github.com/minamijoyo/tfedit/tfeditor"
	"github.com/minamijoyo/tfedit/tfwrite"
)

// AWSS3BucketCorsRuleFilter is a filter implementation for upgrading the
// cors_rule argument of aws_s3_bucket.
// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#cors_rule-argument
type AWSS3BucketCorsRuleFilter struct{}

var _ tfeditor.ResourceFilter = (*AWSS3BucketCorsRuleFilter)(nil)

// NewAWSS3BucketCorsRuleFilter creates a new instance of AWSS3BucketCorsRuleFilter.
func NewAWSS3BucketCorsRuleFilter() tfeditor.ResourceFilter {
	return &AWSS3BucketCorsRuleFilter{}
}

// ResourceFilter upgrades the cors_rule argument of aws_s3_bucket.
func (f *AWSS3BucketCorsRuleFilter) ResourceFilter(inFile *tfwrite.File, resource *tfwrite.Resource) (*tfwrite.File, error) {
	oldNestedBlock := "cors_rule"
	newResourceType := "aws_s3_bucket_cors_configuration"

	nestedBlocks := resource.FindNestedBlocksByType(oldNestedBlock)
	if len(nestedBlocks) == 0 {
		return inFile, nil
	}

	resourceName := resource.Name()
	newResource := tfwrite.NewEmptyResource(newResourceType, resourceName)
	inFile.AppendResource(newResource)
	setBucketArgument(newResource, resource)

	for _, nestedBlock := range nestedBlocks {
		newResource.AppendNestedBlock(nestedBlock)
		resource.RemoveNestedBlock(nestedBlock)
	}

	return inFile, nil
}
