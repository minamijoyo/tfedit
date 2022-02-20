package tfwrite

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/minamijoyo/hcledit/editor"
	"github.com/zclconf/go-cty/cty"
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

// ValueAsString returns a value of Attribute as string.
func (a *Attribute) ValueAsString() (string, error) {
	return editor.GetAttributeValueAsString(a.raw)
}

// ValueAsTokens returns a value of Attribute as raw tokens.
func (a *Attribute) ValueAsTokens() hclwrite.Tokens {
	return a.raw.Expr().BuildTokens(nil)
}

// NestedBlock is a nested block of resource.
type NestedBlock struct {
	raw *hclwrite.Block
}

// NewNestedBlock creates a new instance of NestedBlock.
func NewNestedBlock(block *hclwrite.Block) *NestedBlock {
	return &NestedBlock{raw: block}
}

// NewEmptyNestedBlock creates a new NestedBlock with an empty body.
func NewEmptyNestedBlock(blockType string) *NestedBlock {
	block := hclwrite.NewBlock(blockType, []string{})
	return NewNestedBlock(block)
}

// SetType updates the type name of the block to a given name.
func (b *NestedBlock) SetType(typeName string) {
	b.raw.SetType(typeName)
}

// GetAttribute returns an attribute for a given name.
func (b *NestedBlock) GetAttribute(name string) *Attribute {
	attr := b.raw.Body().GetAttribute(name)
	if attr == nil {
		return nil
	}
	return NewAttribute(attr)
}

// SetAttributeValue sets an attribute for a given name with value.
func (b *NestedBlock) SetAttributeValue(name string, value cty.Value) {
	b.raw.Body().SetAttributeValue(name, value)
}

// SetAttributeRaw sets an attribute for a given name with raw tokens.
func (b *NestedBlock) SetAttributeRaw(name string, tokens hclwrite.Tokens) {
	b.raw.Body().SetAttributeRaw(name, tokens)
}

// RemoveAttribute removes an attribute for a given name from the resource.
func (b *NestedBlock) RemoveAttribute(name string) {
	b.raw.Body().RemoveAttribute(name)
}

// FindNestedBlocksByType returns all matching nested blocks from the body that have the
// given nested block type or returns an empty list if not found.
func (b *NestedBlock) FindNestedBlocksByType(blockType string) []*NestedBlock {
	var matched []*NestedBlock

	for _, block := range b.raw.Body().Blocks() {
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
func (b *NestedBlock) AppendNestedBlock(nestedBlock *NestedBlock) {
	body := b.raw.Body()
	body.AppendNewline()
	body.AppendBlock(nestedBlock.raw)
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
