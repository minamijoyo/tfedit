package tfwrite

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// Provider represents a provider block.
// It implements the Block interface.
type Provider struct {
	*block
}

var _ Block = (*Provider)(nil)

// NewProvider creates a new instance of Provider.
func NewProvider(block *hclwrite.Block) *Provider {
	b := newBlock(block)
	return &Provider{block: b}
}

// NewEmptyProvider creates a new Provider with an empty body.
func NewEmptyProvider(providerType string) *Provider {
	block := hclwrite.NewBlock("provider", []string{providerType})
	return NewProvider(block)
}

// SchemaType returns a type of provider.
// It returns the first label of block.
// Note that it's not the same as the *hclwrite.Block.Type().
func (p *Provider) SchemaType() string {
	labels := p.block.raw.Labels()
	return labels[0]
}
