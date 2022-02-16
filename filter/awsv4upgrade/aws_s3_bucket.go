package awsv4upgrade

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/minamijoyo/hcledit/editor"
)

// AWSS3BucketFilter is a filter implementation for upgrading arguments of
// aws_s3_bucket to AWS provider v4.
// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#s3-bucket-refactor
type AWSS3BucketFilter struct {
}

var _ editor.Filter = (*AWSS3BucketFilter)(nil)

// NewAWSS3BucketFilter creates a new instance of AWSS3BucketFilter.
func NewAWSS3BucketFilter() editor.Filter {
	return &AWSS3BucketFilter{}
}

// Filter upgrades arguments of aws_s3_bucket to AWS provider v4.
// Some rules have not been implemented yet.
func (f *AWSS3BucketFilter) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	m := editor.NewMultiFilter([]editor.Filter{
		// &AWSS3BucketAccelerationStatusFilter{},
		&AWSS3BucketACLFilter{},
		// &AWSS3BucketCorsRuleFilter{},
		// &AWSS3BucketGrantFilter{},
		// &AWSS3BucketLifecycleRuleFilter{},
		&AWSS3BucketLoggingFilter{},
		// &AWSS3BucketObjectLockConfigurationFilter{},
		// &AWSS3BucketPolicyFilter{},
		// &AWSS3BucketReplicationConfigurationFilter{},
		// &AWSS3BucketRequestPayerFilter{},
		// &AWSS3BucketServerSideEncryptionConfigurationFilter{},
		// &AWSS3BucketVersioningFilter{},
		// &AWSS3BucketWebsiteFilter{},
	})
	return m.Filter(inFile)
}

// setBucketArgument is a helper method for setting a bucket argument to the given block.
func setBucketArgument(block *hclwrite.Block, resourceName string) *hclwrite.Attribute {
	return block.Body().SetAttributeTraversal("bucket", hcl.Traversal{
		hcl.TraverseRoot{Name: "aws_s3_bucket"},
		hcl.TraverseAttr{Name: resourceName},
		hcl.TraverseAttr{Name: "id"},
	})
}
