package awsv4upgrade

import (
	"regexp"

	"github.com/minamijoyo/tfedit/tfwrite"
)

// AWSS3BucketWebsiteBlockFilter is a filter implementation for upgrading the
// website argument of aws_s3_bucket and renames all references for the
// website_domain and website_endpoint attributes to the new splitted resource.
// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#website-website_domain-and-website_endpoint-arguments
func AWSS3BucketWebsiteBlockFilter(inFile *tfwrite.File, block tfwrite.Block) (*tfwrite.File, error) {
	if block.Type() == "resource" && block.SchemaType() == "aws_s3_bucket" {
		return awsS3BucketWebsiteResourceFilter(inFile, block.(*tfwrite.Resource))
	}

	return awsS3BucketWebsiteReferenceFilter(inFile, block)
}

// awsS3BucketWebsiteResourceFilter splits the website argument to a new
// aws_s3_bucket_website_configuration resource.
func awsS3BucketWebsiteResourceFilter(inFile *tfwrite.File, resource *tfwrite.Resource) (*tfwrite.File, error) {
	if resource.SchemaType() != "aws_s3_bucket" {
		return inFile, nil
	}

	oldNestedBlock := "website"
	newResourceType := "aws_s3_bucket_website_configuration"

	nestedBlocks := resource.FindNestedBlocksByType(oldNestedBlock)
	if len(nestedBlocks) == 0 {
		return inFile, nil
	}

	resourceName := resource.Name()
	newResource := tfwrite.NewEmptyResource(newResourceType, resourceName)
	inFile.AppendBlock(newResource)
	setParentBucket(newResource, resource)

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

var regexWebsiteDomain = regexp.MustCompile(`aws_s3_bucket\.(.+)\.website_domain`)
var regexWebsiteEndpoint = regexp.MustCompile(`aws_s3_bucket\.(.+)\.website_endpoint`)

// awsS3BucketWebsiteReferenceFilter renames all references for the website_domain and website_endpoint attributes to the new splitted resource.
// aws_s3_bucket.example.website_domain => aws_s3_bucket_website_configuration.example.website_domain
// aws_s3_bucket.example.website_endpoint => aws_s3_bucket_website_configuration.example.website_endpoint
func awsS3BucketWebsiteReferenceFilter(inFile *tfwrite.File, block tfwrite.Block) (*tfwrite.File, error) {
	refs := block.References()
	for _, ref := range refs {
		from := regexWebsiteDomain.FindString(ref)
		if from != "" {
			to := regexWebsiteDomain.ReplaceAllString(from, `aws_s3_bucket_website_configuration.$1.website_domain`)
			block.RenameReference(from, to)
		}
		from = regexWebsiteEndpoint.FindString(ref)
		if from != "" {
			to := regexWebsiteEndpoint.ReplaceAllString(from, `aws_s3_bucket_website_configuration.$1.website_endpoint`)
			block.RenameReference(from, to)
		}
	}

	return inFile, nil
}
