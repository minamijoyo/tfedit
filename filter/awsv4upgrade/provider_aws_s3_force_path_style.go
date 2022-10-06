package awsv4upgrade

import (
	"github.com/minamijoyo/tfedit/tfeditor"
	"github.com/minamijoyo/tfedit/tfwrite"
)

// ProviderAWSS3ForcePathStyleFilter is a filter implementation for upgrading
// the s3_force_path_style argument of provider aws block.
// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#s3_use_path_style
type ProviderAWSS3ForcePathStyleFilter struct{}

var _ tfeditor.ProviderFilter = (*ProviderAWSS3ForcePathStyleFilter)(nil)

// NewProviderAWSS3ForcePathStyleFilter creates a new instance of
// ProviderAWSS3ForcePathStyleFilter.
func NewProviderAWSS3ForcePathStyleFilter() tfeditor.ProviderFilter {
	return &ProviderAWSS3ForcePathStyleFilter{}
}

// ProviderFilter is a filter implementation for upgrading the
// s3_force_path_style argument of provider aws block.
func (f *ProviderAWSS3ForcePathStyleFilter) ProviderFilter(inFile *tfwrite.File, provider *tfwrite.Provider) (*tfwrite.File, error) {
	oldAttribute := "s3_force_path_style"
	newAttribute := "s3_use_path_style"

	// Rename a s3_force_path_style attribute to s3_use_path_style.
	attr := provider.GetAttribute(oldAttribute)
	if attr != nil {
		provider.SetAttributeRaw(newAttribute, attr.ValueAsTokens())
		provider.RemoveAttribute(oldAttribute)
	}

	return inFile, nil
}
