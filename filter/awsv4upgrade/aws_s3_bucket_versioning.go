package awsv4upgrade

import (
	"github.com/minamijoyo/tfedit/tfeditor"
	"github.com/minamijoyo/tfedit/tfwrite"
	"github.com/zclconf/go-cty/cty"
)

// AWSS3BucketVersioningFilter is a filter implementation for upgrading the
// versioning argument of aws_s3_bucket.
// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#versioning-argument
type AWSS3BucketVersioningFilter struct{}

var _ tfeditor.ResourceFilter = (*AWSS3BucketVersioningFilter)(nil)

// NewAWSS3BucketVersioningFilter creates a new instance of AWSS3BucketVersioningFilter.
func NewAWSS3BucketVersioningFilter() tfeditor.ResourceFilter {
	return &AWSS3BucketVersioningFilter{}
}

// ResourceFilter upgrades the versioning argument of aws_s3_bucket.
func (f *AWSS3BucketVersioningFilter) ResourceFilter(inFile *tfwrite.File, resource *tfwrite.Resource) (*tfwrite.File, error) {
	oldNestedBlock := "versioning"
	newResourceType := "aws_s3_bucket_versioning"
	newNestedBlock := "versioning_configuration"

	nestedBlocks := resource.FindNestedBlocksByType(oldNestedBlock)
	if len(nestedBlocks) == 0 {
		return inFile, nil
	}

	resourceName := resource.Name()
	newResource := tfwrite.NewEmptyResource(newResourceType, resourceName)
	inFile.AppendResource(newResource)
	setParentBucket(newResource, resource)

	nestedBlock := nestedBlocks[0]

	// Rename a `versioning` block to a `versioning_configuration` block
	nestedBlock.SetType(newNestedBlock)

	// Map an `enabled` attribute to a `status` attribute
	// enabled = true => status = "Enabled"
	// enabled = false => status = "Suspended"
	enabledAttr := nestedBlock.GetAttribute("enabled")
	if enabledAttr != nil {
		enabled, err := enabledAttr.ValueAsString()
		if err == nil {
			switch enabled {
			case "true":
				nestedBlock.SetAttributeValue("status", cty.StringVal("Enabled"))
			case "false":
				nestedBlock.SetAttributeValue("status", cty.StringVal("Suspended"))
			default:
				// If the value is a variable, not literal, we cannot rewrite it automatically.
				// Set original raw tokens as it is.
				nestedBlock.SetAttributeRaw("status", enabledAttr.ValueAsTokens())
			}
		}
		nestedBlock.RemoveAttribute("enabled")
	}

	newResource.AppendNestedBlock(nestedBlock)
	resource.RemoveNestedBlock(nestedBlock)

	return inFile, nil
}
