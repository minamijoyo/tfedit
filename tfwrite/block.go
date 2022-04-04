package tfwrite

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/minamijoyo/hcledit/editor"
	"github.com/zclconf/go-cty/cty"
)

// Block represents an abstract HCL block.
type Block interface {
	// Raw returns a raw block instance of hclwrite.
	// It should keep it private if possible, but it's required for some
	// operations.
	Raw() *hclwrite.Block

	// Type returns a type of block.
	Type() string

	// SetType updates the type name of the block to a given name.
	SetType(typeName string)

	// GetAttribute returns an attribute for a given name.
	GetAttribute(name string) *Attribute

	// SetAttributeValue sets an attribute for a given name with value.
	SetAttributeValue(name string, value cty.Value)

	// SetAttributeRaw sets an attribute for a given name with raw tokens.
	SetAttributeRaw(name string, tokens hclwrite.Tokens)

	// AppendAttribute appends a given attribute to the block.
	AppendAttribute(attr *Attribute)

	// RemoveAttribute removes an attribute for a given name from the block.
	RemoveAttribute(name string)

	// AppendNestedBlock appends a given nested block to the parent block.
	AppendNestedBlock(nestedBlock Block)

	// AppendUnwrappedNestedBlockBody appends a body of a given block to
	// the parent block. It looks a weird operation, but it's often needed for
	// refactoring like splitting a resource type for sub resource types.
	AppendUnwrappedNestedBlockBody(nestedBlock Block)

	// RemoveNestedBlock removes a given nested block from the parent block.
	RemoveNestedBlock(nestedBlock Block)

	// FindNestedBlocksByType returns all matching blocks from the body that have
	// the given block type or returns an empty list if not found.
	FindNestedBlocksByType(blockType string) []Block

	// VerticalFormat formats a body of the block in vertical. Since
	// VerticalFormat clears tokens internally, If you call VerticalFormat each
	// time RemoveNestedBlock is called, the subsequent RemoveNestedBlock will not
	// work properly, so call VerticalFormat only once for each block.
	VerticalFormat()
}

// block implements the Block interface.
// It implements shared logic for all block types.
// The abstract block does not exist in the Terraform language specification,
// so it's intentionally private.
type block struct {
	raw *hclwrite.Block
}

var _ Block = (*block)(nil)

// newBlock creates a new instance of block.
func newBlock(raw *hclwrite.Block) *block {
	return &block{raw: raw}
}

// newEmptyBlock creates a new block with an empty body.
func newEmptyBlock(blockType string) *block {
	block := hclwrite.NewBlock(blockType, []string{})
	return newBlock(block)
}

// Raw returns a raw block instance of hclwrite.
// It should keep it private if possible, but it's required for some
// operations.
func (b *block) Raw() *hclwrite.Block {
	return b.raw
}

// Type returns a type of block.
func (b *block) Type() string {
	return b.raw.Type()
}

// SetType updates the type name of the block to a given name.
func (b *block) SetType(typeName string) {
	b.raw.SetType(typeName)
}

// GetAttribute returns an attribute for a given name.
func (b *block) GetAttribute(name string) *Attribute {
	attr := b.raw.Body().GetAttribute(name)
	if attr == nil {
		return nil
	}
	return NewAttribute(attr)
}

// SetAttributeValue sets an attribute for a given name with value.
func (b *block) SetAttributeValue(name string, value cty.Value) {
	b.raw.Body().SetAttributeValue(name, value)
}

// SetAttributeRaw sets an attribute for a given name with raw tokens.
func (b *block) SetAttributeRaw(name string, tokens hclwrite.Tokens) {
	b.raw.Body().SetAttributeRaw(name, tokens)
}

// AppendAttribute appends a given attribute to the block.
func (b *block) AppendAttribute(attr *Attribute) {
	expr := attr.raw.BuildTokens(nil)
	b.raw.Body().AppendUnstructuredTokens(expr)
}

// RemoveAttribute removes an attribute for a given name from the block.
func (b *block) RemoveAttribute(name string) {
	b.raw.Body().RemoveAttribute(name)
}

// AppendNestedBlock appends a given nested block to the parent block.
func (b *block) AppendNestedBlock(nestedBlock Block) {
	body := b.raw.Body()
	body.AppendNewline()
	body.AppendBlock(nestedBlock.Raw())
}

// AppendUnwrappedNestedBlockBody appends a body of a given block to
// the parent block. It looks a weird operation, but it's often needed for
// refactoring like splitting a resource type for sub resource types.
func (b *block) AppendUnwrappedNestedBlockBody(nestedBlock Block) {
	unwrapped := nestedBlock.Raw().Body().BuildTokens(nil)
	b.raw.Body().AppendUnstructuredTokens(unwrapped)
}

// RemoveNestedBlock removes a given nested block from the parent block.
func (b *block) RemoveNestedBlock(nestedBlock Block) {
	b.raw.Body().RemoveBlock(nestedBlock.Raw())
}

// FindNestedBlocksByType returns all matching blocks from the body that have
// the given block type or returns an empty list if not found.
func (b *block) FindNestedBlocksByType(blockType string) []Block {
	var matched []Block

	for _, block := range b.raw.Body().Blocks() {
		if block.Type() != blockType {
			continue
		}

		labels := block.Labels()
		if len(labels) == 2 && labels[0] != blockType {
			continue
		}

		newBlock := newBlock(block)
		matched = append(matched, newBlock)
	}

	return matched
}

// VerticalFormat formats a body of the block in vertical.  Since
// VerticalFormat clears tokens internally, If you call VerticalFormat each
// time RemoveNestedBlock is called, the subsequent RemoveNestedBlock will not
// work properly, so call VerticalFormat only once for each block.
func (b *block) VerticalFormat() {
	body := b.raw.Body()
	unformatted := body.BuildTokens(nil)
	formatted := editor.VerticalFormat(unformatted)
	body.Clear()
	body.AppendNewline()
	body.AppendUnstructuredTokens(formatted)
}
