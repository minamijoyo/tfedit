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

// FindBlocksByType returns all matching blocks from the body that have the
// given blockType and schemaType or returns an empty list if not found.
// If the given schemaType is a non-empty string, filter the results.
func (f *File) FindBlocksByType(blockType string, schemaType string) []Block {
	var blocks []Block

	for _, block := range f.Raw().Body().Blocks() {
		if block.Type() == blockType {
			var b Block
			switch blockType {
			case "resource":
				b = NewResource(block)
			case "data":
				b = NewDataSource(block)
			case "provider":
				b = NewProvider(block)
			default:
				continue
			}

			if schemaType != "" && schemaType != b.SchemaType() {
				continue
			}
			blocks = append(blocks, b)
		}
	}

	return blocks
}

// FindResourcesByType returns all matching resources from the body that have the
// given schemaType or returns an empty list if not found.
func (f *File) FindResourcesByType(schemaType string) []*Resource {
	var matched []*Resource

	for _, block := range f.FindBlocksByType("resource", schemaType) {
		b := block.(*Resource)
		matched = append(matched, b)
	}

	return matched
}

// FindDataSourcesByType returns all matching data sources from the body that have the
// given schemaType or returns an empty list if not found.
func (f *File) FindDataSourcesByType(schemaType string) []*DataSource {
	var matched []*DataSource

	for _, block := range f.FindBlocksByType("data", schemaType) {
		b := block.(*DataSource)
		matched = append(matched, b)
	}

	return matched
}

// FindProvidersByType returns all matching providers from the body that have the
// given schemaType or returns an empty list if not found.
func (f *File) FindProvidersByType(schemaType string) []*Provider {
	var matched []*Provider

	for _, block := range f.FindBlocksByType("provider", schemaType) {
		b := block.(*Provider)
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
