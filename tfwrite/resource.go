package tfwrite

import (
	"github.com/hashicorp/hcl/v2"
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

// SchemaType returns a type of resource.
// It returns the first label of block.
// Note that it's not the same as the *hclwrite.Block.Type().
func (r *Resource) SchemaType() string {
	labels := r.block.raw.Labels()
	return labels[0]
}

// Name returns a name of resource.
func (r *Resource) Name() string {
	labels := r.block.raw.Labels()
	return labels[1]
}

// SetAttributeByReference sets an attribute for a given name to a reference of
// another resource.
func (r *Resource) SetAttributeByReference(name string, refResource *Resource, refAttribute string) {
	traversal := hcl.Traversal{
		hcl.TraverseRoot{Name: refResource.SchemaType()},
		hcl.TraverseAttr{Name: refResource.Name()},
		hcl.TraverseAttr{Name: refAttribute},
	}
	r.block.raw.Body().SetAttributeTraversal(name, traversal)
}
