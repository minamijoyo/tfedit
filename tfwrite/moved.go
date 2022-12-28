package tfwrite

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// Moved represents a moved block.
// It implements the Block interface.
type Moved struct {
	*block
}

var _ Block = (*Moved)(nil)

// NewMoved creates a new instance of Moved.
func NewMoved(block *hclwrite.Block) *Moved {
	b := newBlock(block)
	return &Moved{block: b}
}

// NewEmptyMoved creates a new Moved with an empty body.
func NewEmptyMoved() *Moved {
	block := hclwrite.NewBlock("moved", []string{})
	return NewMoved(block)
}
