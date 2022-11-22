package awsv4upgrade

import (
	"github.com/minamijoyo/tfedit/tfwrite"
)

// AWSS3BucketRequestPayerResourceFilter is a filter implementation for
// upgrading the request_payer argument of aws_s3_bucket.
// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#request_payer-argument
func AWSS3BucketRequestPayerResourceFilter(inFile *tfwrite.File, resource *tfwrite.Resource) (*tfwrite.File, error) {
	oldAttribute := "request_payer"
	newResourceType := "aws_s3_bucket_request_payment_configuration"

	attr := resource.GetAttribute(oldAttribute)
	if attr == nil {
		return inFile, nil
	}

	resourceName := resource.Name()
	newResource := tfwrite.NewEmptyResource(newResourceType, resourceName)
	inFile.AppendResource(newResource)
	setParentBucket(newResource, resource)

	// Map an `request_payer` attribute to an `payer` attribute.
	// request_payer = "Requester" => payer = "Requester"
	payer := attr.ValueAsTokens()
	newResource.SetAttributeRaw("payer", payer)

	resource.RemoveAttribute(oldAttribute)

	return inFile, nil
}
