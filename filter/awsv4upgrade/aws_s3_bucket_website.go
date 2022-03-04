package awsv4upgrade

import (
	"github.com/minamijoyo/tfedit/tfeditor"
	"github.com/minamijoyo/tfedit/tfwrite"
)

// AWSS3BucketWebsiteFilter is a filter implementation for upgrading the
// website argument of aws_s3_bucket.
// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#website-website_domain-and-website_endpoint-arguments
type AWSS3BucketWebsiteFilter struct{}

var _ tfeditor.ResourceFilter = (*AWSS3BucketWebsiteFilter)(nil)

// NewAWSS3BucketWebsiteFilter creates a new instance of AWSS3BucketWebsiteFilter.
func NewAWSS3BucketWebsiteFilter() tfeditor.ResourceFilter {
	return &AWSS3BucketWebsiteFilter{}
}

// ResourceFilter upgrades the website argument of aws_s3_bucket.
func (f *AWSS3BucketWebsiteFilter) ResourceFilter(inFile *tfwrite.File, resource *tfwrite.Resource) (*tfwrite.File, error) {
	oldNestedBlock := "website"
	newResourceType := "aws_s3_bucket_website_configuration"

	nestedBlocks := resource.FindNestedBlocksByType(oldNestedBlock)
	if len(nestedBlocks) == 0 {
		return inFile, nil
	}

	resourceName := resource.Name()
	newResource := tfwrite.NewEmptyResource(newResourceType, resourceName)
	inFile.AppendResource(newResource)
	setBucketArgument(newResource, resource)

	websiteBlock := nestedBlocks[0]

	// Map an `index_document` attribute to an `index_document` block
	// index_document = "index.html"
	// =>
	// index_document {
	//   suffix = "index.html"
	// }
	indexDocumentAttr := websiteBlock.GetAttribute("index_document")
	if indexDocumentAttr != nil {
		indexDocumentBlock := tfwrite.NewEmptyNestedBlock("index_document")
		newResource.AppendNestedBlock(indexDocumentBlock)
		suffix := indexDocumentAttr.ValueAsTokens()
		indexDocumentBlock.SetAttributeRaw("suffix", suffix)
		websiteBlock.RemoveAttribute("index_document")
	}

	// Map an `error_document` attribute to an `error_document` block
	// error_document = "error.html"
	// =>
	// error_document {
	//   key = "error.html"
	// }
	errorDocumentAttr := websiteBlock.GetAttribute("error_document")
	if errorDocumentAttr != nil {
		errorDocumentBlock := tfwrite.NewEmptyNestedBlock("error_document")
		newResource.AppendNestedBlock(errorDocumentBlock)
		key := errorDocumentAttr.ValueAsTokens()
		errorDocumentBlock.SetAttributeRaw("key", key)
		websiteBlock.RemoveAttribute("error_document")
	}

	newResource.AppendUnwrappedNestedBlockBody(websiteBlock)
	resource.RemoveNestedBlock(websiteBlock)

	return inFile, nil
}
