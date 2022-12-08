package awsv4upgrade

import (
	"github.com/minamijoyo/tfedit/tfeditor"
	"github.com/minamijoyo/tfedit/tfwrite"
)

// AWSS3BucketFilter is a filter implementation for upgrading arguments of
// aws_s3_bucket to AWS provider v4.
// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#s3-bucket-refactor
type AWSS3BucketFilter struct {
	filters []tfeditor.BlockFilter
}

var _ tfeditor.BlockFilter = (*AWSS3BucketFilter)(nil)

// NewAWSS3BucketFilter creates a new instance of AWSS3BucketFilter.
func NewAWSS3BucketFilter() tfeditor.BlockFilter {
	filters := []tfeditor.BlockFilter{
		tfeditor.ResourceFilterFunc(AWSS3BucketAccelerationStatusResourceFilter),
		tfeditor.ResourceFilterFunc(AWSS3BucketACLResourceFilter),
		tfeditor.ResourceFilterFunc(AWSS3BucketCorsRuleResourceFilter),
		tfeditor.ResourceFilterFunc(AWSS3BucketGrantResourceFilter),
		tfeditor.ResourceFilterFunc(AWSS3BucketLifecycleRuleResourceFilter),
		tfeditor.ResourceFilterFunc(AWSS3BucketLoggingResourceFilter),
		tfeditor.ResourceFilterFunc(AWSS3BucketObjectLockConfigurationResourceFilter),
		tfeditor.ResourceFilterFunc(AWSS3BucketPolicyResourceFilter),
		tfeditor.ResourceFilterFunc(AWSS3BucketReplicationConfigurationResourceFilter),
		tfeditor.ResourceFilterFunc(AWSS3BucketRequestPayerResourceFilter),
		tfeditor.ResourceFilterFunc(AWSS3BucketServerSideEncryptionConfigurationResourceFilter),
		tfeditor.ResourceFilterFunc(AWSS3BucketVersioningResourceFilter),
		tfeditor.ResourceFilterFunc(AWSS3BucketWebsiteResourceFilter),

		// Remove redundant TokenNewLine tokens in the resource block after removing nested blocks.
		// Since VerticalFormat clears tokens internally, we should call it at the end.
		tfeditor.BlockFilterFunc(tfeditor.VerticalFormatterFilter),
	}

	return &AWSS3BucketFilter{filters: filters}
}

// BlockFilter upgrades arguments of aws_s3_bucket to AWS provider v4.
// Some rules have not been implemented yet.
func (f *AWSS3BucketFilter) BlockFilter(inFile *tfwrite.File, block tfwrite.Block) (*tfwrite.File, error) {
	m := tfeditor.NewMultiBlockFilter(f.filters)
	return m.BlockFilter(inFile, block)
}

// setParentBucket is a helper method for setting the followings:
// - copy provider, count and for_each meta arguments
// - set a bucket argument of a new `aws_s3_bucket_*` resource to the original `aws_s3_bucket` resource.
func setParentBucket(newResource *tfwrite.Resource, oldResource *tfwrite.Resource) {
	// copy provider, count and for_each meta arguments
	newResource.CopyAttribute(oldResource, "provider")
	newResource.CopyAttribute(oldResource, "count")
	newResource.CopyAttribute(oldResource, "for_each")

	// set a bucket argument
	newResource.SetAttributeByReference("bucket", oldResource, "id")
}
