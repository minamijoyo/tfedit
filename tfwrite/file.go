package tfwrite

import "github.com/hashicorp/hcl/v2/hclwrite"

// File represents a Terraform configuration file
type File struct {
	raw *hclwrite.File
}

// NewFile creates a new instance of File.
func NewFile(file *hclwrite.File) *File {
	return &File{raw: file}
}

// NewEmptyFile creates a new file with an empty body.
func NewEmptyFile() *File {
	file := hclwrite.NewEmptyFile()
	return NewFile(file)
}

// Raw returns a raw object for hclwrite.
func (f *File) Raw() *hclwrite.File {
	return f.raw
}

// parseBlock is a factory method for Block.
func parseBlock(block *hclwrite.Block) Block {
	switch block.Type() {
	case "resource":
		return NewResource(block)
	case "data":
		return NewDataSource(block)
	case "provider":
		return NewProvider(block)
	default:
		return newBlock(block) // unknown
	}
}

// Blocks returns all blocks.
func (f *File) Blocks() []Block {
	var blocks []Block

	for _, block := range f.Raw().Body().Blocks() {
		b := parseBlock(block)
		blocks = append(blocks, b)
	}

	return blocks
}

// FindBlocksByType returns all matching blocks from the body that have the
// given blockType and schemaType or returns an empty list if not found.
// If the given blockType or schemaType are a non-empty string,
// filter the results.
func (f *File) FindBlocksByType(blockType string, schemaType string) []Block {
	var matched []Block

	for _, b := range f.Blocks() {
		if blockType != "" && blockType != b.Type() {
			continue
		}
		if schemaType != "" && schemaType != b.SchemaType() {
			continue
		}
		matched = append(matched, b)
	}

	return matched
}

// findAnyBlockByType is a generic helper method for implementing block finders.
// It returns all matching blocks from the body that have the given type
// parameter and schemaType or returns an empty list if not found.
// If the given schemaType is a non-empty string, filter the results.
// Note that the Generics doesn't allow an instance method to have a type parameter.
func findAnyBlockByType[T Block](f *File, schemaType string) []T {
	var matched []T

	for _, block := range f.Blocks() {
		b, ok := block.(T)
		if !ok {
			// ignore other block types
			continue
		}

		if schemaType != "" && schemaType != b.SchemaType() {
			continue
		}
		matched = append(matched, b)
	}

	return matched
}

// FindResourcesByType returns all matching resources from the body that have the
// given schemaType or returns an empty list if not found.
func (f *File) FindResourcesByType(schemaType string) []*Resource {
	return findAnyBlockByType[*Resource](f, schemaType)
}

// FindDataSourcesByType returns all matching data sources from the body that have the
// given schemaType or returns an empty list if not found.
func (f *File) FindDataSourcesByType(schemaType string) []*DataSource {
	return findAnyBlockByType[*DataSource](f, schemaType)
}

// FindProvidersByType returns all matching providers from the body that have the
// given schemaType or returns an empty list if not found.
func (f *File) FindProvidersByType(schemaType string) []*Provider {
	return findAnyBlockByType[*Provider](f, schemaType)
}

// AppendBlock appends a given block to the file.
func (f *File) AppendBlock(block Block) {
	body := f.raw.Body()
	body.AppendNewline()
	body.AppendBlock(block.Raw())
}
