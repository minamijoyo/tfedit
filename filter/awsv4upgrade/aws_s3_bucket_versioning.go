package awsv4upgrade

import (
	"github.com/minamijoyo/tfedit/tfwrite"
	"github.com/zclconf/go-cty/cty"
)

// AWSS3BucketVersioningResourceFilter is a filter implementation for upgrading
// the versioning argument of aws_s3_bucket.
// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#versioning-argument
func AWSS3BucketVersioningResourceFilter(inFile *tfwrite.File, resource *tfwrite.Resource) (*tfwrite.File, error) {
	oldNestedBlock := "versioning"
	newResourceType := "aws_s3_bucket_versioning"
	newNestedBlock := "versioning_configuration"

	nestedBlocks := resource.FindNestedBlocksByType(oldNestedBlock)
	if len(nestedBlocks) == 0 {
		return inFile, nil
	}

	resourceName := resource.Name()
	newResource := tfwrite.NewEmptyResource(newResourceType, resourceName)
	inFile.AppendBlock(newResource)
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

	// Map an `mfa_delete` attribute.
	// true => "Enabled"
	// false => "Disabled"
	// There is also the mfa argument in v4, but it seems practically meaningless. Simply ignore it.
	mfaDeleteAttr := nestedBlock.GetAttribute("mfa_delete")
	if mfaDeleteAttr != nil {
		mfaDelete, err := mfaDeleteAttr.ValueAsString()
		if err == nil {
			switch mfaDelete {
			case "true":
				nestedBlock.SetAttributeValue("mfa_delete", cty.StringVal("Enabled"))
			case "false":
				nestedBlock.SetAttributeValue("mfa_delete", cty.StringVal("Disabled"))
			default:
				// If the value is a variable, not literal, we cannot rewrite it automatically.
				// Set original raw tokens as it is.
				nestedBlock.SetAttributeRaw("mfa_delete", mfaDeleteAttr.ValueAsTokens())
			}
		}
	}

	newResource.AppendNestedBlock(nestedBlock)
	resource.RemoveNestedBlock(nestedBlock)

	return inFile, nil
}
