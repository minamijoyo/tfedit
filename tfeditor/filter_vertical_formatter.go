package tfeditor

import (
	"github.com/minamijoyo/tfedit/tfwrite"
)

// VerticalFormatterFilter is a filter implementation to format HCL in vertical.
// At time of writing, the default hcl formatter does not support vertical
// formatting. However, it's useful in some cases such as removing a block
// because leading and trailing newline tokens don't belong to a block, so
// deleting a block leaves extra newline tokens.
func VerticalFormatterFilter(inFile *tfwrite.File, block tfwrite.Block) (*tfwrite.File, error) {
	block.VerticalFormat()

	return inFile, nil
}
