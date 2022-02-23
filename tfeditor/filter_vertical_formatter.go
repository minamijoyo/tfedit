package tfeditor

import (
	"github.com/minamijoyo/tfedit/tfwrite"
)

// VerticalFormatterFilter is a filter implementation to format HCL in vertical.
// At time of writing, the default hcl formatter does not support vertical
// formatting. However, it's useful in some cases such as removing a block
// because leading and trailing newline tokens don't belong to a block, so
// deleting a block leaves extra newline tokens.
type VerticalFormatterFilter struct{}

var _ ResourceFilter = (*VerticalFormatterFilter)(nil)

// NewVerticalFormatterResourceFilter creates a new instance of VerticalFormatterFilter as ResourceFilter.
func NewVerticalFormatterResourceFilter() ResourceFilter {
	return &VerticalFormatterFilter{}
}

// ResourceFilter reads HCL and writes formatted contents in vertical.
func (f *VerticalFormatterFilter) ResourceFilter(inFile *tfwrite.File, resource *tfwrite.Resource) (*tfwrite.File, error) {
	resource.VerticalFormat()

	return inFile, nil
}
