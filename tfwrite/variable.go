package tfwrite

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// Variable represents a variable block.
// It implements the Block interface.
type Variable struct {
	*block
}

var _ Block = (*Variable)(nil)

// NewVariable creates a new instance of Variable.
func NewVariable(block *hclwrite.Block) *Variable {
	b := newBlock(block)
	return &Variable{block: b}
}

// NewEmptyVariable creates a new Variable with an empty body.
func NewEmptyVariable(variableType string) *Variable {
	block := hclwrite.NewBlock("variable", []string{variableType})
	return NewVariable(block)
}
