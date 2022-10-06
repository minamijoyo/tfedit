package awsv4upgrade

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/minamijoyo/hcledit/editor"
)

// AllFilter is a filter implementation for upgrading configurations
// to AWS provider v4.
// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade
type AllFilter struct {
}

var _ editor.Filter = (*AllFilter)(nil)

// NewAllFilter creates a new instance of AllFilter.
func NewAllFilter() editor.Filter {
	return &AllFilter{}
}

// Filter upgrades configurations to AWS provider v4.
func (f *AllFilter) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	m := editor.NewMultiFilter([]editor.Filter{
		NewProviderAWSFilter(),
		NewAWSS3BucketFilter(),
	})
	return m.Filter(inFile)
}
