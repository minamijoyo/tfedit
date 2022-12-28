package tfwrite

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// Module represents a module block.
// It implements the Block interface.
type Module struct {
	*block
}

var _ Block = (*Module)(nil)

// NewModule creates a new instance of Module.
func NewModule(block *hclwrite.Block) *Module {
	b := newBlock(block)
	return &Module{block: b}
}

// NewEmptyModule creates a new Module with an empty body.
func NewEmptyModule(moduleType string) *Module {
	block := hclwrite.NewBlock("module", []string{moduleType})
	return NewModule(block)
}
