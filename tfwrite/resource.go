package tfwrite

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// Resource represents a resource block.
type Resource struct {
	raw *hclwrite.Block
}

// NewResource creates a new instance of Resource.
func NewResource(block *hclwrite.Block) *Resource {
	return &Resource{raw: block}
}

// NewEmptyResource creates a new Resource with an empty body.
func NewEmptyResource(resourceType string, resourceName string) *Resource {
	block := hclwrite.NewBlock("resource", []string{resourceType, resourceName})
	return NewResource(block)
}

// Attribute is an attribute of resource.
type Attribute struct {
	raw *hclwrite.Attribute
}

// NewAttribute creates a new instance of Attribute.
func NewAttribute(attr *hclwrite.Attribute) *Attribute {
	return &Attribute{raw: attr}
}

// NestedBlock is a nested block of resource.
type NestedBlock struct {
	raw *hclwrite.Block
}

// NewNestedBlock creates a new instance of NestedBlock.
func NewNestedBlock(block *hclwrite.Block) *NestedBlock {
	return &NestedBlock{raw: block}
}

// Type returns a resource type.
func (r *Resource) Type() string {
	labels := r.raw.Labels()
	return labels[0]
}

// Name returns a resource name.
func (r *Resource) Name() string {
	labels := r.raw.Labels()
	return labels[1]
}

// GetAttribute returns an attribute for a given name.
func (r *Resource) GetAttribute(name string) *Attribute {
	attr := r.raw.Body().GetAttribute(name)
	if attr == nil {
		return nil
	}
	return NewAttribute(attr)
}

// SetAttributeByReference sets an attribute for a given name to a reference of
// another resource.
func (r *Resource) SetAttributeByReference(name string, refResource *Resource, refAttribute string) {
	traversal := hcl.Traversal{
		hcl.TraverseRoot{Name: refResource.Type()},
		hcl.TraverseAttr{Name: refResource.Name()},
		hcl.TraverseAttr{Name: refAttribute},
	}
	r.raw.Body().SetAttributeTraversal(name, traversal)
}

// AppendAttribute appends a given attribute to the resource.
func (r *Resource) AppendAttribute(attr *Attribute) {
	expr := attr.raw.BuildTokens(nil)
	r.raw.Body().AppendUnstructuredTokens(expr)
}

// RemoveAttribute removes an attribute for a given name from the resource.
func (r *Resource) RemoveAttribute(name string) {
	r.raw.Body().RemoveAttribute(name)
}

// FindNestedBlocksByType returns all matching nested blocks from the body that have the
// given nested block type or returns an empty list if not found.
func (r *Resource) FindNestedBlocksByType(blockType string) []*NestedBlock {
	var matched []*NestedBlock

	for _, block := range r.raw.Body().Blocks() {
		if block.Type() != blockType {
			continue
		}

		labels := block.Labels()
		if len(labels) == 2 && labels[0] != blockType {
			continue
		}

		resource := NewNestedBlock(block)
		matched = append(matched, resource)
	}

	return matched
}

// AppendUnwrappedNestedBlockBody appends a body of a given nested block to the
// resource.
func (r *Resource) AppendUnwrappedNestedBlockBody(nestedBlock *NestedBlock) {
	unwrapped := nestedBlock.raw.Body().BuildTokens(nil)
	r.raw.Body().AppendUnstructuredTokens(unwrapped)
}

// RemoveNestedBlock removes a given nested block from the resource.
func (r *Resource) RemoveNestedBlock(nestedBlock *NestedBlock) {
	r.raw.Body().RemoveBlock(nestedBlock.raw)
}
