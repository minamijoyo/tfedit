package tfwrite

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// Output represents a output block.
// It implements the Block interface.
type Output struct {
	*block
}

var _ Block = (*Output)(nil)

// NewOutput creates a new instance of Output.
func NewOutput(block *hclwrite.Block) *Output {
	b := newBlock(block)
	return &Output{block: b}
}

// NewEmptyOutput creates a new Output with an empty body.
func NewEmptyOutput(outputType string) *Output {
	block := hclwrite.NewBlock("output", []string{outputType})
	return NewOutput(block)
}
