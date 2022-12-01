package awsv4upgrade

import (
	"github.com/minamijoyo/tfedit/tfwrite"
	"github.com/zclconf/go-cty/cty"
)

// AWSS3BucketObjectLockConfigurationResourceFilter is a filter implementation
// for upgrading the object_lock_configuration argument of aws_s3_bucket.
// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#object_lock_configuration-rule-argument
func AWSS3BucketObjectLockConfigurationResourceFilter(inFile *tfwrite.File, resource *tfwrite.Resource) (*tfwrite.File, error) {
	oldNestedBlock := "object_lock_configuration"
	newResourceType := "aws_s3_bucket_object_lock_configuration"

	nestedBlocks := resource.FindNestedBlocksByType(oldNestedBlock)
	if len(nestedBlocks) == 0 {
		return inFile, nil
	}

	resourceName := resource.Name()
	newResource := tfwrite.NewEmptyResource(newResourceType, resourceName)
	inFile.AppendBlock(newResource)
	setParentBucket(newResource, resource)

	objectLockBlock := nestedBlocks[0]

	// Move a rule block to a new resource
	ruleBlocks := objectLockBlock.FindNestedBlocksByType("rule")
	for _, ruleBlock := range ruleBlocks {
		newResource.AppendNestedBlock(ruleBlock)
	}

	// Map an `object_lock_configuration.object_lock_enabled` attribute
	// to a top-level `object_lock_enabled` attribute.
	// In addition, the valid type is now bool.
	// object_lock_enabled = "Enabled" => true
	enabledAttr := objectLockBlock.GetAttribute("object_lock_enabled")
	if enabledAttr != nil {
		enabled, err := enabledAttr.ValueAsString()
		if err == nil {
			switch enabled {
			case `"Enabled"`:
				resource.SetAttributeValue("object_lock_enabled", cty.BoolVal(true))
			// case `"Disabled"`:
			// "Disabled" is not defined as an old valid value.
			// The top-level `object_lock_enabled` attribute is optional and the default is false. No op.
			default:
				// If the value is a variable, not literal, we cannot rewrite it automatically.
				// Set original raw tokens as it is.
				resource.SetAttributeRaw("object_lock_enabled", enabledAttr.ValueAsTokens())
			}
		}
	}
	resource.RemoveNestedBlock(objectLockBlock)

	return inFile, nil
}
