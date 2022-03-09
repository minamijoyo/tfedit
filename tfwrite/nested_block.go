package tfwrite

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

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

// AppendUnwrappedNestedBlockBody appends a body of a given nested block to
// another nestedBlock.
func (b *NestedBlock) AppendUnwrappedNestedBlockBody(nestedBlock *NestedBlock) {
	unwrapped := nestedBlock.raw.Body().BuildTokens(nil)
	b.raw.Body().AppendUnstructuredTokens(unwrapped)
}
