package tfeditor

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/minamijoyo/hcledit/editor"
	"github.com/minamijoyo/tfedit/tfwrite"
)

// ProviderFilter is an interface which reads Terraform configuration and
// rewrite a given provider, and writes Terraform configuration.
type ProviderFilter interface {
	// ProviderFilter reads Terraform configuration and rewrite a given provider,
	// and writes Terraform configuration.
	ProviderFilter(*tfwrite.File, *tfwrite.Provider) (*tfwrite.File, error)
}

// MultiProviderFilter is a ProviderFilter implementation which applies
// multiple provider filters to a given provider in sequence.
type MultiProviderFilter struct {
	filters []ProviderFilter
}

var _ ProviderFilter = (*MultiProviderFilter)(nil)

// NewMultiProviderFilter creates a new instance of MultiProviderFilter.
func NewMultiProviderFilter(filters []ProviderFilter) ProviderFilter {
	return &MultiProviderFilter{
		filters: filters,
	}
}

// ProviderFilter applies multiple filters to a given provider in sequence.
func (f *MultiProviderFilter) ProviderFilter(inFile *tfwrite.File, provider *tfwrite.Provider) (*tfwrite.File, error) {
	current := inFile
	for _, f := range f.filters {
		next, err := f.ProviderFilter(current, provider)
		if err != nil {
			return nil, err
		}
		current = next
	}
	return current, nil
}

// ProvidersByTypeFilter is a Filter implementation for applying a filter to
// multiple providers with a given provider type.
type ProvidersByTypeFilter struct {
	providerType string
	filter       ProviderFilter
}

var _ editor.Filter = (*ProvidersByTypeFilter)(nil)

// NewProvidersByTypeFilter creates a new instance of ProvidersByTypeFilter.
func NewProvidersByTypeFilter(providerType string, filter ProviderFilter) editor.Filter {
	return &ProvidersByTypeFilter{
		providerType: providerType,
		filter:       filter,
	}
}

// Filter applies a filter to multiple providers with a given provider type.
func (f *ProvidersByTypeFilter) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	current := tfwrite.NewFile(inFile)
	providers := current.FindProvidersByType(f.providerType)
	for _, provider := range providers {
		next, err := f.filter.ProviderFilter(current, provider)
		if err != nil {
			return nil, err
		}
		current = next
	}
	return current.Raw(), nil
}
