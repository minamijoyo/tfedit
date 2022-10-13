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

// findBlocksByType returns all matching blocks from the body that have the
// given blockType or returns an empty list if not found.
func (f *File) findBlocksByType(blockType string) []*block {
	var blocks []*block

	for _, block := range f.Raw().Body().Blocks() {
		if block.Type() == blockType {
			blocks = append(blocks, newBlock(block))
		}
	}

	return blocks
}

// FindResourcesByType returns all matching resources from the body that have the
// given resourceType or returns an empty list if not found.
func (f *File) FindResourcesByType(schemaType string) []*Resource {
	var matched []*Resource

	for _, block := range f.findBlocksByType("resource") {
		b := NewResource(block.Raw())
		if b.SchemaType() != schemaType {
			continue
		}

		matched = append(matched, b)
	}

	return matched
}

// FindDataSourcesByType returns all matching data sources from the body that have the
// given dataSourceType or returns an empty list if not found.
func (f *File) FindDataSourcesByType(schemaType string) []*DataSource {
	var matched []*DataSource

	for _, block := range f.findBlocksByType("data") {
		b := NewDataSource(block.Raw())
		if b.SchemaType() != schemaType {
			continue
		}

		matched = append(matched, b)
	}

	return matched
}

// FindProvidersByType returns all matching providers from the body that have the
// given providerType or returns an empty list if not found.
func (f *File) FindProvidersByType(schemaType string) []*Provider {
	var matched []*Provider

	for _, block := range f.findBlocksByType("provider") {
		b := NewProvider(block.Raw())
		if b.SchemaType() != schemaType {
			continue
		}

		matched = append(matched, b)
	}

	return matched
}

// appendBlock appends a given block to the file.
func (f *File) appendBlock(block Block) {
	body := f.raw.Body()
	body.AppendNewline()
	body.AppendBlock(block.Raw())
}

// AppendResource appends a given resource to the file.
func (f *File) AppendResource(block *Resource) {
	f.appendBlock(block)
}

// AppendDataSource appends a given data source to the file.
func (f *File) AppendDataSource(block *DataSource) {
	f.appendBlock(block)
}

// AppendProvider appends a given provider to the file.
func (f *File) AppendProvider(block *Provider) {
	f.appendBlock(block)
}
