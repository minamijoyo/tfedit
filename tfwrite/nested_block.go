package tfwrite

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// NestedBlock represents a nested block.
// It implements the Block interface.
type NestedBlock struct {
	*block
}

var _ Block = (*NestedBlock)(nil)

// NewNestedBlock creates a new instance of NestedBlock.
func NewNestedBlock(block *hclwrite.Block) *NestedBlock {
	b := newBlock(block)
	return &NestedBlock{block: b}
}

// NewEmptyNestedBlock creates a new NestedBlock with an empty body.
func NewEmptyNestedBlock(blockType string) *NestedBlock {
	block := hclwrite.NewBlock(blockType, []string{})
	return NewNestedBlock(block)
}
