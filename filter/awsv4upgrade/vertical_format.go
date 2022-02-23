package awsv4upgrade

import (
	"github.com/minamijoyo/tfedit/tfeditor"
	"github.com/minamijoyo/tfedit/tfwrite"
)

// verticalFormatterFilter is a Filter implementation to format HCL in vertical.
// At time of writing, the default hcl formatter does not support vertical
// formatting. However, it's useful in some cases such as removing a block
// because leading and trailing newline tokens don't belong to a block, so
// deleting a block leaves extra newline tokens.
type verticalFormatterFilter struct{}

var _ tfeditor.ResourceFilter = (*verticalFormatterFilter)(nil)

// ResourceFilter reads HCL and writes formatted contents in vertical.
func (f *verticalFormatterFilter) ResourceFilter(inFile *tfwrite.File, resource *tfwrite.Resource) (*tfwrite.File, error) {
	resource.VerticalFormat()

	return inFile, nil
}
