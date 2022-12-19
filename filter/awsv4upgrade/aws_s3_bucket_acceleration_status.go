package awsv4upgrade

import (
	"github.com/minamijoyo/tfedit/tfwrite"
)

// AWSS3BucketAccelerationStatusResourceFilter is a filter
// implementation for upgrading the acceleration_status argument of aws_s3_bucket.
// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#acceleration_status-argument
func AWSS3BucketAccelerationStatusResourceFilter(inFile *tfwrite.File, resource *tfwrite.Resource) (*tfwrite.File, error) {
	if resource.SchemaType() != "aws_s3_bucket" {
		return inFile, nil
	}

	oldAttribute := "acceleration_status"
	newResourceType := "aws_s3_bucket_accelerate_configuration"

	attr := resource.GetAttribute(oldAttribute)
	if attr == nil {
		return inFile, nil
	}

	resourceName := resource.Name()
	newResource := tfwrite.NewEmptyResource(newResourceType, resourceName)
	inFile.AppendBlock(newResource)
	setParentBucket(newResource, resource)

	// Map an `acceleration_status` attribute to an `status` attribute.
	// acceleration_status = "Enabled" => status = "Enabled"
	status := attr.ValueAsTokens()
	newResource.SetAttributeRaw("status", status)

	resource.RemoveAttribute(oldAttribute)

	return inFile, nil
}
