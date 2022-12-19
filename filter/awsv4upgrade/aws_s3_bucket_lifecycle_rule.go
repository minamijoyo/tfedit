package awsv4upgrade

import (
	"fmt"
	"strings"
	"time"

	"github.com/minamijoyo/tfedit/tfwrite"
	"github.com/zclconf/go-cty/cty"
)

// AWSS3BucketLifecycleRuleResourceFilter is a filter implementation for
// upgrading the lifecycle_rule argument of aws_s3_bucket.
// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#lifecycle_rule-argument
func AWSS3BucketLifecycleRuleResourceFilter(inFile *tfwrite.File, resource *tfwrite.Resource) (*tfwrite.File, error) {
	if resource.SchemaType() != "aws_s3_bucket" {
		return inFile, nil
	}

	oldNestedBlock := "lifecycle_rule"
	newResourceType := "aws_s3_bucket_lifecycle_configuration"
	newNestedBlock := "rule"

	nestedBlocks := resource.FindNestedBlocksByType(oldNestedBlock)
	if len(nestedBlocks) == 0 {
		return inFile, nil
	}

	resourceName := resource.Name()
	newResource := tfwrite.NewEmptyResource(newResourceType, resourceName)
	inFile.AppendBlock(newResource)
	setParentBucket(newResource, resource)

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
		nestedBlock.AppendNestedBlock(filterBlock)
		prefixAttr := nestedBlock.GetAttribute("prefix")
		tagsAttr := nestedBlock.GetAttribute("tags")
		if tagsAttr != nil {
			tags, err := tagsAttr.ValueAsString()
			if err == nil {
				if tags != "{}" {
					// Non empty tags should be wrapped by an `and` block.
					andBlock := tfwrite.NewEmptyNestedBlock("and")
					filterBlock.AppendNestedBlock(andBlock)
					if prefixAttr != nil {
						andBlock.SetAttributeRaw("prefix", prefixAttr.ValueAsTokens())
						nestedBlock.RemoveAttribute("prefix")
					} else {
						// If a prefix attribute is not found, set an empty string by default.
						andBlock.SetAttributeValue("prefix", cty.StringVal(""))
					}
					andBlock.SetAttributeRaw("tags", tagsAttr.ValueAsTokens())
				} else {
					// When both prefix and tags were empty but defined, it will result
					// in a migration plan diff, so remove them and put an empty filter.
					if prefixAttr != nil {
						prefix, err := prefixAttr.ValueAsString()
						if err == nil && prefix == `""` {
							nestedBlock.RemoveAttribute("prefix")
						}
					}
				}
			}
			nestedBlock.RemoveAttribute("tags")
		} else {
			if prefixAttr != nil {
				filterBlock.SetAttributeRaw("prefix", prefixAttr.ValueAsTokens())
				nestedBlock.RemoveAttribute("prefix")
			} else {
				// If a prefix attribute is not found, set an empty string by default.
				// According to the upgrade guide,
				// when aws s3api get-bucket-lifecycle-configuration returns `"Filter" : {}`,
				// we should not set prefix, however we cannot know it without an API call,
				// so we just assume it contains `"Filter" : { "Prefix": "" }` here.
				filterBlock.SetAttributeValue("prefix", cty.StringVal(""))
			}
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
			abort, err := abortAttr.ValueAsString()
			// When the value is 0, adding a block will result in a migration plan diff,
			// so suppress adding the block.
			if err == nil && abort != "0" {
				abortBlock := tfwrite.NewEmptyNestedBlock("abort_incomplete_multipart_upload")
				nestedBlock.AppendNestedBlock(abortBlock)
				abortBlock.SetAttributeRaw("days_after_initiation", abortAttr.ValueAsTokens())
			}
			nestedBlock.RemoveAttribute("abort_incomplete_multipart_upload_days")
		}

		newResource.AppendNestedBlock(nestedBlock)
		resource.RemoveNestedBlock(nestedBlock)
	}

	return inFile, nil
}
