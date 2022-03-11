package awsv4upgrade

import (
	"github.com/minamijoyo/tfedit/tfeditor"
	"github.com/minamijoyo/tfedit/tfwrite"
)

// AWSS3BucketAccelerationStatusFilter is a filter implementation for upgrading
// the acceleration_status argument of aws_s3_bucket.
// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#acceleration_status-argument
type AWSS3BucketAccelerationStatusFilter struct{}

var _ tfeditor.ResourceFilter = (*AWSS3BucketAccelerationStatusFilter)(nil)

// NewAWSS3BucketAccelerationStatusFilter creates a new instance of AWSS3BucketAccelerationStatusFilter.
func NewAWSS3BucketAccelerationStatusFilter() tfeditor.ResourceFilter {
	return &AWSS3BucketAccelerationStatusFilter{}
}

// ResourceFilter upgrades the acceleration_status argument of aws_s3_bucket.
func (f *AWSS3BucketAccelerationStatusFilter) ResourceFilter(inFile *tfwrite.File, resource *tfwrite.Resource) (*tfwrite.File, error) {
	oldAttribute := "acceleration_status"
	newResourceType := "aws_s3_bucket_accelerate_configuration"

	attr := resource.GetAttribute(oldAttribute)
	if attr == nil {
		return inFile, nil
	}

	resourceName := resource.Name()
	newResource := tfwrite.NewEmptyResource(newResourceType, resourceName)
	inFile.AppendResource(newResource)
	setBucketArgument(newResource, resource)

	// Map an `acceleration_status` attribute to an `status` attribute.
	// acceleration_status = "Enabled" => status = "Enabled"
	status := attr.ValueAsTokens()
	newResource.SetAttributeRaw("status", status)

	resource.RemoveAttribute(oldAttribute)

	return inFile, nil
}
