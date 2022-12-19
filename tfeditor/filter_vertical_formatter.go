package tfeditor

import (
	"github.com/minamijoyo/tfedit/tfwrite"
)

// VerticalFormatterBlockFilter is a BlockFilter implementation to format HCL
// in vertical.  At time of writing, the default hcl formatter does not support
// vertical formatting. However, it's useful in some cases such as removing a
// block because leading and trailing newline tokens don't belong to a block,
// so deleting a block leaves extra newline tokens.
type VerticalFormatterBlockFilter struct {
	// A block type to be formatted.
	blockType string
	// A schema type to be formatted.
	schemaType string
}

// NewVerticalFormatterBlockFilter returns a new instance of VerticalFormatterBlockFilter.
func NewVerticalFormatterBlockFilter(blockType string, schemaType string) *VerticalFormatterBlockFilter {
	return &VerticalFormatterBlockFilter{
		blockType:  blockType,
		schemaType: schemaType,
	}
}

// BlockFilter reads Terraform configuration and rewrite a given block,
// and writes Terraform configuration.
// When block type or schema type is set, format a block only if matching.
func (f *VerticalFormatterBlockFilter) BlockFilter(inFile *tfwrite.File, block tfwrite.Block) (*tfwrite.File, error) {

	if f.blockType != "" && f.blockType != block.Type() {
		return inFile, nil
	}
	if f.schemaType != "" && f.schemaType != block.SchemaType() {
		return inFile, nil
	}

	block.VerticalFormat()
	return inFile, nil
}
