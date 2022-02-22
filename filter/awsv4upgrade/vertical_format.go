package awsv4upgrade

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/minamijoyo/hcledit/editor"
	"github.com/minamijoyo/tfedit/tfwrite"
)

// verticalFormatterFilter is a Filter implementation to format HCL in vertical.
// At time of writing, the default hcl formatter does not support vertical
// formatting. However, it's useful in some cases such as removing a block
// because leading and trailing newline tokens don't belong to a block, so
// deleting a block leaves extra newline tokens.
type verticalFormatterFilter struct{}

var _ editor.Filter = (*verticalFormatterFilter)(nil)

// Filter reads HCL and writes formatted contents in vertical.
func (f *verticalFormatterFilter) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	file := tfwrite.NewFile(inFile)
	oldResourceType := "aws_s3_bucket"

	targets := file.FindResourcesByType(oldResourceType)
	for _, oldResource := range targets {
		oldResource.VerticalFormat()
	}

	return inFile, nil
}
