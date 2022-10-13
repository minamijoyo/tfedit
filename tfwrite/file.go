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

// FindResourcesByType returns all matching resources from the body that have the
// given resourceType or returns an empty list if not found.
func (f *File) FindResourcesByType(resourceType string) []*Resource {
	var matched []*Resource

	for _, block := range f.raw.Body().Blocks() {
		if block.Type() != "resource" {
			continue
		}

		labels := block.Labels()
		if len(labels) == 2 && labels[0] != resourceType {
			continue
		}

		resource := NewResource(block)
		matched = append(matched, resource)
	}

	return matched
}

// FindDataSourcesByType returns all matching data sources from the body that have the
// given dataSourceType or returns an empty list if not found.
func (f *File) FindDataSourcesByType(dataSourceType string) []*DataSource {
	var matched []*DataSource

	for _, block := range f.raw.Body().Blocks() {
		if block.Type() != "data" {
			continue
		}

		labels := block.Labels()
		if len(labels) == 2 && labels[0] != dataSourceType {
			continue
		}

		dataSource := NewDataSource(block)
		matched = append(matched, dataSource)
	}

	return matched
}

// FindProvidersByType returns all matching providers from the body that have the
// given providerType or returns an empty list if not found.
func (f *File) FindProvidersByType(providerType string) []*Provider {
	var matched []*Provider

	for _, block := range f.raw.Body().Blocks() {
		if block.Type() != "provider" {
			continue
		}

		labels := block.Labels()
		if len(labels) == 1 && labels[0] != providerType {
			continue
		}

		provider := NewProvider(block)
		matched = append(matched, provider)
	}

	return matched
}

// AppendResource appends a given resource to the file.
func (f *File) AppendResource(resource *Resource) {
	body := f.raw.Body()
	body.AppendNewline()
	body.AppendBlock(resource.raw)
}

// AppendDataSource appends a given data source to the file.
func (f *File) AppendDataSource(dataSource *DataSource) {
	body := f.raw.Body()
	body.AppendNewline()
	body.AppendBlock(dataSource.raw)
}

// AppendProvider appends a given provider to the file.
func (f *File) AppendProvider(provider *Provider) {
	body := f.raw.Body()
	body.AppendNewline()
	body.AppendBlock(provider.raw)
}
