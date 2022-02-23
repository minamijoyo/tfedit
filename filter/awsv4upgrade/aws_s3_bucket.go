package awsv4upgrade

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/minamijoyo/hcledit/editor"
	"github.com/minamijoyo/tfedit/tfeditor"
	"github.com/minamijoyo/tfedit/tfwrite"
)

// AWSS3BucketFilter is a filter implementation for upgrading arguments of
// aws_s3_bucket to AWS provider v4.
// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#s3-bucket-refactor
type AWSS3BucketFilter struct {
	filters []tfeditor.ResourceFilter
}

var _ editor.Filter = (*AWSS3BucketFilter)(nil)

// NewAWSS3BucketFilter creates a new instance of AWSS3BucketFilter.
func NewAWSS3BucketFilter() editor.Filter {
	filters := []tfeditor.ResourceFilter{
		// &AWSS3BucketAccelerationStatusFilter{},
		&AWSS3BucketACLFilter{},
		// &AWSS3BucketCorsRuleFilter{},
		// &AWSS3BucketGrantFilter{},
		&AWSS3BucketLifecycleRuleFilter{},
		&AWSS3BucketLoggingFilter{},
		// &AWSS3BucketObjectLockConfigurationFilter{},
		// &AWSS3BucketPolicyFilter{},
		// &AWSS3BucketReplicationConfigurationFilter{},
		// &AWSS3BucketRequestPayerFilter{},
		// &AWSS3BucketServerSideEncryptionConfigurationFilter{},
		// &AWSS3BucketVersioningFilter{},
		// &AWSS3BucketWebsiteFilter{},

		// Remove redundant TokenNewLine tokens in the resource block after removing nested blocks.
		// Since VerticalFormat clears tokens internally, we should call it at the end.
		tfeditor.NewVerticalFormatterResourceFilter(),
	}
	return &AWSS3BucketFilter{filters: filters}
}

// Filter upgrades arguments of aws_s3_bucket to AWS provider v4.
// Some rules have not been implemented yet.
func (f *AWSS3BucketFilter) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	m := tfeditor.NewResourcesByTypeFilter("aws_s3_bucket", f)
	return m.Filter(inFile)
}

// ResourceFilter upgrades arguments of aws_s3_bucket to AWS provider v4.
// Some rules have not been implemented yet.
func (f *AWSS3BucketFilter) ResourceFilter(inFile *tfwrite.File, resource *tfwrite.Resource) (*tfwrite.File, error) {
	m := tfeditor.NewMultiResourceFilter(f.filters)
	return m.ResourceFilter(inFile, resource)
}

// setBucketArgument is a helper method for setting a bucket argument to the given block.
func setBucketArgument(block *hclwrite.Block, resourceName string) *hclwrite.Attribute {
	return block.Body().SetAttributeTraversal("bucket", hcl.Traversal{
		hcl.TraverseRoot{Name: "aws_s3_bucket"},
		hcl.TraverseAttr{Name: resourceName},
		hcl.TraverseAttr{Name: "id"},
	})
}
