package awsv4upgrade

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/minamijoyo/hcledit/editor"
	"github.com/minamijoyo/tfedit/tfwrite"
)

// AWSS3BucketLoggingFilter is a filter implementation for upgrading the
// logging argument of aws_s3_bucket.
// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#logging-argument
type AWSS3BucketLoggingFilter struct{}

var _ editor.Filter = (*AWSS3BucketLoggingFilter)(nil)

// NewAWSS3BucketLoggingFilter creates a new instance of AWSS3BucketLoggingFilter.
func NewAWSS3BucketLoggingFilter() editor.Filter {
	return &AWSS3BucketLoggingFilter{}
}

// Filter upgrades the logging argument of aws_s3_bucket.
func (f *AWSS3BucketLoggingFilter) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	file := tfwrite.NewFile(inFile)
	oldResourceType := "aws_s3_bucket"
	oldNestedBlock := "logging"
	oldResourceRefAttribute := "id"
	newResourceType := "aws_s3_bucket_logging"
	newResourceRefAttribute := "bucket"

	targets := file.FindResourcesByType(oldResourceType)
	for _, oldResource := range targets {
		nestedBlock := oldResource.GetFirstNestedBlock(oldNestedBlock)
		if nestedBlock == nil {
			continue
		}

		resourceName := oldResource.Name()
		newResource := tfwrite.NewEmptyResource(newResourceType, resourceName)
		file.AppendResource(newResource)
		newResource.SetAttributeByReference(newResourceRefAttribute, oldResource, oldResourceRefAttribute)
		newResource.AppendUnwrappedNestedBlockBody(nestedBlock)
		oldResource.RemoveNestedBlock(nestedBlock)
	}

	return inFile, nil
}
