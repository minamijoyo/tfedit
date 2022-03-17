package awsv4upgrade

import (
	"github.com/minamijoyo/tfedit/tfeditor"
	"github.com/minamijoyo/tfedit/tfwrite"
)

// AWSS3BucketReplicationConfigurationFilter is a filter implementation for
// upgrading the replication_configuration argument of aws_s3_bucket.
// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#replication_configuration-argument
type AWSS3BucketReplicationConfigurationFilter struct{}

var _ tfeditor.ResourceFilter = (*AWSS3BucketReplicationConfigurationFilter)(nil)

// NewAWSS3BucketReplicationConfigurationFilter creates a new instance of
// AWSS3BucketReplicationConfigurationFilter.
func NewAWSS3BucketReplicationConfigurationFilter() tfeditor.ResourceFilter {
	return &AWSS3BucketReplicationConfigurationFilter{}
}

// ResourceFilter upgrades the replication_configuration argument of
// aws_s3_bucket.
func (f *AWSS3BucketReplicationConfigurationFilter) ResourceFilter(inFile *tfwrite.File, resource *tfwrite.Resource) (*tfwrite.File, error) {
	oldNestedBlock := "replication_configuration"
	newResourceType := "aws_s3_bucket_replication_configuration"

	nestedBlocks := resource.FindNestedBlocksByType(oldNestedBlock)
	if len(nestedBlocks) == 0 {
		return inFile, nil
	}

	resourceName := resource.Name()
	newResource := tfwrite.NewEmptyResource(newResourceType, resourceName)
	inFile.AppendResource(newResource)
	setBucketArgument(newResource, resource)

	for _, nestedBlock := range nestedBlocks {
		roleAttr := nestedBlock.GetAttribute("role")
		if roleAttr != nil {
			newResource.AppendAttribute(roleAttr)
		}

		rulesBlocks := nestedBlock.FindNestedBlocksByType("rules")
		for _, rulesBlock := range rulesBlocks {
			// Rename a `rules` block to a `rule` block
			rulesBlock.SetType("rule")
			newResource.AppendNestedBlock(rulesBlock)

			// Map a `delete_marker_replication_status` attribute to a `delete_marker_replication` block
			// delete_marker_replication_status = "Enabled"
			// =>
			// delete_marker_replication {
			//   status = "Enabled"
			// }
			deleteMarkerAttr := rulesBlock.GetAttribute("delete_marker_replication_status")
			if deleteMarkerAttr != nil {
				deleteMarkerBlock := tfwrite.NewEmptyNestedBlock("delete_marker_replication")
				rulesBlock.AppendNestedBlock(deleteMarkerBlock)
				deleteMarkerBlock.SetAttributeRaw("status", deleteMarkerAttr.ValueAsTokens())
				rulesBlock.RemoveAttribute("delete_marker_replication_status")
			}

			destinationBlocks := rulesBlock.FindNestedBlocksByType("destination")
			for _, destinationBlock := range destinationBlocks {
				// Map a `replication_time.minutes` attribute to a `replication_time.time.minutes` attribute
				// replication_time {
				//   status  = "Enabled"
				//   minutes = 15
				// }
				// =>
				// replication_time {
				//   status = "Enabled"
				//   time {
				//     minutes = 15
				//   }
				// }
				replicationTimeBlocks := destinationBlock.FindNestedBlocksByType("replication_time")
				for _, replicationTimeBlock := range replicationTimeBlocks {
					minutesAttribute := replicationTimeBlock.GetAttribute("minutes")
					if minutesAttribute != nil {
						timeBlock := tfwrite.NewEmptyNestedBlock("time")
						replicationTimeBlock.AppendNestedBlock(timeBlock)
						timeBlock.SetAttributeRaw("minutes", minutesAttribute.ValueAsTokens())
						replicationTimeBlock.RemoveAttribute("minutes")
					}
				}

				// Map a `metrics.minutes` attribute to a `metrics.event_threshold.minutes` attribute
				// metrics {
				//   status  = "Enabled"
				//   minutes = 15
				// }
				// =>
				// metrics {
				//   status = "Enabled"
				//   event_threshold {
				//     minutes = 15
				//   }
				// }
				metricsBlocks := destinationBlock.FindNestedBlocksByType("metrics")
				for _, metricsBlock := range metricsBlocks {
					minutesAttribute := metricsBlock.GetAttribute("minutes")
					if minutesAttribute != nil {
						eventThresholdBlock := tfwrite.NewEmptyNestedBlock("event_threshold")
						metricsBlock.AppendNestedBlock(eventThresholdBlock)
						eventThresholdBlock.SetAttributeRaw("minutes", minutesAttribute.ValueAsTokens())
						metricsBlock.RemoveAttribute("minutes")
					}
				}
			}
		}

		resource.RemoveNestedBlock(nestedBlock)
	}

	return inFile, nil
}
