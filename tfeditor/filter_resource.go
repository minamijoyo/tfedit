package tfeditor

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/minamijoyo/hcledit/editor"
	"github.com/minamijoyo/tfedit/tfwrite"
)

// ResourceFilter is an interface which reads Terraform configuration and
// rewrite a given resource, and writes Terraform configuration.
type ResourceFilter interface {
	// ResourceFilter reads Terraform configuration and rewrite a given resource,
	// and writes Terraform configuration.
	ResourceFilter(*tfwrite.File, *tfwrite.Resource) (*tfwrite.File, error)
}

// MultiResourceFilter is a ResourceFilter implementation which applies
// multiple resource filters to a given resource in sequence.
type MultiResourceFilter struct {
	filters []ResourceFilter
}

var _ ResourceFilter = (*MultiResourceFilter)(nil)

// NewMultiResourceFilter creates a new instance of MultiResourceFilter.
func NewMultiResourceFilter(filters []ResourceFilter) ResourceFilter {
	return &MultiResourceFilter{
		filters: filters,
	}
}

// ResourceFilter applies multiple filters to a given resource in sequence.
func (f *MultiResourceFilter) ResourceFilter(inFile *tfwrite.File, resource *tfwrite.Resource) (*tfwrite.File, error) {
	current := inFile
	for _, f := range f.filters {
		next, err := f.ResourceFilter(current, resource)
		if err != nil {
			return nil, err
		}
		current = next
	}
	return current, nil
}

// ResourcesByTypeFilter is a Filter implementation for applying a filter to
// multiple resources with a given resource type.
type ResourcesByTypeFilter struct {
	resourceType string
	filter       ResourceFilter
}

var _ editor.Filter = (*ResourcesByTypeFilter)(nil)

// NewResourcesByTypeFilter creates a new instance of ResourcesByTypeFilter.
func NewResourcesByTypeFilter(resourceType string, filter ResourceFilter) editor.Filter {
	return &ResourcesByTypeFilter{
		resourceType: resourceType,
		filter:       filter,
	}
}

// Filter applies a filter to multiple resources with a given resource type.
func (f *ResourcesByTypeFilter) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	current := tfwrite.NewFile(inFile)
	resources := current.FindResourcesByType(f.resourceType)
	for _, resource := range resources {
		next, err := f.filter.ResourceFilter(current, resource)
		if err != nil {
			return nil, err
		}
		current = next
	}
	return current.Raw(), nil
}
