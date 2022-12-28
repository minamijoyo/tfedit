package tfwrite

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// Locals represents a locals block.
// It implements the Block interface.
type Locals struct {
	*block
}

var _ Block = (*Locals)(nil)

// NewLocals creates a new instance of Locals.
func NewLocals(block *hclwrite.Block) *Locals {
	b := newBlock(block)
	return &Locals{block: b}
}

// NewEmptyLocals creates a new Locals with an empty body.
func NewEmptyLocals() *Locals {
	block := hclwrite.NewBlock("locals", []string{})
	return NewLocals(block)
}
