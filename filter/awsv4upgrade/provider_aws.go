package awsv4upgrade

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/minamijoyo/hcledit/editor"
	"github.com/minamijoyo/tfedit/tfeditor"
	"github.com/minamijoyo/tfedit/tfwrite"
)

// ProviderAWSFilter is a filter implementation for upgrading arguments of
// provider aws block to v4.
// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#new-provider-arguments
type ProviderAWSFilter struct {
	filters []tfeditor.ProviderFilter
}

var _ editor.Filter = (*ProviderAWSFilter)(nil)

// NewProviderAWSFilter creates a new instance of ProviderAWSFilter.
func NewProviderAWSFilter() editor.Filter {
	filters := []tfeditor.ProviderFilter{
		&ProviderAWSS3ForcePathStyleFilter{},
	}
	return &ProviderAWSFilter{filters: filters}
}

// Filter upgrades arguments of provider aws block to v4.
// Some rules have not been implemented yet.
func (f *ProviderAWSFilter) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	m := tfeditor.NewProvidersByTypeFilter("aws", f)
	return m.Filter(inFile)
}

// ProviderFilter upgrades arguments of provider aws block to v4.
// Some rules have not been implemented yet.
func (f *ProviderAWSFilter) ProviderFilter(inFile *tfwrite.File, provider *tfwrite.Provider) (*tfwrite.File, error) {
	m := tfeditor.NewMultiProviderFilter(f.filters)
	return m.ProviderFilter(inFile, provider)
}
