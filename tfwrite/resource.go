package tfwrite

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/minamijoyo/hcledit/editor"
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

// AppendNestedBlock appends a given nested block to the resource.
func (r *Resource) AppendNestedBlock(nestedBlock *NestedBlock) {
	body := r.raw.Body()
	body.AppendNewline()
	body.AppendBlock(nestedBlock.raw)
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

// VerticalFormat formats a body of the resource in vertical.
// Since VerticalFormat clears tokens internally, If you call VerticalFormat
// each time RemoveNestedBlock is called, the subsequent RemoveNestedBlock will
// not work properly, so call VerticalFormat only once for each resource.
func (r *Resource) VerticalFormat() {
	body := r.raw.Body()
	unformatted := body.BuildTokens(nil)
	formatted := editor.VerticalFormat(unformatted)
	body.Clear()
	body.AppendNewline()
	body.AppendUnstructuredTokens(formatted)
}
