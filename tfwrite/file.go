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

// AppendResource is a helper method for appending a given resource block to
// the file.
func (f *File) AppendResource(resource *Resource) {
	body := f.raw.Body()
	body.AppendNewline()
	body.AppendBlock(resource.raw)
}