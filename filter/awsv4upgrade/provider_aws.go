package awsv4upgrade

import (
	"github.com/minamijoyo/tfedit/tfeditor"
	"github.com/minamijoyo/tfedit/tfwrite"
)

// ProviderAWSFilter is a filter implementation for upgrading arguments of
// provider aws block to v4.
// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#new-provider-arguments
type ProviderAWSFilter struct {
	filters []tfeditor.BlockFilter
}

var _ tfeditor.BlockFilter = (*ProviderAWSFilter)(nil)

// NewProviderAWSFilter creates a new instance of ProviderAWSFilter.
func NewProviderAWSFilter() tfeditor.BlockFilter {
	filters := []tfeditor.BlockFilter{

		tfeditor.ProviderFilterFunc(AWSS3ForcePathStyleProviderFilter),
	}
	return &ProviderAWSFilter{filters: filters}
}

// BlockFilter upgrades arguments of provider aws block to v4.
// Some rules have not been implemented yet.
func (f *ProviderAWSFilter) BlockFilter(inFile *tfwrite.File, block tfwrite.Block) (*tfwrite.File, error) {
	m := tfeditor.NewMultiBlockFilter(f.filters)
	return m.BlockFilter(inFile, block)
}
