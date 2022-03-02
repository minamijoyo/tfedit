package awsv4upgrade

import (
	"fmt"
	"strings"
	"time"

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
	newResourceType := "aws_s3_bucket_lifecycle_configuration"
	newNestedBlock := "rule"

	nestedBlocks := resource.FindNestedBlocksByType(oldNestedBlock)
	if len(nestedBlocks) == 0 {
		return inFile, nil
	}

	resourceName := resource.Name()
	newResource := tfwrite.NewEmptyResource(newResourceType, resourceName)
	inFile.AppendResource(newResource)
	setBucketArgument(newResource, resource)

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

		// Map a prefix attribute to a filter block without tags
		// prefix  = "tmp/"
		// =>
		// filter {
		//   prefix  = "tmp/"
		// }
		// Map a prefix attribute to a filter block with tags
		// prefix  = "tmp/"
		// tags = {
		//   rule      = "log"
		//   autoclean = "true"
		// }
		// =>
		// filter {
		//   and {
		//     prefix  = "tmp/"
		//     tags = {
		//       rule      = "log"
		//       autoclean = "true"
		//     }
		//   }
		// }
		// Create a filter block
		filterBlock := tfwrite.NewEmptyNestedBlock("filter")
		// A child block points a `filter` or an `and` block.
		var childBlock *tfwrite.NestedBlock
		nestedBlock.AppendNestedBlock(filterBlock)
		prefixAttr := nestedBlock.GetAttribute("prefix")
		tagsAttr := nestedBlock.GetAttribute("tags")
		if tagsAttr != nil && prefixAttr != nil {
			andBlock := tfwrite.NewEmptyNestedBlock("and")
			filterBlock.AppendNestedBlock(andBlock)
			childBlock = andBlock
		} else {
			childBlock = filterBlock
		}
		if prefixAttr != nil {
			childBlock.SetAttributeRaw("prefix", prefixAttr.ValueAsTokens())
			nestedBlock.RemoveAttribute("prefix")
		} else {
			// If a prefix attribute is not found, set an empty string by default.
			// According to the upgrade guide,
			// when aws s3api get-bucket-lifecycle-configuration returns `"Filter" : {}`,
			// we should not set prefix, however we cannot know it without an API call,
			// so we just assume it contains `"Filter" : { "Prefix": "" }` here.
			childBlock.SetAttributeValue("prefix", cty.StringVal(""))
		}
		if tagsAttr != nil {
			childBlock.SetAttributeRaw("tags", tagsAttr.ValueAsTokens())
			nestedBlock.RemoveAttribute("tags")
		}

		// Convert a timestamp format for a date attribute in transition.
		// date = "2022-12-31" => date = "2022-12-31T00:00:00Z"
		transitionBlocks := nestedBlock.FindNestedBlocksByType("transition")
		for _, transitionBlock := range transitionBlocks {
			dateAttr := transitionBlock.GetAttribute("date")
			if dateAttr != nil {
				date, err := dateAttr.ValueAsString()
				if err == nil {
					unquotedDate := strings.Trim(date, "\"")
					// Try to parse date as an old format in v3.
					_, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT00:00:00Z", unquotedDate))
					if err != nil {
						// If failed to parse, we assume that the value is a variable, not literal,
						// we cannot rewrite it automatically, so keep original raw tokens as it is.
						continue
					}
					// If the value has a string literal with valid format in v3,
					// covert it to a new formart in v4.
					newDate := fmt.Sprintf("%sT00:00:00Z", unquotedDate)
					transitionBlock.SetAttributeValue("date", cty.StringVal(newDate))
				}
			}
		}

		// Convert a timestamp format for a date attribute in expiration.
		// date = "2022-12-31" => date = "2022-12-31T00:00:00Z"
		expirationBlocks := nestedBlock.FindNestedBlocksByType("expiration")
		for _, expirationBlock := range expirationBlocks {
			dateAttr := expirationBlock.GetAttribute("date")
			if dateAttr != nil {
				date, err := dateAttr.ValueAsString()
				if err == nil {
					unquotedDate := strings.Trim(date, "\"")
					// Try to parse date as an old format in v3.
					_, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT00:00:00Z", unquotedDate))
					if err != nil {
						// If failed to parse, we assume that the value is a variable, not literal,
						// we cannot rewrite it automatically, so keep original raw tokens as it is.
						continue
					}
					// If the value has a string literal with valid format in v3,
					// covert it to a new formart in v4.
					newDate := fmt.Sprintf("%sT00:00:00Z", unquotedDate)
					expirationBlock.SetAttributeValue("date", cty.StringVal(newDate))
				}
			}
		}

		// Rename a days attribute in noncurrent_version_transition to noncurrent_days.
		noncurrentVersionTransitionBlocks := nestedBlock.FindNestedBlocksByType("noncurrent_version_transition")
		for _, noncurrentVersionTransitionBlock := range noncurrentVersionTransitionBlocks {
			daysAttr := noncurrentVersionTransitionBlock.GetAttribute("days")
			if daysAttr != nil {
				noncurrentVersionTransitionBlock.SetAttributeRaw("noncurrent_days", daysAttr.ValueAsTokens())
				noncurrentVersionTransitionBlock.RemoveAttribute("days")
			}
		}

		// Rename a days attribute in noncurrent_version_expiration to noncurrent_days.
		noncurrentVersionExpirationBlocks := nestedBlock.FindNestedBlocksByType("noncurrent_version_expiration")
		for _, noncurrentVersionExpirationBlock := range noncurrentVersionExpirationBlocks {
			daysAttr := noncurrentVersionExpirationBlock.GetAttribute("days")
			if daysAttr != nil {
				noncurrentVersionExpirationBlock.SetAttributeRaw("noncurrent_days", daysAttr.ValueAsTokens())
				noncurrentVersionExpirationBlock.RemoveAttribute("days")
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
