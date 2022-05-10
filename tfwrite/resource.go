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
// It returns the second label of block.
func (r *Resource) Name() string {
	labels := r.block.raw.Labels()
	return labels[1]
}

// Count returns a meta argument of count.
// It returns nil if not found.
func (r *Resource) Count() *Attribute {
	return r.GetAttribute("count")
}

// ForEach returns a meta argument of for_each.
// It returns nil if not found.
func (r *Resource) ForEach() *Attribute {
	return r.GetAttribute("for_each")
}

// ReferableName returns a name of resource instance which can be referenced as
// a part of address.
// It contains an index reference if count or for_each is set.
// If neither count nor for_each is set, it just returns the name.
func (r *Resource) ReferableName() string {
	name := r.Name()

	if count := r.Count(); count != nil {
		return name + "[count.index]"
	}

	if forEach := r.ForEach(); forEach != nil {
		return name + "[each.key]"
	}

	return name
}

// SetAttributeByReference sets an attribute for a given name to a reference of
// another resource.
func (r *Resource) SetAttributeByReference(name string, refResource *Resource, refAttribute string) {
	traversal := hcl.Traversal{
		hcl.TraverseRoot{Name: refResource.SchemaType()},
		hcl.TraverseAttr{Name: refResource.ReferableName()},
		hcl.TraverseAttr{Name: refAttribute},
	}
	r.block.raw.Body().SetAttributeTraversal(name, traversal)
}
