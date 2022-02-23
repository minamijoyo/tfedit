package awsv4upgrade

import (
	"github.com/minamijoyo/tfedit/tfeditor"
	"github.com/minamijoyo/tfedit/tfwrite"
	"github.com/zclconf/go-cty/cty"
)

// AWSS3BucketLifecycleRuleFilter is a filter implementation for upgrading the
// lifecycle_rule argument of aws_s3_bucket.
// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#lifecycle_rule-argument
type AWSS3BucketLifecycleRuleFilter struct{}

var _ tfeditor.ResourceFilter = (*AWSS3BucketLifecycleRuleFilter)(nil)

// NewAWSS3BucketLifecycleRuleFilter creates a new instance of AWSS3Bucketlifecycle_ruleFilter.
func NewAWSS3BucketLifecycleRuleFilter() tfeditor.ResourceFilter {
	return &AWSS3BucketLifecycleRuleFilter{}
}

// ResourceFilter upgrades the lifecycle_rule argument of aws_s3_bucket.
func (f *AWSS3BucketLifecycleRuleFilter) ResourceFilter(inFile *tfwrite.File, resource *tfwrite.Resource) (*tfwrite.File, error) {
	oldNestedBlock := "lifecycle_rule"
	oldResourceRefAttribute := "id"
	newResourceType := "aws_s3_bucket_lifecycle_configuration"
	newNestedBlock := "rule"
	newResourceRefAttribute := "bucket"

	nestedBlocks := resource.FindNestedBlocksByType(oldNestedBlock)
	if len(nestedBlocks) == 0 {
		return inFile, nil
	}

	resourceName := resource.Name()
	newResource := tfwrite.NewEmptyResource(newResourceType, resourceName)
	inFile.AppendResource(newResource)
	newResource.SetAttributeByReference(newResourceRefAttribute, resource, oldResourceRefAttribute)

	for _, nestedBlock := range nestedBlocks {
		// Rename a `lifecycle_rule` block to a `rule` block
		nestedBlock.SetType(newNestedBlock)

		// Map an `enabled` attribute to a `status` attribute
		// enabled = true => status = "Enabled"
		// enabled = false => status = "Disabled"
		enabledAttr := nestedBlock.GetAttribute("enabled")
		if enabledAttr != nil {
			enabled, err := enabledAttr.ValueAsString()
			if err == nil {
				switch enabled {
				case "true":
					nestedBlock.SetAttributeValue("status", cty.StringVal("Enabled"))
				case "false":
					nestedBlock.SetAttributeValue("status", cty.StringVal("Disabled"))
				default:
					// If the value is a variable, not literal, we cannot rewrite it automatically.
					// Set original raw tokens as it is.
					nestedBlock.SetAttributeRaw("status", enabledAttr.ValueAsTokens())
				}
			}
			nestedBlock.RemoveAttribute("enabled")
		}

		// Map a prefix attribute to a filter block
		// prefix  = "tmp/"
		// =>
		// filter {
		//   prefix  = "tmp/"
		// }
		prefixAttr := nestedBlock.GetAttribute("prefix")
		filterBlock := tfwrite.NewEmptyNestedBlock("filter")
		nestedBlock.AppendNestedBlock(filterBlock)
		if prefixAttr != nil {
			filterBlock.SetAttributeRaw("prefix", prefixAttr.ValueAsTokens())
		} else {
			// If a prefix attribute is not found, set an empty string by default.
			// According to the upgrade guide,
			// when aws s3api get-bucket-lifecycle-configuration returns `"Filter" : {}`,
			// we should not set prefix, however we cannot know it without an API call,
			// so we just assume it contains `"Filter" : { "Prefix": "" }` here.
			filterBlock.SetAttributeValue("prefix", cty.StringVal(""))
		}
		nestedBlock.RemoveAttribute("prefix")

		// Rename a days attribute in noncurrent_version_transition to noncurrent_days.
		transitionBlocks := nestedBlock.FindNestedBlocksByType("noncurrent_version_transition")
		for _, transitionBlock := range transitionBlocks {
			daysAttr := transitionBlock.GetAttribute("days")
			if daysAttr != nil {
				transitionBlock.SetAttributeRaw("noncurrent_days", daysAttr.ValueAsTokens())
				transitionBlock.RemoveAttribute("days")
			}
		}

		// Rename a days attribute in noncurrent_version_expiration to noncurrent_days.
		expirationBlocks := nestedBlock.FindNestedBlocksByType("noncurrent_version_expiration")
		for _, expirationBlock := range expirationBlocks {
			daysAttr := expirationBlock.GetAttribute("days")
			if daysAttr != nil {
				expirationBlock.SetAttributeRaw("noncurrent_days", daysAttr.ValueAsTokens())
				expirationBlock.RemoveAttribute("days")
			}
		}

		// Map a abort_incomplete_multipart_upload_days attribute to a abort_incomplete_multipart_upload block
		// abort_incomplete_multipart_upload_days = 7
		// =>
		// abort_incomplete_multipart_upload {
		//   days_after_initiation = 7
		// }
		abortAttr := nestedBlock.GetAttribute("abort_incomplete_multipart_upload_days")
		if abortAttr != nil {
			abortBlock := tfwrite.NewEmptyNestedBlock("abort_incomplete_multipart_upload")
			nestedBlock.AppendNestedBlock(abortBlock)
			abortBlock.SetAttributeRaw("days_after_initiation", abortAttr.ValueAsTokens())
			nestedBlock.RemoveAttribute("abort_incomplete_multipart_upload_days")
		}

		newResource.AppendNestedBlock(nestedBlock)
		resource.RemoveNestedBlock(nestedBlock)
	}

	return inFile, nil
}
