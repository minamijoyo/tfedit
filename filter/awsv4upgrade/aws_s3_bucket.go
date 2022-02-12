package awsv4upgrade

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/minamijoyo/hcledit/editor"
	"github.com/zclconf/go-cty/cty"
)

// AWSS3BucketFilter is a filter implementation for upgrading configurations of
// aws_s3_bucket to AWS provider v4.
// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#s3-bucket-refactor
type AWSS3BucketFilter struct {
}

var _ editor.Filter = (*AWSS3BucketFilter)(nil)

// NewAWSS3BucketFilter creates a new instance of AWSS3BucketFilter.
func NewAWSS3BucketFilter() editor.Filter {
	return &AWSS3BucketFilter{}
}

// Filter upgrades configurations of aws_s3_bucket to AWS provider v4.
func (f *AWSS3BucketFilter) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	blocks := inFile.Body().Blocks()
	for _, b := range blocks {
		if b.Type() != "resource" {
			continue
		}

		labels := b.Labels()
		if len(labels) == 2 && labels[0] != "aws_s3_bucket" {
			continue
		}
		bName := labels[1]

		if b.Body().GetAttribute("acl") != nil {
			inFile.Body().AppendNewline()
			newblock := inFile.Body().AppendNewBlock("resource", []string{"aws_s3_bucket_acl", bName})
			newblock.Body().SetAttributeTraversal("bucket", hcl.Traversal{
				hcl.TraverseRoot{Name: "aws_s3_bucket"},
				hcl.TraverseAttr{Name: bName},
				hcl.TraverseAttr{Name: "id"},
			})
			newblock.Body().SetAttributeValue("acl", cty.StringVal("private"))
			b.Body().RemoveAttribute("acl")
		}

		if nested := b.Body().FirstMatchingBlock("logging", []string{}); nested != nil {
			inFile.Body().AppendNewline()
			newblock := inFile.Body().AppendNewBlock("resource", []string{"aws_s3_bucket_logging", bName})
			newblock.Body().SetAttributeTraversal("bucket", hcl.Traversal{
				hcl.TraverseRoot{Name: "aws_s3_bucket"},
				hcl.TraverseAttr{Name: bName},
				hcl.TraverseAttr{Name: "id"},
			})
			newblock.Body().AppendUnstructuredTokens(nested.Body().BuildTokens(nil))
			b.Body().RemoveBlock(nested)
		}
	}

	return inFile, nil
}
