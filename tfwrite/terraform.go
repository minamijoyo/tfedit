package tfwrite

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// Terraform represents a terraform block.
// It implements the Block interface.
type Terraform struct {
	*block
}

var _ Block = (*Terraform)(nil)

// NewTerraform creates a new instance of Terraform.
func NewTerraform(block *hclwrite.Block) *Terraform {
	b := newBlock(block)
	return &Terraform{block: b}
}

// NewEmptyTerraform creates a new Terraform with an empty body.
func NewEmptyTerraform() *Terraform {
	block := hclwrite.NewBlock("terraform", []string{})
	return NewTerraform(block)
}
