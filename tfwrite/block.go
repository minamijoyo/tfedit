package tfwrite

import (
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

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

	// SchemaType returns the first label of block.
	// If it does not have a label, return an empty string.
	// Note that it's not the same as the *hclwrite.Block.Type().
	SchemaType() string

	// SetType updates the type name of the block to a given name.
	SetType(typeName string)

	// Attributes returns all attributes.
	// Note that this does not return attributes in nested blocks.
	Attributes() []*Attribute

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

	// CopyAttribute is a helper method which copies an attribute from a given block.
	// Do nothing if the given block doesn't have the attribute.
	CopyAttribute(from Block, name string)

	// NestedBlocks returns all nested blocks.
	NestedBlocks() []*NestedBlock

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

	// References returns all variable references for all attributes.
	// It returns a unique and sorted list.
	References() []string

	// RenameReference renames all variable references for all attributes.
	// The `from` and `to` arguments are specified as dot-delimited resource addresses.
	// The `from` argument can be a partial prefix match, but must match the length
	// of the `to` argument.
	RenameReference(from string, to string)
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

// SchemaType returns the first label of block.
// If it does not have a label, return an empty string.
// Note that it's not the same as the *hclwrite.Block.Type().
func (b *block) SchemaType() string {
	labels := b.raw.Labels()
	if len(labels) >= 1 {
		return labels[0]
	}
	return ""
}

// SetType updates the type name of the block to a given name.
func (b *block) SetType(typeName string) {
	b.raw.SetType(typeName)
}

// Attributes returns all attributes.
// Note that this does not return attributes in nested blocks.
func (b *block) Attributes() []*Attribute {
	var ret []*Attribute
	attrs := b.raw.Body().Attributes()
	for _, attr := range attrs {
		ret = append(ret, NewAttribute(attr))
	}
	return ret
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

// CopyAttribute is a helper method which copies an attribute from a given block.
// Do nothing if the given block doesn't have the attribute.
func (b *block) CopyAttribute(from Block, name string) {
	if attr := from.GetAttribute(name); attr != nil {
		b.AppendAttribute(attr)
	}
}

// NestedBlocks returns all nested blocks.
func (b *block) NestedBlocks() []*NestedBlock {
	var ret []*NestedBlock
	blocks := b.raw.Body().Blocks()
	for _, block := range blocks {
		ret = append(ret, NewNestedBlock(block))
	}
	return ret
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

// References returns all variable references for all attributes.
// It returns a unique and sorted list.
func (b *block) References() []string {
	refs := map[string]struct{}{}
	for _, key := range b.rawReferences() {
		// To remove duplicates, append keys to a map.
		refs[key] = struct{}{}
	}

	keys := maps.Keys(refs)
	slices.Sort(keys)
	return keys
}

// rawReferences returns a list of raw variable references.
// The result may be unsorted and duplicated,
// but it is efficient to merge the result later.
func (b *block) rawReferences() []string {
	refs := []string{}
	for _, attr := range b.Attributes() {
		ks := attr.rawReferences()
		refs = append(refs, ks...)
	}

	for _, block := range b.NestedBlocks() {
		// recursive call
		ks := block.rawReferences()
		refs = append(refs, ks...)
	}

	return refs
}

// RenameReference renames all variable references for all attributes.
// The `from` and `to` arguments are specified as dot-delimited resource addresses.
// The `from` argument can be a partial prefix match, but must match the length
// of the `to` argument.
func (b *block) RenameReference(from string, to string) {
	for _, attr := range b.Attributes() {
		attr.RenameReference(from, to)
	}

	for _, block := range b.NestedBlocks() {
		// recursive call
		block.RenameReference(from, to)
	}
}
