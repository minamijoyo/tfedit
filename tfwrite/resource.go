package tfwrite

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// Resource represents a resource block.
// It implements the Block interface.
type Resource struct {
	*block
}

var _ Block = (*Resource)(nil)

// NewResource creates a new instance of Resource.
func NewResource(block *hclwrite.Block) *Resource {
	b := newBlock(block)
	return &Resource{block: b}
}

// NewEmptyResource creates a new Resource with an empty body.
func NewEmptyResource(resourceType string, resourceName string) *Resource {
	block := hclwrite.NewBlock("resource", []string{resourceType, resourceName})
	return NewResource(block)
}
